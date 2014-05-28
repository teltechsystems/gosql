package gosql

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSelect(t *testing.T) {
	Convey("With a simple query", t, func() {
		query := Select()
		So(query, ShouldNotBeNil)
		So(query.from, ShouldBeNil)

		query.From("users", []string{"*"})
		So(query.from, ShouldHaveSameTypeAs, &table{})
		So(query.from.tableName, ShouldEqual, "users")
		So(query.from.columns, ShouldResemble, []string{"*"})

		So(query.String(), ShouldEqual, "SELECT * FROM users")
	})
}
