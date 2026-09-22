package utils

import (
	"fmt"
	"net/url"
	"strconv"
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

func ParseFilterParams(query url.Values) (FilterParams, error) {
	var params FilterParams

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

	return params, nil
}
