package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func show(w http.ResponseWriter, r *http.Request) {
	// parse id from request
	// extract article by id in db
	// set article values in an article object variable
	// render article
	w.Header().Set("Content-Type", "text/html")
	vars := mux.Vars(r)
	id := vars["id"]
	primitiveID, _ := primitive.ObjectIDFromHex(id)
	blogCollection := client.Database("markdownblog").Collection("blogposts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var a articleSchema
	if err := blogCollection.FindOne(ctx, bson.M{"_id": primitiveID}).Decode(&a); err != nil {
		log.Fatal(err)
	}
	type data struct {
		ID          string
		Title       string
		DateTime    time.Time
		Description string
		Markdown    template.HTML
	}
	var d data
	d.ID = a.ID.Hex() // converts primitive.ObjectID to string
	d.Title = a.Title
	d.DateTime = a.ID.Timestamp().Local()
	d.Description = a.Description
	contentBytes := []byte(a.Markdown)
	htmlBytes := markdown.ToHTML(contentBytes, nil, nil)
	html := template.HTML(htmlBytes)
	d.Markdown = html
	showTemplate, err := template.ParseFiles("../views/show.gohtml")
	if err != nil {
		panic(err)
	}

	if err := showTemplate.Execute(w, d); err != nil {
		panic(err)
	}
}
