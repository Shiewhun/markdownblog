package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func update(w http.ResponseWriter, r *http.Request) {
	// parse request/form into a form object
	// use mongo Update function to update article by id.
	w.Header().Set("Content-Type", "text/html")
	form := newArticleForm{}
	if err := parseForm(r, &form); err != nil {
		panic(err)
	}
	// need to find a way to use the mongo connection already set in main()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://dbUser:secret-password@cluster0.ehlxz.mongodb.net/admin?retryWrites=true&w=majority"))
	blogDatabase := client.Database("markdownblog")
	blogCollection := blogDatabase.Collection("blogposts")
	update := bson.M{
		"$set": bson.M{
			"title":       form.Title,
			"description": form.Description,
			"markdown":    form.Markdown,
		},
	}
	vars := mux.Vars(r)
	idStr := vars["id"]
	primitiveID, _ := primitive.ObjectIDFromHex(idStr)
	_, err := blogCollection.UpdateOne(
		ctx,
		bson.M{"_id": primitiveID},
		update,
	)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("/articles/%v", idStr)
	http.Redirect(w, r, url, http.StatusFound)
}
