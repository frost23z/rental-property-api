package services

import (
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLoadData(t *testing.T) {
	pretty, err := json.MarshalIndent(sourceData, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(pretty))

	Convey("Subject: Test Load Data\n", t, func() {
		Convey("There should be 100 records", func() {
			So(len(sourceData), ShouldEqual, 100)
		})
	})
}
