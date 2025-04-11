package handlers

import (
	"log"
	"net/http"
	"time"

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

	if !crypto.VerifyMasterPassword(formMasterPassword, masterPassword) { 
		return c.NoContent(http.StatusUnauthorized)
	}

	raw_token, hashed_token := crypto.GenerateToken()
	writeAuthCookie(c, raw_token)

	res, err :=	database.InsertToken(hashed_token)
	if err != nil {
		log.Fatal(err)
	}

	if res != 1 {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("HX-Redirect", "/vault")
	return c.NoContent(http.StatusOK)
}

func Register(c echo.Context) error {
	formMasterPassword := c.FormValue("masterpassword")

	res, err := database.InserMasterPassword(formMasterPassword)
	if err != nil {
		log.Fatal(err)
	}

	if res != 1 {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
}

func Logout(c echo.Context) error {
	clearAuthCookie(c)
	res, err := database.DeleteToken()
	if err != nil {
		return err
	}

	if res != 1 {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("HX-Redirect", "/")
	return c.NoContent(http.StatusOK)
} 

func writeAuthCookie(c echo.Context, raw_token string) {
	cookie := new(http.Cookie)
	cookie.Name = "auth-token"
	cookie.Value = raw_token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}

func readAuthCookie(c echo.Context) (*http.Cookie, error) {
	cookie, err := c.Cookie("auth-token")
	if err != nil {
		return nil, err
	}
	return cookie, nil
}

func clearAuthCookie(c echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = "auth-token"
	cookie.Value = "" 
	cookie.Expires = time.Unix(0, 0)
	cookie.MaxAge = -1 
	cookie.HttpOnly = true
	c.SetCookie(cookie)
}