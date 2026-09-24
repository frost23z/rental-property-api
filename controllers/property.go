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
// @Description get all properties
// @Success 200 {object} models.Response
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
