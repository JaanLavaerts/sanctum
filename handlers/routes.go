package handlers

import (
	"github.com/labstack/echo/v4"
)


func RegisterRoutes(e *echo.Echo) {
	e.GET("/", LoginPage)
	e.GET("/vault", VaultPage)

	e.POST("/login", Login)
	e.POST("/register", Register)
	e.POST("/add", AddEntry)
}
