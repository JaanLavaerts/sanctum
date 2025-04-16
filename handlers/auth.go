package handlers

import (
	"encoding/base64"
	"log"
	"net/http"
	"time"

	"github.com/JaanLavaerts/sanctum/crypto"
	"github.com/JaanLavaerts/sanctum/database"
	"github.com/labstack/echo/v4"
)

var DerivedKey []byte

type loginPageData struct {
	IsNew bool
	Error string
}

func LoginPage(c echo.Context) error {
	masterPassword, _, err := database.GetMasterPassword()
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

	masterPassword, salt, err := database.GetMasterPassword()
	if err != nil {
		log.Fatal(err)
	}

	saltString, err := base64.RawURLEncoding.DecodeString(salt)
	if err != nil {
		log.Fatal(err)
	}

	if !crypto.VerifyMasterPassword(formMasterPassword, masterPassword) {
		data := loginPageData{
			Error: "Wrong master password, please try again.",
		}
		return c.Render(http.StatusOK, "login", data)
	}

	raw_token, hashed_token := crypto.GenerateAuthToken()
	writeAuthCookie(c, raw_token)

	res, err := database.InsertToken(hashed_token)
	if err != nil {
		log.Fatal(err)
	}

	if res != 1 {
		return c.NoContent(http.StatusInternalServerError)
	}

	if err != nil {
		return err
	}

	DerivedKey, err = crypto.DeriveKey(formMasterPassword, saltString)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Redirect", "/vault")
	return c.NoContent(http.StatusOK)
}

func Register(c echo.Context) error {
	formMasterPassword := c.FormValue("masterpassword")

	salt, err := crypto.GenerateSalt()
	saltString := base64.RawURLEncoding.EncodeToString(salt)
	if err != nil {
		return err
	}

	hashed_password, err := crypto.GenerateHash(formMasterPassword)
	if err != nil {
		log.Fatalf("Error creating hash: %q", err)
	}
	res, err := database.InserMasterPassword(hashed_password, saltString)
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

	for i := range DerivedKey {
		DerivedKey[i] = 0
	}

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
	cookie.Expires = time.Now().Add(30 * time.Minute)
	cookie.HttpOnly = true
	c.SetCookie(cookie)
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
