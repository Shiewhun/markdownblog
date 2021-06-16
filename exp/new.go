package main

import (
	"html/template"
	"net/http"
)

func new(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	newArticleTemplate, err := template.ParseFiles("../views/new.gohtml")
	if err != nil {
		panic(err)
	}

	if err := newArticleTemplate.Execute(w, nil); err != nil {
		panic(err)
	}
}
