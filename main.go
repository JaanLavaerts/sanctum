package main

import (
	"html/template"
	"io"

	"github.com/JaanLavaerts/sanctum/database"
	"github.com/JaanLavaerts/sanctum/handlers"
	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	database.InitDB()
	e := echo.New()
	
	t := &Template{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
	e.Renderer = t

	handlers.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(":8080"))
}
