package utils

import (
	"fmt"
	"math"
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

	paramValues, err := initialValidate(query)
	if err != nil {
		return params, err
	}

	if minPrice, ok := paramValues["min_price"]; ok {
		if params.MinPrice, err = parseFloatParam("min_price", minPrice); err != nil {
			return params, err
		}
	}

	if maxPrice, ok := paramValues["max_price"]; ok {
		if params.MaxPrice, err = parseFloatParam("max_price", maxPrice); err != nil {
			return params, err
		}
	}

	if params.MinPrice != nil && params.MaxPrice != nil && *params.MinPrice > *params.MaxPrice {
		return params, fmt.Errorf("invalid price range: min_price must not exceed max_price")
	}

	if minStarRating, ok := paramValues["min_star_rating"]; ok {
		if params.MinStarRating, err = parseIntParam("min_star_rating", minStarRating); err != nil {
			return params, err
		}
	}

	if minReviewScore, ok := paramValues["min_review_score"]; ok {
		if params.MinReviewScore, err = parseFloatParam("min_review_score", minReviewScore); err != nil {
			return params, err
		}
	}

	if minReviews, ok := paramValues["min_reviews"]; ok {
		if params.MinReviews, err = parseIntParam("min_reviews", minReviews); err != nil {
			return params, err
		}
	}

	if published := paramValues["published"]; published != "" {
		parsed, err := strconv.ParseBool(published)
		if err != nil {
			return params, fmt.Errorf("invalid published: must be true or false")
		}
		params.Published = &parsed
	}

	if propertyType := paramValues["property_type"]; propertyType != "" {
		if !validPropertyTypes[propertyType] {
			return params, fmt.Errorf("invalid property_type: must be one of Hotel, House, Apartment, Villa, Resort, Hostel")
		}
		params.PropertyType = &propertyType
	}

	if feed := paramValues["feed"]; feed != "" {
		parsed, err := strconv.Atoi(feed)
		if err != nil {
			return params, fmt.Errorf("invalid feed: must be an integer")
		}
		if !validFeeds[parsed] {
			return params, fmt.Errorf("invalid feed: must be one of 11, 12, 22, 24")
		}
		params.Feed = &parsed
	}

	if minBedroom, ok := paramValues["min_bedroom"]; ok {
		if params.MinBedroom, err = parseIntParam("min_bedroom", minBedroom); err != nil {
			return params, err
		}
	}

	if amenitiesRaw := paramValues["amenities"]; amenitiesRaw != "" {
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

	if limit := paramValues["limit"]; limit != "" {
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

func initialValidate(query url.Values) (map[string]string, error) {
	values := make(map[string]string, len(query))
	for name, rawValues := range query {
		if !validQueryParams[name] {
			return nil, fmt.Errorf("invalid query parameter: %s", name)
		}
		if len(rawValues) > 1 {
			return nil, fmt.Errorf("multiple values for query parameter: %s", name)
		}
		value := strings.TrimSpace(rawValues[0])
		if value == "" {
			return nil, fmt.Errorf("empty value for query parameter: %s", name)
		}
		values[name] = value
	}
	return values, nil
}

func parseFloatParam(name, raw string) (*float64, error) {
	value, err := strconv.ParseFloat(raw, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		return nil, fmt.Errorf("invalid %s: must be a number", name)
	}
	if value < 0 {
		return nil, fmt.Errorf("invalid %s: must not be negative", name)
	}
	return &value, nil
}

func parseIntParam(name, raw string) (*int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: must be an integer", name)
	}
	if value < 0 {
		return nil, fmt.Errorf("invalid %s: must not be negative", name)
	}
	return &value, nil
}
