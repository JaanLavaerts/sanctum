package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/JaanLavaerts/sanctum/database"
	"github.com/labstack/echo/v4"
)

type PageData struct {
	HashedPassword string
}

func RegisterRoutes(e *echo.Echo) {
	e.GET("/", indexHandler)
	e.GET("/time", timeHandler)
}

func indexHandler(c echo.Context) error {
	hashedPassword, err := database.GetMasterPassword()
	if err != nil {
		log.Fatal(err)
	}

	data := PageData{
		HashedPassword: hashedPassword,
	}

	return c.Render(http.StatusOK, "layout.html", data)
}

func timeHandler(c echo.Context) error {
	currentTime := time.Now().Format("15:04:05")
	return c.String(http.StatusOK, currentTime)
}