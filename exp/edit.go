package main

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func edit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// get article from db by id
	//
	// parse url and extract id
	// query db with id and extract single article.
	// if successful
	// render edit page.
	// prefill html with values for .Description .Title .Markdown
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
	var data newArticleForm
	data.ID = a.ID.Hex() // converts primitive.ObjectID to string
	data.Title = a.Title
	data.Description = a.Description
	data.Markdown = a.Markdown
	editTemplate, err := template.ParseFiles("../views/edit.gohtml")
	if err != nil {
		panic(err)
	}

	if err := editTemplate.Execute(w, data); err != nil {
		panic(err)
	}
}
