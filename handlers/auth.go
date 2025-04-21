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
	IsNew      bool
	Error      string
	Success    string
	IsLoggedIn bool
}

func LoginPage(c echo.Context) error {
	masterPassword, _, err := database.GetMasterPassword()
	if err != nil {
		log.Fatal(err)
	}
	auth_token, err := database.GetToken()

	if err != nil {
		log.Fatal(err)
	}

	data := loginPageData{
		IsNew:      len(masterPassword) == 0,
		IsLoggedIn: len(auth_token) != 0,
	}

	return c.Render(http.StatusOK, "login", data)
}

func Login(c echo.Context) error {
	formMasterPassword := c.FormValue("masterpassword")

	if !AuthenticateUser(formMasterPassword) {
		data := loginPageData{
			Error: "wrong master password",
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

	data := loginPageData{
		IsNew:      false,
		IsLoggedIn: false,
		Success:    "vault created",
	}

	return c.Render(http.StatusOK, "login", data)
}

func Logout(c echo.Context) error {
	return LogoutUser(c, "")
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

func AuthenticateUser(formMasterPassword string) bool {
	masterPassword, salt, err := database.GetMasterPassword()
	if err != nil {
		log.Fatal(err)
	}

	saltString, err := base64.RawURLEncoding.DecodeString(salt)
	if err != nil {
		log.Fatal(err)
	}

	if !crypto.VerifyMasterPassword(formMasterPassword, masterPassword) {
		return false
	}

	DerivedKey, err = crypto.DeriveKey(formMasterPassword, saltString)
	if err != nil {
		log.Fatal(err)
	}
	return true
}

func LogoutUser(c echo.Context, message string) error {
	clearAuthCookie(c)
	DerivedKey = nil

	_, err := database.DeleteToken()
	if err != nil {
		return err
	}

	data := loginPageData{
		IsNew:      false,
		IsLoggedIn: false,
		Error:      message,
	}

	return c.Render(http.StatusOK, "login", data)
}
