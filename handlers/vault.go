package handlers

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
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

func GeneratePassword(c echo.Context) error {
	charset := "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
	"!@#$%^&*()-_=+[]{}|;:,.<>?/"
	password := make([]byte, 12)
	for i := range password {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return err
		}
		password[i] = charset[randomIndex.Int64()]
	}
	html := fmt.Sprintf(`<input id="password" type="text" value="%s" name="password" required />`, password)

	return c.HTML(http.StatusOK, html)
}