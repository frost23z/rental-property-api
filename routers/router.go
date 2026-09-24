// @APIVersion 1.0.0
// @Title Rental Property API
// @Description This is a sample server for managing rental properties.
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"rental-property-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/properties",
			beego.NSInclude(
				&controllers.PropertyController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
