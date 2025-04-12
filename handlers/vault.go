package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/JaanLavaerts/sanctum/database"
	"github.com/labstack/echo/v4"
)

type vaultPageData struct{
	Entries []database.Entry
}

func VaultPage(c echo.Context) error {
	entries, err := database.GetEntries()
	if err != nil {
		log.Fatal(err)
	}

	data := vaultPageData{
		Entries: entries,
	} 

	return c.Render(http.StatusOK, "vault", data)
}

func AddEntry(c echo.Context) error {
	// TODO cant delete entry right after adding because ID is not set locally, only in DB
	password := c.FormValue("password")
	site := c.FormValue("site")
	notes := c.FormValue("notes")
	timestamp := time.Now()

	newEntry := database.Entry{
		Password: password,	
		Site: site,
		Notes: notes,
		Timestamp: timestamp,
	}
	
	id, err := database.InsertEntry(newEntry)
	if err != nil {
		log.Fatal(err)
	}

	newEntry.Id = id

	return c.Render(http.StatusOK, "entry", newEntry)
}

func DeleteEntry(c echo.Context) error {
	id := c.Param("id")
	err := database.DeleteEntry(id)
	if err != nil {
		log.Fatal(err)
	}

	return c.NoContent(http.StatusOK)
}