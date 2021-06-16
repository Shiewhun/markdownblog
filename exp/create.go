package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type newArticleForm struct {
	ID          string    `schema:"ID"`
	DateTime    time.Time `schema:"datetime"`
	Title       string    `schema:"title"`
	Description string    `schema:"description"`
	Markdown    string    `schema:"markdown"`
}

type articleSchema struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title,omitempty" bson:"title,omitempty"`
	Description string             `json:"description,omitempty" bson:"description,omitempty"`
	Markdown    string             `json:"markdown,omitempty" bson:"markdown,omitempty"`
}

func create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// parse form body
	// obtain title, description, markdown
	// set article fields to title, description, markdown
	// fill in createdAt.
	// store in db
	//
	// if no errors, then redirect to
	// GET articles/id
	form := newArticleForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	article := articleSchema{
		Title:       form.Title,
		Description: form.Description,
		Markdown:    form.Markdown,
	}
	// need to find a way to use the mongo connection already set in main()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://dbUser:secret-password@cluster0.ehlxz.mongodb.net/admin?retryWrites=true&w=majority"))
	blogDatabase := client.Database("markdownblog")
	blogCollection := blogDatabase.Collection("blogposts")
	blogResult, err := blogCollection.InsertOne(ctx, article)
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Fprintf(w, "%v", blogResult.InsertedID)
	primitiveID := blogResult.InsertedID.(primitive.ObjectID) // type switch from intereface{} to pritive.ObjectID
	idStr := primitiveID.Hex()                                // convert from primitive.ObjectID to str so it can be put into url
	url := fmt.Sprintf("/articles/%s", idStr)
	http.Redirect(w, r, url, http.StatusFound)
}
