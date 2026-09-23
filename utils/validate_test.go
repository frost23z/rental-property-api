package utils

import (
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func ptrTo[T any](v T) *T { return &v }

func TestParseFilterParams_Invalid(t *testing.T) {
	const propertyTypeErr = "invalid property_type: must be one of Hotel, House, Apartment, Villa, Resort, Hostel"

	// Each case has exactly one problem, so map iteration order cannot change the error.
	tests := []struct {
		name    string
		query   string
		wantErr string
	}{
		{"unknown parameter", "colour=red", "invalid query parameter: colour"},
		{"parameter names are case sensitive", "Limit=5", "invalid query parameter: Limit"},
		{"repeated parameter", "feed=11&feed=12", "multiple values for query parameter: feed"},
		{"empty value", "limit=", "empty value for query parameter: limit"},
		{"blank value", "min_price=%20", "empty value for query parameter: min_price"},

		{"non numeric price", "min_price=cheap", "invalid min_price: must be a number"},
		{"NaN price", "max_price=NaN", "invalid max_price: must be a number"},
		{"Inf price", "min_price=Inf", "invalid min_price: must be a number"},
		{"negative Inf price", "min_price=-Inf", "invalid min_price: must be a number"},
		{"negative price", "min_price=-1", "invalid min_price: must not be negative"},
		{"negative max price", "max_price=-0.5", "invalid max_price: must not be negative"},
		{"inverted price range", "min_price=200&max_price=100", "invalid price range: min_price must not exceed max_price"},

		{"non integer star rating", "min_star_rating=3.5", "invalid min_star_rating: must be an integer"},
		{"negative star rating", "min_star_rating=-1", "invalid min_star_rating: must not be negative"},
		{"non numeric review score", "min_review_score=high", "invalid min_review_score: must be a number"},
		{"NaN review score", "min_review_score=NaN", "invalid min_review_score: must be a number"},
		{"negative review score", "min_review_score=-0.5", "invalid min_review_score: must not be negative"},
		{"non integer min reviews", "min_reviews=1.5", "invalid min_reviews: must be an integer"},
		{"negative min reviews", "min_reviews=-5", "invalid min_reviews: must not be negative"},
		{"non integer bedroom", "min_bedroom=two", "invalid min_bedroom: must be an integer"},
		{"negative bedroom", "min_bedroom=-1", "invalid min_bedroom: must not be negative"},

		{"published 1", "published=1", "invalid published: must be true or false"},
		{"published 0", "published=0", "invalid published: must be true or false"},
		{"published is case sensitive", "published=TRUE", "invalid published: must be true or false"},
		{"published yes", "published=yes", "invalid published: must be true or false"},

		{"unknown property type", "property_type=Castle", propertyTypeErr},
		{"property type is case sensitive", "property_type=hotel", propertyTypeErr},

		{"non integer feed", "feed=abc", "invalid feed: must be an integer"},
		{"decimal feed", "feed=11.0", "invalid feed: must be an integer"},
		{"unsupported feed", "feed=13", "invalid feed: must be one of 11, 12, 22, 24"},
		{"negative feed", "feed=-11", "invalid feed: must be one of 11, 12, 22, 24"},

		{"only separators in amenities", "amenities=,%20,", "invalid amenities: must contain at least one amenity"},
		{"only commas in amenities", "amenities=,,,", "invalid amenities: must contain at least one amenity"},

		{"non integer limit", "limit=abc", "invalid limit: must be an integer"},
		{"decimal limit", "limit=2.5", "invalid limit: must be an integer"},
		{"zero limit", "limit=0", "invalid limit: must be at least 1"},
		{"negative limit", "limit=-2", "invalid limit: must be at least 1"},
	}

	Convey("Subject: ParseFilterParams rejects invalid input\n", t, func() {
		for _, tt := range tests {
			tt := tt
			Convey(tt.name, func() {
				query, err := url.ParseQuery(tt.query)
				So(err, ShouldBeNil)

				_, err = ParseFilterParams(query)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, tt.wantErr)
			})
		}
	})
}

func TestParseFilterParams_Valid(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  FilterParams
	}{
		{"no params", "", FilterParams{}},
		{"min price", "min_price=50", FilterParams{MinPrice: ptrTo(50.0)}},
		{"min price zero", "min_price=0", FilterParams{MinPrice: ptrTo(0.0)}},
		{"scientific notation", "min_price=1e2", FilterParams{MinPrice: ptrTo(100.0)}},
		{"decimal max price", "max_price=150.5", FilterParams{MaxPrice: ptrTo(150.5)}},
		{"equal min and max price", "min_price=75&max_price=75", FilterParams{MinPrice: ptrTo(75.0), MaxPrice: ptrTo(75.0)}},
		{"surrounding whitespace is trimmed", "min_price=%2050%20", FilterParams{MinPrice: ptrTo(50.0)}},
		{"min star rating", "min_star_rating=3", FilterParams{MinStarRating: ptrTo(3)}},
		{"min review score", "min_review_score=7.5", FilterParams{MinReviewScore: ptrTo(7.5)}},
		{"min reviews zero", "min_reviews=0", FilterParams{MinReviews: ptrTo(0)}},
		{"min bedroom", "min_bedroom=2", FilterParams{MinBedroom: ptrTo(2)}},
		{"published true", "published=true", FilterParams{Published: ptrTo(true)}},
		{"published false", "published=false", FilterParams{Published: ptrTo(false)}},
		{"property type is trimmed", "property_type=%20Hotel%20", FilterParams{PropertyType: ptrTo("Hotel")}},
		{"amenities are trimmed and empty items dropped", "amenities=Internet,%20Parking,,", FilterParams{Amenities: ptrTo([]string{"Internet", "Parking"})}},
		{"amenity names may contain spaces", "amenities=Air%20Conditioner,Pool", FilterParams{Amenities: ptrTo([]string{"Air Conditioner", "Pool"})}},
		{"single amenity", "amenities=Pool", FilterParams{Amenities: ptrTo([]string{"Pool"})}},
		{"limit of one", "limit=1", FilterParams{Limit: ptrTo(1)}},
		{
			"all parameters together",
			"min_price=50&max_price=150.5&min_star_rating=3&min_review_score=7.5&min_reviews=10&published=false&property_type=Hotel&feed=11&min_bedroom=2&amenities=Internet,%20Parking,&limit=5",
			FilterParams{
				MinPrice:       ptrTo(50.0),
				MaxPrice:       ptrTo(150.5),
				MinStarRating:  ptrTo(3),
				MinReviewScore: ptrTo(7.5),
				MinReviews:     ptrTo(10),
				Published:      ptrTo(false),
				PropertyType:   ptrTo("Hotel"),
				Feed:           ptrTo(11),
				MinBedroom:     ptrTo(2),
				Amenities:      ptrTo([]string{"Internet", "Parking"}),
				Limit:          ptrTo(5),
			},
		},
	}

	Convey("Subject: ParseFilterParams parses valid input\n", t, func() {
		for _, tt := range tests {
			tt := tt
			Convey(tt.name, func() {
				query, err := url.ParseQuery(tt.query)
				So(err, ShouldBeNil)

				params, err := ParseFilterParams(query)
				So(err, ShouldBeNil)
				So(params, ShouldResemble, tt.want)
			})
		}
	})
}

func TestParseFilterParams_AllowedValues(t *testing.T) {
	Convey("Subject: Every allowed property type is accepted\n", t, func() {
		for _, propertyType := range []string{"Hotel", "House", "Apartment", "Villa", "Resort", "Hostel"} {
			propertyType := propertyType
			Convey(propertyType, func() {
				params, err := ParseFilterParams(url.Values{"property_type": {propertyType}})
				So(err, ShouldBeNil)
				So(params.PropertyType, ShouldNotBeNil)
				So(*params.PropertyType, ShouldEqual, propertyType)
			})
		}
	})

	Convey("Subject: Every allowed feed is accepted\n", t, func() {
		for _, feed := range []struct {
			raw  string
			want int
		}{{"11", 11}, {"12", 12}, {"22", 22}, {"24", 24}} {
			feed := feed
			Convey(feed.raw, func() {
				params, err := ParseFilterParams(url.Values{"feed": {feed.raw}})
				So(err, ShouldBeNil)
				So(params.Feed, ShouldNotBeNil)
				So(*params.Feed, ShouldEqual, feed.want)
			})
		}
	})
}
