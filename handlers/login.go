package handlers

import (
	"log"
	"net/http"

	"github.com/JaanLavaerts/sanctum/database"
	"github.com/labstack/echo/v4"
)


type PageData struct {
	MasterPassword string
	IsNew bool
}

func LoginPage(c echo.Context) error {
	master_password, err := database.GetMasterPassword()
	if err != nil {
		log.Fatal(err)
	}

	data := PageData{
		MasterPassword: master_password,
		IsNew: len(master_password) == 0,
	}

	return c.Render(http.StatusOK, "login.html", data)
}

