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

	return params, nil
}
