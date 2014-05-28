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
