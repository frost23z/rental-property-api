package services

import (
	"errors"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"rental-property-api/models"
	"rental-property-api/utils"

	. "github.com/smartystreets/goconvey/convey"
)

func ptr[T any](v T) *T { return &v }

// customSource is a small deterministic data set. Coordinates are set so
// Transform (which reads Coordinates[0] and [1]) never panics on it.
func customSource() models.Source {
	source := models.Source{
		{ID: "P1", Feed: 11, Published: false, PropertyTypeCategory: "Hotel", UsdPrice: 50, StarRating: 3, ReviewScoreGeneral: 7.5, NumberOfReview: 15, BedroomCount: 1, AmenityCategories: []string{"Internet", "Pool"}},
		{ID: "P2", Feed: 11, Published: true, PropertyTypeCategory: "Hotel", UsdPrice: 100, StarRating: 4, ReviewScoreGeneral: 8.5, NumberOfReview: 40, BedroomCount: 2, AmenityCategories: []string{"Parking"}},
		{ID: "P3", Feed: 12, Published: true, PropertyTypeCategory: "Villa", UsdPrice: 150, StarRating: 5, ReviewScoreGeneral: 9.0, NumberOfReview: 80, BedroomCount: 4, AmenityCategories: []string{"Pool", "Kitchen"}},
		{ID: "P4", Feed: 22, Published: false, PropertyTypeCategory: "House", UsdPrice: 200, StarRating: 2, ReviewScoreGeneral: 6.0, NumberOfReview: 5, BedroomCount: 3, AmenityCategories: []string{"Internet", "Parking", "Kitchen"}},
		{ID: "P5", Feed: 24, Published: true, PropertyTypeCategory: "Hostel", UsdPrice: 20, StarRating: 1, ReviewScoreGeneral: 5.5, NumberOfReview: 0, BedroomCount: 0, AmenityCategories: []string{}},
	}
	for i := range source {
		source[i].LonLat.Coordinates = []float64{float64(100 + i), float64(10 + i)}
	}
	return source
}

// swapSource replaces the in-memory store and returns a func that restores it.
// Use it as: Reset(swapSource(customSource()))
func swapSource(s models.Source) (restore func()) {
	original := sourceData
	sourceData = s
	return func() { sourceData = original }
}

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

func TestLoadData_Errors(t *testing.T) {
	dir := t.TempDir()
	write := func(name, content string) string {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}

	tests := []struct {
		name    string
		content string
	}{
		{"invalid JSON", "not json"},
		{"empty file", ""},
		{"JSON object instead of array", `{"id":"x"}`},
		{"categories is not valid JSON", `[{"id":"x","categories":"not json"}]`},
		{"categories is not a string", `[{"id":"x","categories":[]}]`},
	}

	Convey("Subject: LoadData failure paths\n", t, func() {
		Reset(swapSource(customSource()))
		wantIDs := recordIDs(customSource())

		Convey("Missing file returns a not-exist error and keeps existing data", func() {
			err := LoadData(filepath.Join(dir, "missing.json"))
			So(err, ShouldNotBeNil)
			So(errors.Is(err, os.ErrNotExist), ShouldBeTrue)
			So(recordIDs(sourceData), ShouldResemble, wantIDs)
		})

		for _, tt := range tests {
			Convey(tt.name+" returns an error and keeps existing data", func() {
				err := LoadData(write("case.json", tt.content))
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, "failed to unmarshal source data")
				So(recordIDs(sourceData), ShouldResemble, wantIDs)
			})
		}

		Convey("A valid file replaces the existing data", func() {
			err := LoadData(write("valid.json", `[{"id":"T1","categories":"[]","lonlat":{"coordinates":[1,2]}}]`))
			So(err, ShouldBeNil)
			So(recordIDs(sourceData), ShouldResemble, []string{"T1"})
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

func TestApplyFilters(t *testing.T) {
	tests := []struct {
		name   string
		params utils.FilterParams
		want   []string
	}{
		{"no filters returns everything", utils.FilterParams{}, []string{"P1", "P2", "P3", "P4", "P5"}},
		{"feed and published", utils.FilterParams{Feed: ptr(11), Published: ptr(false)}, []string{"P1"}},
		{"published false", utils.FilterParams{Published: ptr(false)}, []string{"P1", "P4"}},
		{"feed only", utils.FilterParams{Feed: ptr(22)}, []string{"P4"}},
		{"min price is inclusive", utils.FilterParams{MinPrice: ptr(100.0)}, []string{"P2", "P3", "P4"}},
		{"max price is inclusive", utils.FilterParams{MaxPrice: ptr(100.0)}, []string{"P1", "P2", "P5"}},
		{"price range with property type", utils.FilterParams{MinPrice: ptr(50.0), MaxPrice: ptr(150.0), PropertyType: ptr("Hotel")}, []string{"P1", "P2"}},
		{"min star rating is inclusive", utils.FilterParams{MinStarRating: ptr(5)}, []string{"P3"}},
		{"min star rating", utils.FilterParams{MinStarRating: ptr(4)}, []string{"P2", "P3"}},
		{"min review score is inclusive", utils.FilterParams{MinReviewScore: ptr(8.5)}, []string{"P2", "P3"}},
		{"min reviews is inclusive", utils.FilterParams{MinReviews: ptr(40)}, []string{"P2", "P3"}},
		{"min review score and min reviews", utils.FilterParams{MinReviewScore: ptr(8.0), MinReviews: ptr(50)}, []string{"P3"}},
		{"min bedroom zero matches everything", utils.FilterParams{MinBedroom: ptr(0)}, []string{"P1", "P2", "P3", "P4", "P5"}},
		{"min bedroom", utils.FilterParams{MinBedroom: ptr(3)}, []string{"P3", "P4"}},
		{"property type is an exact match", utils.FilterParams{PropertyType: ptr("Villa")}, []string{"P3"}},
		{"property type is case sensitive", utils.FilterParams{PropertyType: ptr("villa")}, []string{}},
		{"amenities match when any one is present", utils.FilterParams{Amenities: &[]string{"Internet", "Parking"}}, []string{"P1", "P2", "P4"}},
		{"amenities with two shared across records", utils.FilterParams{Amenities: &[]string{"Pool", "Kitchen"}}, []string{"P1", "P3", "P4"}},
		{"amenities are case sensitive", utils.FilterParams{Amenities: &[]string{"internet"}}, []string{}},
		{"unknown amenity matches nothing", utils.FilterParams{Amenities: &[]string{"Helipad"}}, []string{}},
		{"feed plus amenities", utils.FilterParams{Feed: ptr(11), Amenities: &[]string{"Internet", "Parking"}}, []string{"P1", "P2"}},
		{"amenity matches but another filter fails", utils.FilterParams{Published: ptr(true), Amenities: &[]string{"Internet"}}, []string{}},
		{"nothing matches", utils.FilterParams{MinPrice: ptr(1000.0)}, []string{}},
		{
			"every filter at once narrows to one record",
			utils.FilterParams{
				MinPrice: ptr(100.0), MaxPrice: ptr(100.0), MinStarRating: ptr(4), MinReviewScore: ptr(8.5),
				MinReviews: ptr(40), Published: ptr(true), PropertyType: ptr("Hotel"), Feed: ptr(11),
				MinBedroom: ptr(2), Amenities: &[]string{"Parking"},
			},
			[]string{"P2"},
		},
	}

	Convey("Subject: ApplyFilters (fixtures)\n", t, func() {
		source := customSource()
		for _, tt := range tests {
			Convey(tt.name, func() {
				got := ApplyFilters(source, tt.params)
				So(got, ShouldNotBeNil)
				So(recordIDs(got), ShouldResemble, tt.want)
			})
		}

		Convey("Does not modify the source", func() {
			ApplyFilters(source, utils.FilterParams{Feed: ptr(11)})
			So(recordIDs(source), ShouldResemble, []string{"P1", "P2", "P3", "P4", "P5"})
		})
	})
}

func TestApplyFilters_RealData(t *testing.T) {
	hasAny := func(have, want []string) bool {
		for _, h := range have {
			if slices.Contains(want, h) {
				return true
			}
		}
		return false
	}

	tests := []struct {
		name   string
		params utils.FilterParams
		match  func(r models.SourceRecord) bool
	}{
		{"min price", utils.FilterParams{MinPrice: ptr(100.0)}, func(r models.SourceRecord) bool { return r.UsdPrice >= 100 }},
		{"max price", utils.FilterParams{MaxPrice: ptr(100.0)}, func(r models.SourceRecord) bool { return r.UsdPrice <= 100 }},
		{"feed", utils.FilterParams{Feed: ptr(11)}, func(r models.SourceRecord) bool { return r.Feed == 11 }},
		{"published", utils.FilterParams{Published: ptr(true)}, func(r models.SourceRecord) bool { return r.Published }},
		{"not published", utils.FilterParams{Published: ptr(false)}, func(r models.SourceRecord) bool { return !r.Published }},
		{"property type", utils.FilterParams{PropertyType: ptr("Hotel")}, func(r models.SourceRecord) bool { return r.PropertyTypeCategory == "Hotel" }},
		{"min star rating", utils.FilterParams{MinStarRating: ptr(4)}, func(r models.SourceRecord) bool { return r.StarRating >= 4 }},
		{"min review score", utils.FilterParams{MinReviewScore: ptr(5.0)}, func(r models.SourceRecord) bool { return r.ReviewScoreGeneral >= 5 }},
		{"min reviews", utils.FilterParams{MinReviews: ptr(100)}, func(r models.SourceRecord) bool { return r.NumberOfReview >= 100 }},
		{"min bedroom", utils.FilterParams{MinBedroom: ptr(3)}, func(r models.SourceRecord) bool { return r.BedroomCount >= 3 }},
		{"amenities (any)", utils.FilterParams{Amenities: &[]string{"Pool", "Gym"}}, func(r models.SourceRecord) bool { return hasAny(r.AmenityCategories, []string{"Pool", "Gym"}) }},
		{
			"combined",
			utils.FilterParams{Feed: ptr(11), MinPrice: ptr(50.0), MaxPrice: ptr(250.0), MinStarRating: ptr(2), Amenities: &[]string{"Pool", "Kitchen"}},
			func(r models.SourceRecord) bool {
				return r.Feed == 11 && r.UsdPrice >= 50 && r.UsdPrice <= 250 && r.StarRating >= 2 &&
					hasAny(r.AmenityCategories, []string{"Pool", "Kitchen"})
			},
		},
	}

	Convey("Subject: ApplyFilters against the real data\n", t, func() {
		for _, tt := range tests {
			Convey(tt.name, func() {
				want := make([]string, 0)
				for _, r := range sourceData {
					if tt.match(r) {
						want = append(want, r.ID)
					}
				}
				So(recordIDs(ApplyFilters(sourceData, tt.params)), ShouldResemble, want)
			})
		}
	})
}

func TestGetAllProperties(t *testing.T) {
	Convey("Subject: GetAllProperties\n", t, func() {
		Reset(swapSource(customSource()))

		tests := []struct {
			name   string
			params utils.FilterParams
			want   []string
		}{
			{"no limit returns all matches", utils.FilterParams{}, []string{"P1", "P2", "P3", "P4", "P5"}},
			{"limit below the number of matches keeps the first N in order", utils.FilterParams{Limit: ptr(2)}, []string{"P1", "P2"}},
			{"limit equal to the number of matches", utils.FilterParams{Limit: ptr(5)}, []string{"P1", "P2", "P3", "P4", "P5"}},
			{"limit above the number of matches", utils.FilterParams{Limit: ptr(50)}, []string{"P1", "P2", "P3", "P4", "P5"}},
			{"limit is applied after filtering", utils.FilterParams{Feed: ptr(11), Limit: ptr(1)}, []string{"P1"}},
			{"filter without limit", utils.FilterParams{Feed: ptr(11)}, []string{"P1", "P2"}},
			{"no matches with a limit", utils.FilterParams{MinPrice: ptr(1000.0), Limit: ptr(3)}, []string{}},
		}

		for _, tt := range tests {
			Convey(tt.name, func() {
				res := GetAllProperties(tt.params)
				So(res.Result.Items, ShouldNotBeNil)
				So(itemIDs(res.Result.Items), ShouldResemble, tt.want)
				So(res.Result.Count, ShouldEqual, len(tt.want))
			})
		}
	})

	Convey("Subject: GetAllProperties on the real data\n", t, func() {
		Convey("No filters returns every record in file order", func() {
			res := GetAllProperties(utils.FilterParams{})
			So(res.Result.Count, ShouldEqual, len(sourceData))
			So(itemIDs(res.Result.Items), ShouldResemble, recordIDs(sourceData))
		})

		Convey("Limit returns a prefix of the full result", func() {
			res := GetAllProperties(utils.FilterParams{Limit: ptr(10)})
			So(res.Result.Count, ShouldEqual, 10)
			So(itemIDs(res.Result.Items), ShouldResemble, recordIDs(sourceData)[:10])
		})
	})
}
