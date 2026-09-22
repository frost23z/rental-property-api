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

func TestGetPropertyByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "Found property with valid ID",
			id:      "BC-1000004",
			wantErr: nil,
		},
		{
			name:    "Property not found with invalid ID",
			id:      "does-not-exist",
			wantErr: fmt.Errorf("Property not found"),
		},
	}

	Convey("Subject: Test GetPropertyByID function\n", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				res, err := GetPropertyByID(tt.id)
				if tt.wantErr != nil {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, tt.wantErr.Error())
				} else {
					So(err, ShouldBeNil)
					So(res.ID, ShouldEqual, tt.id)
				}
			})
		}
	})
}
