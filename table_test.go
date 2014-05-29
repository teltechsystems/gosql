package gosql

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestTableGetAlias(t *testing.T) {
	Convey("A basic table name should be its own alias", t, func() {
		table := table{"tableName", []string{}}

		So(table.GetAlias(), ShouldEqual, "tableName")
	})

	Convey("An aliased table name should return its abbreviation", t, func() {
		table := table{"tableName tn", []string{}}

		So(table.GetAlias(), ShouldEqual, "tn")
	})
}
