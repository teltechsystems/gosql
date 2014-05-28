package gosql

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestQuery(t *testing.T) {
	Convey("Without a valid FROM, the query string should be blank", t, func() {
		query := &Query{}
		So(query.String(), ShouldEqual, "")
	})
}

func TestQueryJoin(t *testing.T) {
	Convey("With a single join, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"*"})

		So(len(query.joins), ShouldEqual, 0)
		query.Join(INNER_JOIN, "payments", "payments.user_id = users.id", []string{"amount"})
		So(len(query.joins), ShouldEqual, 1)

		// A few assertions about the join
		So(query.joins[0].joinType, ShouldEqual, INNER_JOIN)
		So(query.joins[0].table.tableName, ShouldEqual, "payments")
		So(query.joins[0].table.columns, ShouldResemble, []string{"amount"})
		So(query.joins[0].predicate, ShouldEqual, "payments.user_id = users.id")

		So(query.String(), ShouldEqual, "SELECT *, amount FROM users INNER JOIN payments ON payments.user_id = users.id")
	})

	Convey("With a single join chained with a where, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"*"}).
			Join(INNER_JOIN, "payments", "payments.user_id = users.id", []string{"amount"}).
			Where("payments.amount > ? AND payments.is_approved", 10)

		So(query.String(), ShouldEqual, "SELECT *, amount FROM users INNER JOIN payments ON payments.user_id = users.id WHERE payments.amount > ? AND payments.is_approved")
	})
}

func TestQueryInnerJoin(t *testing.T) {
	Convey("With a single join chained with a where, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"*"}).
			InnerJoin("payments", "payments.user_id = users.id", []string{"amount"}).
			Where("payments.amount > ? AND payments.is_approved", 10)

		So(query.String(), ShouldEqual, "SELECT *, amount FROM users INNER JOIN payments ON payments.user_id = users.id WHERE payments.amount > ? AND payments.is_approved")
	})
}

func TestQueryLeftJoin(t *testing.T) {
	Convey("With a single join chained with a where, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"*"}).
			LeftJoin("payments", "payments.user_id = users.id", []string{"amount"}).
			Where("payments.amount > ? AND payments.is_approved", 10)

		So(query.String(), ShouldEqual, "SELECT *, amount FROM users LEFT JOIN payments ON payments.user_id = users.id WHERE payments.amount > ? AND payments.is_approved")
	})
}

func TestQueryWhere(t *testing.T) {
	Convey("With a single where condition, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"*"})

		So(len(query.whereParts), ShouldEqual, 0)
		query.Where("first_name = ?", "Bryan")
		So(len(query.whereParts), ShouldEqual, 1)

		So(query.whereParts[0].args[0].(string), ShouldEqual, "Bryan")

		So(query.String(), ShouldEqual, "SELECT * FROM users WHERE first_name = ?")
	})

	Convey("With a single where condition containing multiple arguments, a valid query should be returned", t, func() {
		query := &Query{}
		query.From("users", []string{"*"})

		So(len(query.whereParts), ShouldEqual, 0)
		query.Where("first_name = ? AND last_name = ?", "Bryan", "Moyles")
		So(len(query.whereParts), ShouldEqual, 1)

		So(query.whereParts[0].args[0].(string), ShouldEqual, "Bryan")
		So(query.whereParts[0].args[1].(string), ShouldEqual, "Moyles")

		So(query.String(), ShouldEqual, "SELECT * FROM users WHERE first_name = ? AND last_name = ?")
	})

	Convey("With multiple where conditions, a valid query should be returned", t, func() {
		query := &Query{}

		So(len(query.whereParts), ShouldEqual, 0)
		query.From("users", []string{"*"}).
			Where("first_name = ?", "Bryan").
			Where("last_name = ?", "Moyles")
		So(len(query.whereParts), ShouldEqual, 2)

		So(query.whereParts[0].args[0].(string), ShouldEqual, "Bryan")
		So(query.whereParts[1].args[0].(string), ShouldEqual, "Moyles")

		So(query.String(), ShouldEqual, "SELECT * FROM users WHERE first_name = ? AND last_name = ?")
	})
}
