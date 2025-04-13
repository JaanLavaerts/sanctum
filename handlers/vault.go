package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/JaanLavaerts/sanctum/crypto"
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
	password := c.FormValue("password")
	site := c.FormValue("site")
	notes := c.FormValue("notes")
	timestamp := time.Now()

	encryptedPassword, nonce, err := crypto.EncryptEntryPassword(password, DerivedKey)
	stringNonce := base64.StdEncoding.EncodeToString(nonce)
	stringPassword := base64.StdEncoding.EncodeToString(encryptedPassword)

	newEntry := database.Entry{
		Password: stringPassword,	
		Site: site,
		Notes: notes,
		Timestamp: timestamp,
		Nonce: stringNonce,
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

func RevealPassword(c echo.Context) error {
	id := c.Param("id")
	entry, err := database.GetEntry(id)
	fmt.Println(entry)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("STUFF", entry, DerivedKey)
	plainPassword, _ := crypto.DecryptPassword(entry.Password, DerivedKey, entry.Nonce)

	return c.String(200, plainPassword)
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