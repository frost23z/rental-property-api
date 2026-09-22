package models

import (
	"encoding/json"
	"fmt"
)

type Source []SourceRecord

type SourceRecord struct {
	ID                   string   `json:"id"`
	Feed                 int      `json:"feed"`
	Country              string   `json:"country"`
	CountryCode          string   `json:"country_code"`
	State                string   `json:"state"`
	StateAbbr            *string  `json:"state_abbr"`
	City                 string   `json:"city"`
	Display              string   `json:"display"`
	LocationID           string   `json:"location_id"`
	PropertyName         string   `json:"property_name"`
	PropertySlug         string   `json:"property_slug"`
	PropertyTypeCategory string   `json:"property_type_category"`
	UsdPrice             float64  `json:"usd_price"`
	Occupancy            int      `json:"occupancy"`
	BedroomCount         int      `json:"bedroom_count"`
	BathroomCount        int      `json:"bathroom_count"`
	NumberOfReview       int      `json:"number_of_review"`
	ReviewScoreGeneral   float64  `json:"review_score_general"`
	StarRating           int      `json:"star_rating"`
	AmenityCategories    []string `json:"amenity_categories"`
	LonLat               struct {
		Coordinates []float64 `json:"coordinates"`
	} `json:"lonlat"`
	Categories Categories `json:"categories"`
	Published  bool       `json:"published"`
	Images     []string   `json:"images"`
}

type Categories []Category

type Category struct {
	LocationID string   `json:"LocationID"`
	Name       string   `json:"Name"`
	Type       string   `json:"Type"`
	Slug       string   `json:"Slug"`
	Display    []string `json:"Display"`
}

func (c *Categories) UnmarshalJSON(data []byte) error {
	var encoded string

	if err := json.Unmarshal(data, &encoded); err != nil {
		return fmt.Errorf("failed to decode categories string: %w", err)
	}

	if err := json.Unmarshal([]byte(encoded), (*[]Category)(c)); err != nil {
		return fmt.Errorf("failed to decode categories JSON: %w", err)
	}

	return nil
}
