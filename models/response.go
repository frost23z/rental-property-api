package models

type ResponseItem struct {
	ID        string
	Feed      int
	Published bool
	GeoInfo   struct {
		Breadcrumbs []string // TODO: Simple string array or stuct with LocationID, Name?
		City        string
		Country     string
		CountryCode string
		Name        string
		LocationID  string
		Lat         float64
		Lon         float64
		State       string
		StateAbbr   *string
	}
	Property struct {
		Amenities    []string
		Name         string
		Slug         string
		PropertyType string
		Price        float64
		ReviewScore  float64
		StarRating   int
		Counts       struct {
			Bathroom  int
			Bedroom   int
			Reviews   int
			Occupancy int
		}
		Image struct {
			Count  int
			Images []string
		}
	}
}

type Response struct {
	Result struct {
		Count int
		Items []ResponseItem
	}
}
