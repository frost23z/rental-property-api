package models

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCategories_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantNames []string
		wantErr   string
	}{
		{"valid list", `{"c":"[{\"Name\":\"Japan\"},{\"Name\":\"Tokyo\"}]"}`, []string{"Japan", "Tokyo"}, ""},
		{"empty list", `{"c":"[]"}`, []string{}, ""},
		{"value is an array, not a string", `{"c":[{"Name":"Japan"}]}`, nil, "failed to decode categories string"},
		{"string that is not JSON", `{"c":"not json"}`, nil, "failed to decode categories JSON"},
		{"string holding an object", `{"c":"{}"}`, nil, "failed to decode categories JSON"},
		{"empty string", `{"c":""}`, nil, "failed to decode categories JSON"},
	}

	Convey("Subject: Categories.UnmarshalJSON\n", t, func() {
		for _, tt := range tests {
			tt := tt
			Convey(tt.name, func() {
				var out struct {
					C Categories `json:"c"`
				}
				err := json.Unmarshal([]byte(tt.input), &out)

				if tt.wantErr != "" {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldContainSubstring, tt.wantErr)
					return
				}

				So(err, ShouldBeNil)
				names := make([]string, 0, len(out.C))
				for _, c := range out.C {
					names = append(names, c.Name)
				}
				So(names, ShouldResemble, tt.wantNames)
			})
		}

		Convey("All category fields are decoded", func() {
			var out struct {
				C Categories `json:"c"`
			}
			input := `{"c":"[{\"LocationID\": \"6100001\",\"Name\": \"Shinjuku\",\"Type\": \"city\",\"Slug\": \"japan/tokyo/shinjuku\",\"Display\": [\"japan\",\"tokyo\",\"shinjuku\"]}]"}`
			So(json.Unmarshal([]byte(input), &out), ShouldBeNil)
			So(len(out.C), ShouldEqual, 1)
			So(out.C[0], ShouldResemble, Category{
				LocationID: "6100001",
				Name:       "Shinjuku",
				Type:       "city",
				Slug:       "japan/tokyo/shinjuku",
				Display:    []string{"japan", "tokyo", "shinjuku"},
			})
		})
	})
}

func TestSourceRecord_UnmarshalJSON(t *testing.T) {
	const record = `{
		"id": "BC-1000001",
		"feed": 11,
		"state_abbr": null,
		"property_type_category": "Resort",
		"usd_price": 116.69,
		"review_score_general": 4.2,
		"amenity_categories": ["Breakfast Included", "Pool"],
		"lonlat": {"coordinates": [139.6917, 35.6895]},
		"categories": "[{\"Name\": \"Japan\",\"Type\": \"country\"},{\"Name\": \"Tokyo\",\"Type\": \"state\"}]",
		"published": false,
		"images": ["image-1.jpg", "image-2.jpg"]
	}`

	Convey("Subject: SourceRecord decoding\n", t, func() {
		Convey("Decodes a record shaped like the data file", func() {
			var r SourceRecord
			So(json.Unmarshal([]byte(record), &r), ShouldBeNil)

			So(r.ID, ShouldEqual, "BC-1000001")
			So(r.Feed, ShouldEqual, 11)
			So(r.StateAbbr, ShouldBeNil)
			So(r.PropertyTypeCategory, ShouldEqual, "Resort")
			So(r.UsdPrice, ShouldAlmostEqual, 116.69)
			So(r.ReviewScoreGeneral, ShouldAlmostEqual, 4.2)
			So(r.AmenityCategories, ShouldResemble, []string{"Breakfast Included", "Pool"})
			So(r.LonLat.Coordinates, ShouldResemble, []float64{139.6917, 35.6895})
			So(len(r.Categories), ShouldEqual, 2)
			So(r.Categories[1].Name, ShouldEqual, "Tokyo")
			So(r.Categories[1].Type, ShouldEqual, "state")
			So(r.Published, ShouldBeFalse)
			So(r.Images, ShouldResemble, []string{"image-1.jpg", "image-2.jpg"})
		})

		Convey("A non-null state_abbr becomes a non-nil pointer", func() {
			var r SourceRecord
			So(json.Unmarshal([]byte(`{"state_abbr":"TK","categories":"[]"}`), &r), ShouldBeNil)
			So(r.StateAbbr, ShouldNotBeNil)
			So(*r.StateAbbr, ShouldEqual, "TK")
		})

		Convey("One record with bad categories fails the whole Source", func() {
			var s Source
			err := json.Unmarshal([]byte(`[{"id":"a","categories":"[]"},{"id":"b","categories":"nope"}]`), &s)
			So(err, ShouldNotBeNil)
		})
	})
}
