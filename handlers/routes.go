package handlers

import (
	"net/http"

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

	// routes that need auth
	auth := e.Group("")
	auth.Use(AuthMiddleware())
	auth.POST("/logout", Logout)
	auth.GET("/vault", VaultPage)
	auth.POST("/add", AddEntry)
	auth.DELETE("/delete/:id", DeleteEntry)
	auth.GET("/generate", GeneratePassword)
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
		ErrorHandler: func(err error, c echo.Context) error {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		},
	}
	return middleware.KeyAuthWithConfig(middlewareAuthConfig)
}
