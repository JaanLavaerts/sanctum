package handler

import (
	"log"
	"net/http"
	"text/template"

	"github.com/JaanLavaerts/sanctum/crypto"
	"github.com/JaanLavaerts/sanctum/database"
)

func LoginPage(w http.ResponseWriter, r *http.Request) {
	plain_password, err := database.GetMasterPassword()
	var tmpl *template.Template
	if err != nil {
		log.Fatal(err)
	}
	if len(plain_password) == 0 {
		tmpl = template.Must(template.ParseFiles(
			"templates/base.html",
			"templates/master-password/create.html",
		))
	} else {
		tmpl = template.Must(template.ParseFiles(
			"templates/base.html",
			"templates/master-password/login.html",
		))
	}
	tmpl.ExecuteTemplate(w, "base", nil)
}

func CreateMasterPassword(w http.ResponseWriter, r *http.Request) {
	var plain_password string
	if r.Method == http.MethodPost {
		r.ParseForm()
		plain_password = r.FormValue("master_password")
		database.InserMasterPassword(plain_password)

		w.Header().Set("HX-Redirect", "/vault")
		w.WriteHeader(http.StatusOK)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	hashed_password, err := database.GetMasterPassword()
	if err != nil {
		log.Fatal(err)
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		plain_password := r.FormValue("master_password")

		if !crypto.VerifyMasterPassword(plain_password, hashed_password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.Header().Set("HX-Redirect", "/vault")
		w.WriteHeader(http.StatusOK)
	}
}