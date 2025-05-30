package handlers

import (
	"github.com/JaanLavaerts/sanctum/crypto"
	"github.com/JaanLavaerts/sanctum/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegisterRoutes(e *echo.Echo) {
	// public routes
	e.GET("/", LoginPage)
	e.POST("/login", Login)
	e.POST("/register", Register)
	e.POST("/logout", Logout)

	// routes that need auth
	auth := e.Group("")
	auth.Use(AuthMiddleware())
	auth.GET("/vault", VaultPage)
	auth.POST("/add", AddEntry)
	auth.DELETE("/delete/:id", DeleteEntry)
	auth.GET("/generate", GeneratePassword)
	auth.GET("/reveal/:id", RevealPassword)
	auth.GET("/hide/:id", HidePassword)
}

func AuthMiddleware() echo.MiddlewareFunc {
	middlewareAuthConfig := middleware.KeyAuthConfig{
		KeyLookup: "cookie:auth-token",
		Validator: func(token string, c echo.Context) (bool, error) {
			db_token, err := database.GetToken()
			if err != nil {
				return false, err
			}
			isValid, verifyErr := crypto.VerifyAuthToken(token, db_token)
			if verifyErr != nil {
				return false, verifyErr
			}
			return isValid, nil
		},
		ErrorHandler: func(err error, c echo.Context) error {
			return LogoutUser(c, "session expired")
		},
	}
	return middleware.KeyAuthWithConfig(middlewareAuthConfig)
}
