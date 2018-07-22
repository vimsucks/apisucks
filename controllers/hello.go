package controllers

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/vimsucks/apisucks/models"
	"fmt"
)

func SayHello(c echo.Context) error {
	student := models.Student{}
	fmt.Println(student)
	return c.String(http.StatusOK, "Hello World")
}
