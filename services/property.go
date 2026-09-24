package services

import (
	"encoding/json"
	"fmt"
	"os"
	"rental-property-api/models"
	"rental-property-api/utils"
)

var sourceData models.Source

func LoadData(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read source data file %q: %w", path, err)
	}

	var loaded models.Source
	if err := json.Unmarshal(data, &loaded); err != nil {
		return fmt.Errorf("failed to unmarshal source data %q: %w", path, err)
	}

	sourceData = loaded
	return nil
}

func Transform(record models.SourceRecord) models.ResponseItem {
	breadcrumbs := []string{}
	for _, category := range record.Categories {
		breadcrumbs = append(breadcrumbs, category.Name)
	}

	return models.ResponseItem{
		ID:        record.ID,
		Feed:      record.Feed,
		Published: record.Published,
		GeoInfo: models.GeoInfo{
			Breadcrumbs: breadcrumbs,
			City:        record.City,
			Country:     record.Country,
			CountryCode: record.CountryCode,
			Name:        record.Display,
			LocationID:  record.LocationID,
			Lat:         record.LonLat.Coordinates[1],
			Lon:         record.LonLat.Coordinates[0],
			State:       record.State,
			StateAbbr:   record.StateAbbr,
		},
		Property: models.Property{
			Amenities:    record.AmenityCategories,
			Name:         record.PropertyName,
			Slug:         record.PropertySlug,
			PropertyType: record.PropertyTypeCategory,
			Price:        record.UsdPrice,
			ReviewScore:  record.ReviewScoreGeneral,
			StarRating:   record.StarRating,
			Counts: models.PropertyCounts{
				Bathroom:  record.BathroomCount,
				Bedroom:   record.BedroomCount,
				Reviews:   record.NumberOfReview,
				Occupancy: record.Occupancy,
			},
			Image: models.PropertyImage{
				Count:  len(record.Images),
				Images: record.Images,
			},
		},
	}
}

func TransformAll(records models.Source) []models.ResponseItem {
	response := []models.ResponseItem{}
	for _, record := range records {
		response = append(response, Transform(record))
	}
	return response
}

func GetPropertyByID(id string) (models.ResponseItem, error) {
	for _, record := range sourceData {
		if record.ID == id {
			return Transform(record), nil
		}
	}
	return models.ResponseItem{}, fmt.Errorf("Property not found")
}

func GetAllProperties(filterParams utils.FilterParams) models.Response {
	filtered := ApplyFilters(sourceData, filterParams)

	if filterParams.Limit != nil && *filterParams.Limit < len(filtered) {
		filtered = filtered[:*filterParams.Limit]
	}

	items := TransformAll(filtered)

	return models.Response{
		Result: models.Result{
			Count: len(items),
			Items: items,
		},
	}
}

func ApplyFilters(source models.Source, filterParams utils.FilterParams) []models.SourceRecord {
	result := make([]models.SourceRecord, 0, len(source))

	for _, record := range source {
		if !matchesFilters(record, filterParams) {
			continue
		}
		result = append(result, record)
	}

	return result
}

func matchesFilters(record models.SourceRecord, params utils.FilterParams) bool {
	if params.MinPrice != nil && record.UsdPrice < *params.MinPrice {
		return false
	}
	if params.MaxPrice != nil && record.UsdPrice > *params.MaxPrice {
		return false
	}
	if params.MinStarRating != nil && record.StarRating < *params.MinStarRating {
		return false
	}
	if params.MinReviewScore != nil && record.ReviewScoreGeneral < *params.MinReviewScore {
		return false
	}
	if params.MinReviews != nil && record.NumberOfReview < *params.MinReviews {
		return false
	}
	if params.Published != nil && record.Published != *params.Published {
		return false
	}
	if params.PropertyType != nil && record.PropertyTypeCategory != *params.PropertyType {
		return false
	}
	if params.Feed != nil && record.Feed != *params.Feed {
		return false
	}
	if params.MinBedroom != nil && record.BedroomCount < *params.MinBedroom {
		return false
	}
	if params.Amenities != nil && len(*params.Amenities) > 0 && !hasAnyAmenity(record.AmenityCategories, *params.Amenities) {
		return false
	}
	return true
}

func hasAnyAmenity(recordAmenities, wantedAmenities []string) bool {
	amenitySet := make(map[string]bool, len(recordAmenities))
	for _, amenity := range recordAmenities {
		amenitySet[amenity] = true
	}
	for _, wantedAmenity := range wantedAmenities {
		if amenitySet[wantedAmenity] {
			return true
		}
	}
	return false
}
