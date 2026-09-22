package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"rental-property-api/models"
	"rental-property-api/utils"

	"github.com/beego/beego/v2/core/logs"
)

var sourceData models.Source

func init() {
	filePath := filepath.Join(utils.GetRootPath(), "data", "rental_properties.json")

	if err := LoadData(filePath); err != nil {
		logs.Error("failed to load source data: %v", err)
		panic(err)
	}
}

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

func GetAllProperties() models.Response {
	items := TransformAll(sourceData)

	return models.Response{
		Result: models.Result{
			Count: len(items),
			Items: items,
		},
	}
}
