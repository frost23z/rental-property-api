package main

import (
	"os"

	_ "rental-property-api/routers"
	"rental-property-api/services"

	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	dataFile, err := beego.AppConfig.String("datafile")
	if err != nil || dataFile == "" {
		logs.Error("datafile is not set in conf/app.conf: %v", err)
		os.Exit(1)
	}
	if err := services.LoadData(dataFile); err != nil {
		logs.Error("failed to load source data: %v", err)
		os.Exit(1)
	}

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
