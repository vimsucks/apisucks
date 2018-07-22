package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	. "github.com/vimsucks/apisucks/config"
	"github.com/vimsucks/apisucks/controllers"
	"github.com/vimsucks/apisucks/models"
)

func main() {
	defer models.DB.Close()
	fmt.Printf("Config: %+v\n", Conf)
	e := echo.New()
	if Conf.HttpsEnabled {
		log.Info("HTTPS enabled, redirect http -> https")
		e.Pre(middleware.HTTPSRedirect())
	}
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	controllers.RegisterControllers(e)
	//e.POST("/cust/student/:id/score", controllers.GetStudentScore)

	if Conf.HttpsEnabled {
		e.Logger.Fatal(e.StartTLS(fmt.Sprintf("%s:%d", Conf.Host, Conf.Port),
			Conf.CertPath,
			Conf.KeyPath))
	} else {
		e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%d", Conf.Host, Conf.Port)))
	}
}
