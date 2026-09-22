package controllers

import (
	"rental-property-api/models"
	"rental-property-api/services"

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
	if id != "" {
		ob, err := services.GetPropertyByID(id)
		if err != nil {
			p.Ctx.Output.SetStatus(404)
			p.Data["json"] = models.ErrorResponse{Error: err.Error()}
		} else {
			p.Data["json"] = ob
		}
	}
	p.ServeJSON()
}
