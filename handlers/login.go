package handlers

import (
	"log"
	"net/http"

	"github.com/JaanLavaerts/sanctum/crypto"
	"github.com/JaanLavaerts/sanctum/database"
	"github.com/labstack/echo/v4"
)


type loginPageData struct {
	IsNew bool
}

func LoginPage(c echo.Context) error {
	masterPassword, err := database.GetMasterPassword()
	if err != nil {
		log.Fatal(err)
	}

	data := loginPageData{
		IsNew: len(masterPassword) == 0,
	}

	return c.Render(http.StatusOK, "login", data)
}

func Login(c echo.Context) error {
	formMasterPassword := c.FormValue("masterpassword")
	masterPassword, err := database.GetMasterPassword()
	if err != nil {
		log.Fatal(err)
	}

	if crypto.VerifyMasterPassword(formMasterPassword, masterPassword) { 
		c.Response().Header().Set("HX-Redirect", "/vault")
	}

	return c.NoContent(http.StatusUnauthorized)
}

func Register(c echo.Context) error {
	formMasterPassword := c.FormValue("masterpassword")

	res, err := database.InserMasterPassword(formMasterPassword)
	if err != nil {
		log.Fatal(err)
	}

	if res == 1 {
		c.Response().Header().Set("HX-Redirect", "/")
	}

	return c.NoContent(http.StatusOK)
}
