package models

type ResponseItem struct {
	ID        string
	Feed      int
	Published bool
	GeoInfo   GeoInfo
	Property  Property
}

type GeoInfo struct {
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

type Property struct {
	Amenities    []string
	Name         string
	Slug         string
	PropertyType string
	Price        float64
	ReviewScore  float64
	StarRating   int
	Counts       PropertyCounts
	Image        PropertyImage
}

type PropertyCounts struct {
	Bathroom  int
	Bedroom   int
	Reviews   int
	Occupancy int
}

type PropertyImage struct {
	Count  int
	Images []string
}

type Response struct {
	Result struct {
		Count int
		Items []ResponseItem
	}
}
