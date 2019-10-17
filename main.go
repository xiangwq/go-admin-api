package main

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"go-admin-api/models"
	_ "go-admin-api/routers"
)

func main() {
	models.Init()
	err := models.RegisterRbac()
	if err != nil {
		panic(err)
	}
	beego.Run()
}
