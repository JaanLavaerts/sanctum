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
		log.Printf("Error getting master password: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	auth_token, err := database.GetToken()

	if err != nil {
		log.Printf("Error getting auth token: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	data := loginPageData{
		IsNew:      len(masterPassword) == 0,
		IsLoggedIn: len(auth_token) != 0,
	}

	return c.Render(http.StatusOK, "login", data)
}

func Login(c echo.Context) error {
	formMasterPassword := c.FormValue("masterpassword")

	ok, err := AuthenticateUser(formMasterPassword)
	if err != nil {
		log.Printf("Authentication error: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	if !ok {
		data := loginPageData{Error: "wrong master password"}
		return c.Render(http.StatusOK, "login", data)
	}

	raw_token, hashed_token, err := crypto.GenerateAuthToken()
	if err != nil {
		log.Printf("Error generating auth token: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	writeAuthCookie(c, raw_token)

	_, err = database.InsertToken(hashed_token)
	if err != nil {
		log.Printf("Error inserting token: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("HX-Redirect", "/vault")
	return c.NoContent(http.StatusOK)
}

func Register(c echo.Context) error {
	formMasterPassword := c.FormValue("masterpassword")

	salt, err := crypto.GenerateSalt()
	saltString := base64.RawURLEncoding.EncodeToString(salt)
	if err != nil {
		log.Printf("Error generating salt: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	hashed_password, err := crypto.GenerateHash(formMasterPassword)
	if err != nil {
		log.Printf("Error generating hash: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}
	_, err = database.InserMasterPassword(hashed_password, saltString)
	if err != nil {
		log.Printf("Error inserting master password: %v", err)
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

func AuthenticateUser(formMasterPassword string) (bool, error) {
	masterPassword, salt, err := database.GetMasterPassword()
	if err != nil {
		log.Printf("Error getting master password: %v", err)
		return false, err
	}

	saltString, err := base64.RawURLEncoding.DecodeString(salt)
	if err != nil {
		log.Printf("Error decoding salt: %v", err)
		return false, err
	}

	if !crypto.VerifyMasterPassword(formMasterPassword, masterPassword) {
		return false, nil
	}

	DerivedKey, err = crypto.DeriveKey(formMasterPassword, saltString)
	if err != nil {
		log.Printf("Error deriving key: %v", err)
		return false, err
	}
	return true, nil
}

func LogoutUser(c echo.Context, message string) error {
	clearAuthCookie(c)
	DerivedKey = nil

	_, err := database.DeleteToken()
	if err != nil {
		log.Printf("Error deleting token: %v", err)
		return c.NoContent(http.StatusInternalServerError)
	}

	data := loginPageData{
		IsNew:      false,
		IsLoggedIn: false,
		Error:      message,
	}

	return c.Render(http.StatusOK, "login", data)
}
