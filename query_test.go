package gosql

import (
	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"database/sql"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestQuery(t *testing.T) {
	Convey("Without a valid FROM, the query string should be blank", t, func() {
		query := &Query{}
		So(query.String(), ShouldEqual, "")
	})

	Convey("Test a simple query to string", t, func() {
		query := &Query{}
		query.From("users", []string{"id"})
		So(query.String(), ShouldEqual, "SELECT users.id FROM users")
	})

	Convey("Test a simple aliased query to string", t, func() {
		query := &Query{}
		query.From("users u", []string{"id"})
		So(query.String(), ShouldEqual, "SELECT u.id FROM users u")
	})

	Convey("With a function call as a field, the field should not be prefixed", t, func() {
		query := &Query{}
		query.From("users u", []string{"MAX(id)"})
		So(query.String(), ShouldEqual, "SELECT MAX(id) FROM users u")
	})
}

func TestQueryGroupBy(t *testing.T) {
	Convey("With a simple GroupBy, a valid query should be returned", t, func() {
		query := &Query{}

		So(query.groupBy, ShouldEqual, "")
		query.From("users", []string{"id"}).
			GroupBy("users.id")
		So(query.groupBy, ShouldEqual, "users.id")

		So(query.String(), ShouldEqual, "SELECT users.id FROM users GROUP BY users.id")
	})
}

func TestQueryJoin(t *testing.T) {
	Convey("With a single join, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"id"})

		So(len(query.joins), ShouldEqual, 0)
		query.Join(INNER_JOIN, "payments", "payments.user_id = users.id", []string{"amount"})
		So(len(query.joins), ShouldEqual, 1)

		// A few assertions about the join
		So(query.joins[0].joinType, ShouldEqual, INNER_JOIN)
		So(query.joins[0].table.tableName, ShouldEqual, "payments")
		So(query.joins[0].table.columns, ShouldResemble, []string{"amount"})
		So(query.joins[0].wherePart.predicate, ShouldEqual, "payments.user_id = users.id")
		So(len(query.joins[0].wherePart.args), ShouldEqual, 0)

		So(query.String(), ShouldEqual, "SELECT users.id, payments.amount FROM users INNER JOIN payments ON payments.user_id = users.id")
	})

	Convey("With a single join chained with a where, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"id"}).
			Join(INNER_JOIN, "payments", "payments.user_id = users.id", []string{"amount"}).
			Where("payments.amount > ? AND payments.is_approved", 10)

		So(query.String(), ShouldEqual, "SELECT users.id, payments.amount FROM users INNER JOIN payments ON payments.user_id = users.id WHERE (payments.amount > ? AND payments.is_approved)")
	})

	Convey("With a single join with a bounded variadic, a valid query should be returned", t, func() {
		query := &Query{}

		So(len(query.getArgs()), ShouldEqual, 0)
		query.From("users", []string{"id"}).
			Join(INNER_JOIN, "payments", "payments.user_id = users.id AND some_statement = ?", []string{"amount"}, "hello")
		So(len(query.getArgs()), ShouldEqual, 1)

		So(query.String(), ShouldEqual, "SELECT users.id, payments.amount FROM users INNER JOIN payments ON payments.user_id = users.id AND some_statement = ?")
	})
}

func TestQueryInnerJoin(t *testing.T) {
	Convey("With a single join chained with a where, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"id"}).
			InnerJoin("payments", "payments.user_id = users.id", []string{"amount"}).
			Where("payments.amount > ? AND payments.is_approved", 10)

		So(query.String(), ShouldEqual, "SELECT users.id, payments.amount FROM users INNER JOIN payments ON payments.user_id = users.id WHERE (payments.amount > ? AND payments.is_approved)")
	})
}

func TestQueryLeftJoin(t *testing.T) {
	Convey("With a single join chained with a where, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"id"}).
			LeftJoin("payments", "payments.user_id = users.id", []string{"amount"}).
			Where("payments.amount > ? AND payments.is_approved", 10)

		So(query.String(), ShouldEqual, "SELECT users.id, payments.amount FROM users LEFT JOIN payments ON payments.user_id = users.id WHERE (payments.amount > ? AND payments.is_approved)")
	})
}

func TestQueryOrderBy(t *testing.T) {
	Convey("With a simple OrderBy, a valid query should be returned", t, func() {
		query := &Query{}

		So(len(query.orderByParts), ShouldEqual, 0)
		query.From("users", []string{"id"}).
			OrderBy([]string{"users.id ASC"})
		So(len(query.orderByParts), ShouldEqual, 1)

		So(query.String(), ShouldEqual, "SELECT users.id FROM users ORDER BY users.id ASC")
	})

	Convey("With a complex OrderBy, a valid query should be returned", t, func() {
		query := &Query{}

		So(len(query.orderByParts), ShouldEqual, 0)
		query.From("users", []string{"id"}).
			OrderBy([]string{"users.id ASC", "users.first_name DESC"})
		So(len(query.orderByParts), ShouldEqual, 2)

		So(query.String(), ShouldEqual, "SELECT users.id FROM users ORDER BY users.id ASC, users.first_name DESC")
	})
}

func TestQueryUse(t *testing.T) {
	Convey("A db instance should be able to be associated to the query", t, func() {
		db, err := sql.Open("sqlite3", ":memory:")
		So(err, ShouldBeNil)

		query := &Query{}

		So(query.using, ShouldBeNil)
		query.Use(db)
		So(query.using, ShouldNotBeNil)
	})
}

func TestQueryQuery(t *testing.T) {
	Convey("With a query unassociated to a database, an error should be returned", t, func() {
		query := &Query{}

		rows, err := query.Query()
		So(rows, ShouldBeNil)
		So(err, ShouldEqual, MissingDatabase)
	})

	Convey("With a query associated to a database, an error should be returned due to invalid schema", t, func() {
		db, err := sql.Open("sqlite3", ":memory:")
		So(err, ShouldBeNil)

		query := Select().From("users", []string{"id"})
		query.Use(db)

		rows, err := query.Query()
		So(rows, ShouldBeNil)
		So(err, ShouldNotBeNil)
		So(err, ShouldNotEqual, MissingDatabase)
	})
}

func TestQueryQueryRow(t *testing.T) {
	Convey("With a query unassociated to a database, should panic", t, func() {
		query := &Query{}

		defer func() {
			So(recover(), ShouldNotBeNil)
		}()

		query.QueryRow()
	})

	Convey("With a query associated to a database, an error should be returned due to invalid schema", t, func() {
		db, err := sql.Open("sqlite3", ":memory:")
		So(err, ShouldBeNil)

		query := Select().From("users", []string{"id"})
		query.Use(db)

		row := query.QueryRow()
		So(row, ShouldNotBeNil)
	})
}

func TestQueryWhere(t *testing.T) {
	Convey("With a single where condition, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"id"})

		So(len(query.whereParts), ShouldEqual, 0)
		query.Where("first_name = ?", "Bryan")
		So(len(query.whereParts), ShouldEqual, 1)

		So(query.whereParts[0].args[0].(string), ShouldEqual, "Bryan")

		So(query.String(), ShouldEqual, "SELECT users.id FROM users WHERE (first_name = ?)")
	})

	Convey("With a single where condition containing multiple arguments, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"id"})

		So(len(query.whereParts), ShouldEqual, 0)
		query.Where("first_name = ? AND last_name = ?", "Bryan", "Moyles")
		So(len(query.whereParts), ShouldEqual, 1)

		So(query.whereParts[0].args[0].(string), ShouldEqual, "Bryan")
		So(query.whereParts[0].args[1].(string), ShouldEqual, "Moyles")

		So(query.String(), ShouldEqual, "SELECT users.id FROM users WHERE (first_name = ? AND last_name = ?)")
	})

	Convey("With multiple where conditions, a valid query should be returned", t, func() {
		query := &Query{}

		So(len(query.whereParts), ShouldEqual, 0)
		query.From("users", []string{"id"}).
			Where("first_name = ?", "Bryan").
			Where("last_name = ?", "Moyles")
		So(len(query.whereParts), ShouldEqual, 2)

		So(query.whereParts[0].args[0].(string), ShouldEqual, "Bryan")
		So(query.whereParts[1].args[0].(string), ShouldEqual, "Moyles")

		So(query.String(), ShouldEqual, "SELECT users.id FROM users WHERE (first_name = ?) AND (last_name = ?)")
	})
}
