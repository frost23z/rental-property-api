package controllers

import (
	"rental-property-api/models"
	"rental-property-api/services"
	"rental-property-api/utils"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

type PropertyController struct {
	beego.Controller
}

// @Title Get
// @Description find property by id
// @Param id path string true "The id of the property you want to get"
// @Success 200 {object} models.ResponseItem
// @Success 404 {object} models.ErrorResponse
// @router /:id [get]
// Using the @Success annotation for 404 too. Using the @Failure annotation displays {object} models.ErrorResponse as a string, can't render actual JSON object in Swagger UI.
func (p *PropertyController) Get() {
	id := p.Ctx.Input.Param(":id")

	item, err := services.GetPropertyByID(id)
	if err != nil {
		p.Ctx.Output.SetStatus(404)
		p.Data["json"] = models.ErrorResponse{Error: err.Error()}
		p.ServeJSON()
		return
	}

	p.Data["json"] = item
	p.ServeJSON()
}

// @Title GetAll
// @Description get all properties, optionally filtered by the parameters below. All parameters are optional; each may be supplied at most once.
// @Param min_price query number false "Minimum price. Must be >= 0. Must not exceed max_price if both are set."
// @Param max_price query number false "Maximum price. Must be >= 0. Must not be less than min_price if both are set."
// @Param min_star_rating query int false "Minimum star rating. Must be an integer >= 0."
// @Param min_review_score query number false "Minimum review score. Must be >= 0."
// @Param min_reviews query int false "Minimum number of reviews. Must be an integer >= 0."
// @Param published query bool false "Filter by published status. Must be exactly 'true' or 'false'."
// @Param property_type query string false "Property type. Must be one of: Hotel, House, Apartment, Villa, Resort, Hostel."
// @Param feed query int false "Feed id. Must be one of: 11, 12, 22, 24."
// @Param min_bedroom query int false "Minimum number of bedrooms. Must be an integer >= 0."
// @Param amenities query string false "Comma-separated list of amenities (e.g. 'wifi,pool,parking'). Must contain at least one non-empty value after trimming."
// @Param limit query int false "Maximum number of results to return. Must be an integer >= 1."
// @Success 200 {object} models.Response
// @Failure 400 {object} models.ErrorResponse "Returned when an unknown query parameter is supplied, a parameter is repeated, a value is empty/whitespace, or any value fails the validation rules above."
// @router / [get]
func (p *PropertyController) GetAll() {
	query := p.Ctx.Request.URL.Query()

	filterParams, err := utils.ParseFilterParams(query)

	if err != nil {
		logs.Error("Error parsing filter parameters: %v", err)
		p.Ctx.Output.SetStatus(400)
		p.Data["json"] = models.ErrorResponse{Error: err.Error()}
		p.ServeJSON()
		return
	}

	logs.Info(
		"Filter parameters: min_price=%v max_price=%v min_star_rating=%v min_review_score=%v min_reviews=%v published=%v property_type=%v feed=%v min_bedroom=%v amenities=%v limit=%v",
		valueOrNil(filterParams.MinPrice),
		valueOrNil(filterParams.MaxPrice),
		valueOrNil(filterParams.MinStarRating),
		valueOrNil(filterParams.MinReviewScore),
		valueOrNil(filterParams.MinReviews),
		valueOrNil(filterParams.Published),
		valueOrNil(filterParams.PropertyType),
		valueOrNil(filterParams.Feed),
		valueOrNil(filterParams.MinBedroom),
		valueOrNil(filterParams.Amenities),
		valueOrNil(filterParams.Limit),
	)

	data := services.GetAllProperties(filterParams)
	p.Data["json"] = data
	p.ServeJSON()
}

func valueOrNil[T any](v *T) any {
	if v == nil {
		return nil
	}

	return *v
}
