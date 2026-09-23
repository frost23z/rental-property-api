package services

import (
	"testing"

	"rental-property-api/models"

	. "github.com/smartystreets/goconvey/convey"
)

func recordIDs(records []models.SourceRecord) []string {
	ids := make([]string, 0, len(records))
	for _, r := range records {
		ids = append(ids, r.ID)
	}
	return ids
}

func itemIDs(items []models.ResponseItem) []string {
	ids := make([]string, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	return ids
}

func TestLoadData(t *testing.T) {
	Convey("Subject: The provided data file\n", t, func() {
		Convey("There should be 100 records", func() {
			So(len(sourceData), ShouldEqual, 100)
		})

		Convey("Every record has a unique, non-empty ID", func() {
			seen := map[string]bool{}
			var bad []string
			for _, r := range sourceData {
				if r.ID == "" || seen[r.ID] {
					bad = append(bad, r.ID)
				}
				seen[r.ID] = true
			}
			So(bad, ShouldBeEmpty)
		})

		Convey("Every record has at least 2 coordinates (Transform reads [0] and [1])", func() {
			var bad []string
			for _, r := range sourceData {
				if len(r.LonLat.Coordinates) < 2 {
					bad = append(bad, r.ID)
				}
			}
			So(bad, ShouldBeEmpty)
		})

		Convey("No record has a negative price", func() {
			var bad []string
			for _, r := range sourceData {
				if r.UsdPrice < 0 {
					bad = append(bad, r.ID)
				}
			}
			So(bad, ShouldBeEmpty)
		})
	})
}

func TestTransform(t *testing.T) {
	Convey("Subject: Transform\n", t, func() {
		Convey("Maps every field of a real record (BC-1000001)", func() {
			item, err := GetPropertyByID("BC-1000001")
			So(err, ShouldBeNil)

			So(item.ID, ShouldEqual, "BC-1000001")
			So(item.Feed, ShouldEqual, 11)
			So(item.Published, ShouldBeFalse)

			So(item.GeoInfo.Breadcrumbs, ShouldResemble, []string{"Japan", "Tokyo", "Shinjuku"})
			So(item.GeoInfo.City, ShouldEqual, "Shinjuku")
			So(item.GeoInfo.Country, ShouldEqual, "Japan")
			So(item.GeoInfo.CountryCode, ShouldEqual, "JP")
			So(item.GeoInfo.Name, ShouldEqual, "Shinjuku, Japan")
			So(item.GeoInfo.LocationID, ShouldEqual, "6100001")
			So(item.GeoInfo.State, ShouldEqual, "Tokyo")
			So(item.GeoInfo.StateAbbr, ShouldBeNil)
			So(item.GeoInfo.Lat, ShouldAlmostEqual, 35.6895)
			So(item.GeoInfo.Lon, ShouldAlmostEqual, 139.6917)

			So(item.Property.Name, ShouldEqual, "Shinjuku Grand Resort")
			So(item.Property.Slug, ShouldEqual, "shinjuku-grand-resort-0001")
			So(item.Property.PropertyType, ShouldEqual, "Resort")
			So(item.Property.Price, ShouldAlmostEqual, 116.69)
			So(item.Property.ReviewScore, ShouldAlmostEqual, 4.2)
			So(item.Property.StarRating, ShouldEqual, 5)
			So(item.Property.Amenities, ShouldResemble, []string{"Breakfast Included", "Child Friendly", "Pool"})
			So(item.Property.Counts.Bathroom, ShouldEqual, 2)
			So(item.Property.Counts.Bedroom, ShouldEqual, 2)
			So(item.Property.Counts.Reviews, ShouldEqual, 304)
			So(item.Property.Counts.Occupancy, ShouldEqual, 3)
			So(item.Property.Image.Count, ShouldEqual, 5)
			So(len(item.Property.Image.Images), ShouldEqual, 5)
		})

		Convey("Passes a state abbreviation through and swaps lon/lat correctly", func() {
			abbr := "TK"
			record := models.SourceRecord{
				ID:         "X",
				StateAbbr:  &abbr,
				Categories: models.Categories{{Name: "A"}, {Name: "B"}},
				Images:     []string{"a.jpg", "b.jpg"},
			}
			record.LonLat.Coordinates = []float64{1.5, 2.5}

			item := Transform(record)
			So(item.GeoInfo.StateAbbr, ShouldNotBeNil)
			So(*item.GeoInfo.StateAbbr, ShouldEqual, "TK")
			So(item.GeoInfo.Lon, ShouldAlmostEqual, 1.5)
			So(item.GeoInfo.Lat, ShouldAlmostEqual, 2.5)
			So(item.GeoInfo.Breadcrumbs, ShouldResemble, []string{"A", "B"})
			So(item.Property.Image.Count, ShouldEqual, 2)
		})

		Convey("No categories gives empty, non-nil breadcrumbs", func() {
			record := models.SourceRecord{ID: "X"}
			record.LonLat.Coordinates = []float64{1, 2}

			item := Transform(record)
			So(item.GeoInfo.Breadcrumbs, ShouldNotBeNil)
			So(item.GeoInfo.Breadcrumbs, ShouldBeEmpty)
			So(item.Property.Image.Count, ShouldEqual, 0)
		})
	})
}

func TestTransformAll(t *testing.T) {
	Convey("Subject: TransformAll\n", t, func() {
		Convey("Transforms every real record, in order, without panicking", func() {
			response := TransformAll(sourceData)
			So(len(response), ShouldEqual, len(sourceData))
			So(itemIDs(response), ShouldResemble, recordIDs(sourceData))
		})

		Convey("Empty input gives an empty, non-nil slice", func() {
			response := TransformAll(models.Source{})
			So(response, ShouldNotBeNil)
			So(response, ShouldBeEmpty)
		})
	})
}

func TestGetPropertyByID(t *testing.T) {
	tests := []struct {
		name    string
		id      string
		wantErr string
	}{
		{"Property not found with invalid ID", "does-not-exist", "Property not found"},
		{"Lookup is case sensitive", "bc-1000004", "Property not found"},
		{"Lookup does not trim whitespace", " BC-1000004", "Property not found"},
		{"Empty ID", "", "Property not found"},
	}

	Convey("Subject: GetPropertyByID\n", t, func() {
		Convey("Found property with valid ID", func() {
			res, err := GetPropertyByID("BC-1000004")
			So(err, ShouldBeNil)
			So(res.ID, ShouldEqual, "BC-1000004")
			So(res.Property.Name, ShouldEqual, "Umeda Garden House")
			So(res.Published, ShouldBeTrue)
		})

		Convey("Every loaded ID can be looked up", func() {
			var missing []string
			for _, r := range sourceData {
				res, err := GetPropertyByID(r.ID)
				if err != nil || res.ID != r.ID {
					missing = append(missing, r.ID)
				}
			}
			So(missing, ShouldBeEmpty)
		})

		for _, tt := range tests {
			Convey(tt.name, func() {
				res, err := GetPropertyByID(tt.id)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, tt.wantErr)
				So(res.ID, ShouldBeEmpty)
			})
		}
	})
}
