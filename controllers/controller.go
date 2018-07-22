package controllers

import "github.com/labstack/echo"

func RegisterControllers(e *echo.Echo) {
	e.GET("/hello", SayHello)
	e.GET("/progressbar/:percentage", ProgressBar)
	e.POST("/cust/:sid/score/all", GetAllScore)
	e.POST("/cust/:sid/score/year/:year", GetYearScore)
	e.POST("/cust/:sid/score/year/:year/semester/:semester", GetYearSemesterScore)
}

