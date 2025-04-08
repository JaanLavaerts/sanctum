package main

import (
	"log"
	"net/http"

	"github.com/JaanLavaerts/sanctum/database"
	"github.com/JaanLavaerts/sanctum/internal/handler"
)

func main() {
	database.InitDB()
	// pages
	http.HandleFunc("/", handler.LoginPage)
	http.HandleFunc("/vault", handler.VaultPage)

	// functions
	http.HandleFunc("/create-master-password", handler.CreateMasterPassword)
	http.HandleFunc("/login", handler.Login)

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
