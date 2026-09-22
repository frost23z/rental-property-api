package services

import (
	"encoding/json"
	"fmt"
	"os"
	"rental-property-api/models"

	"github.com/beego/beego/v2/core/logs"
)

var sourceData models.Source

func init() {
	if err := LoadData("../data/rental_properties.json"); err != nil {
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
