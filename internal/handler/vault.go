package handler

import (
	"net/http"
	"text/template"
)

func VaultPage(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/vault/vault.html",
	))
	tmpl.ExecuteTemplate(w, "base", nil)
}
