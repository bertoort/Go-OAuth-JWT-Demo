package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var templates = template.New("").Funcs(
	template.FuncMap{
		"Upper": func(s string) string {
			return strings.ToUpper(s)
		},
	},
)

func init() {
	filepath.Walk("views", func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".html") {
			templates.ParseFiles(path)
		}
		return nil
	})
}

// render a template given a model
func renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	err := templates.ExecuteTemplate(w, tmpl, p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("GET " + tmpl)
}

// sendJSON Will send out a json message to the client
func sendJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	msg, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	w.Write(msg)
}
