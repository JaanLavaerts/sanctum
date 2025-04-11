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
	e.POST("/logout", Logout)
	e.POST("/register", Register)

	// routes that need auth
	auth := e.Group("")
	auth.Use(AuthMiddleware())
	auth.GET("/vault", VaultPage)
	auth.POST("/add", AddEntry)
}

func AuthMiddleware() echo.MiddlewareFunc {
	middlewareAuthConfig := middleware.KeyAuthConfig{
		KeyLookup: "cookie:auth-token",
		Validator: func(token string, c echo.Context) (bool, error) {
					db_token, err := database.GetToken()
					if err != nil {
						return false, err
					}
					return crypto.VerifyToken(token, db_token), err
				},
	}
	return middleware.KeyAuthWithConfig(middlewareAuthConfig)
}
