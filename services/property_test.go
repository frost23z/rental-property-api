package services

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadData(t *testing.T) {
	Convey("Subject: Test Load Data\n", t, func() {
		Convey("There should be 100 records", func() {
			So(len(sourceData), ShouldEqual, 100)
		})
	})
}
