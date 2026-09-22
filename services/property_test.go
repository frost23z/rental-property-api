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

func TestTransformAll(t *testing.T) {
	sourceDataSlice := sourceData[:3]
	response := TransformAll(sourceDataSlice)

	pretty, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(pretty))

	Convey("Subject: Test Map Source Data to Response\n", t, func() {
		Convey("There should be 3 items in the response", func() {
			So(len(response), ShouldEqual, 3)
		})
	})
}
