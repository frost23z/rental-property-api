package utils

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type FilterParams struct {
	MinPrice       *float64
	MaxPrice       *float64
	MinStarRating  *int
	MinReviewScore *float64
	MinReviews     *int
	Published      *bool
	PropertyType   *string
	Feed           *int
	MinBedroom     *int
	Amenities      *[]string
	Limit          *int
}

var validQueryParams = map[string]bool{
	"min_price":        true,
	"max_price":        true,
	"min_star_rating":  true,
	"min_review_score": true,
	"min_reviews":      true,
	"published":        true,
	"property_type":    true,
	"feed":             true,
	"min_bedroom":      true,
	"amenities":        true,
	"limit":            true,
}

var validPropertyTypes = map[string]bool{
	"Hotel":     true,
	"House":     true,
	"Apartment": true,
	"Villa":     true,
	"Resort":    true,
	"Hostel":    true,
}

var validFeeds = map[int]bool{
	11: true,
	12: true,
	22: true,
	24: true,
}

func ParseFilterParams(query url.Values) (FilterParams, error) {
	var params FilterParams

	paramValues := make(map[string]string, len(query))
	for paramName, rawValues := range query {
		if !validQueryParams[paramName] {
			return params, fmt.Errorf("invalid query parameter: %s", paramName)
		}
		if len(rawValues) > 1 {
			return params, fmt.Errorf("multiple values for query parameter: %s", paramName)
		}
		trimmedValue := strings.TrimSpace(rawValues[0])
		if trimmedValue == "" {
			return params, fmt.Errorf("empty value for query parameter: %s", paramName)
		}
		paramValues[paramName] = trimmedValue
	}

	if minPrice := query.Get("min_price"); minPrice != "" {
		parsed, err := strconv.ParseFloat(minPrice, 64)
		if err != nil {
			return params, fmt.Errorf("invalid min_price: must be a number")
		}
		if parsed < 0 {
			return params, fmt.Errorf("invalid min_price: must not be negative")
		}
		params.MinPrice = &parsed
	}

	if maxPrice := query.Get("max_price"); maxPrice != "" {
		parsed, err := strconv.ParseFloat(maxPrice, 64)
		if err != nil {
			return params, fmt.Errorf("invalid max_price: must be a number")
		}
		if parsed < 0 {
			return params, fmt.Errorf("invalid max_price: must not be negative")
		}
		params.MaxPrice = &parsed
	}

	if params.MinPrice != nil && params.MaxPrice != nil && *params.MinPrice > *params.MaxPrice {
		return params, fmt.Errorf("invalid price range: min_price must not exceed max_price")
	}

	if minStarRating := query.Get("min_star_rating"); minStarRating != "" {
		parsed, err := strconv.Atoi(minStarRating)
		if err != nil {
			return params, fmt.Errorf("invalid min_star_rating: must be an integer")
		}
		if parsed < 0 {
			return params, fmt.Errorf("invalid min_star_rating: must not be negative")
		}
		params.MinStarRating = &parsed
	}

	if minReviewScore := query.Get("min_review_score"); minReviewScore != "" {
		parsed, err := strconv.ParseFloat(minReviewScore, 64)
		if err != nil {
			return params, fmt.Errorf("invalid min_review_score: must be a number")
		}
		if parsed < 0 {
			return params, fmt.Errorf("invalid min_review_score: must not be negative")
		}
		params.MinReviewScore = &parsed
	}

	if minReviews := query.Get("min_reviews"); minReviews != "" {
		parsed, err := strconv.Atoi(minReviews)
		if err != nil {
			return params, fmt.Errorf("invalid min_reviews: must be an integer")
		}
		if parsed < 0 {
			return params, fmt.Errorf("invalid min_reviews: must not be negative")
		}
		params.MinReviews = &parsed
	}

	if published := query.Get("published"); published != "" {
		parsed, err := strconv.ParseBool(published)
		if err != nil {
			return params, fmt.Errorf("invalid published: must be true or false")
		}
		params.Published = &parsed
	}

	if propertyType := query.Get("property_type"); propertyType != "" {
		if !validPropertyTypes[propertyType] {
			return params, fmt.Errorf("invalid property_type: must be one of Hotel, House, Apartment, Villa, Resort, Hostel")
		}
		params.PropertyType = &propertyType
	}

	if feed := query.Get("feed"); feed != "" {
		parsed, err := strconv.Atoi(feed)
		if err != nil {
			return params, fmt.Errorf("invalid feed: must be an integer")
		}
		if !validFeeds[parsed] {
			return params, fmt.Errorf("invalid feed: must be one of 11, 12, 22, 24")
		}
		params.Feed = &parsed
	}

	if minBedroom := query.Get("min_bedroom"); minBedroom != "" {
		parsed, err := strconv.Atoi(minBedroom)
		if err != nil {
			return params, fmt.Errorf("invalid min_bedroom: must be an integer")
		}
		if parsed < 0 {
			return params, fmt.Errorf("invalid min_bedroom: must not be negative")
		}
		params.MinBedroom = &parsed
	}

	if amenitiesRaw := query.Get("amenities"); amenitiesRaw != "" {
		var amenities []string
		for _, amenity := range strings.Split(amenitiesRaw, ",") {
			amenity = strings.TrimSpace(amenity)
			if amenity != "" {
				amenities = append(amenities, amenity)
			}
		}
		if len(amenities) == 0 {
			return params, fmt.Errorf("invalid amenities: must not be empty")
		}
		params.Amenities = &amenities
	}

	if limit := query.Get("limit"); limit != "" {
		parsed, err := strconv.Atoi(limit)
		if err != nil {
			return params, fmt.Errorf("invalid limit: must be an integer")
		}
		if parsed < 0 {
			return params, fmt.Errorf("invalid limit: must not be negative")
		}
		params.Limit = &parsed
	}

	return params, nil
}
