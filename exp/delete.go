package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

func deleteArticle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "text/html")
	// parse id
	// query db for document with id
	// delete document
	vars := mux.Vars(r)
	id := vars["id"]
	primitiveID, _ := primitive.ObjectIDFromHex(id)
	usersCollection := client.Database("markdownblog").Collection("blogposts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := usersCollection.DeleteOne(ctx, bson.M{"_id": primitiveID})
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/", http.StatusFound)
}
