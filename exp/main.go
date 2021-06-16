package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var homeTemplate *template.Template
var client *mongo.Client

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	var err error
	// query for all articles in the db
	// render them.
	var data []newArticleForm
	collection := client.Database("markdownblog").Collection("blogposts")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// options.Find() is apparently used to set a Find
	// operation filter.
	findOptions := options.Find()
	// SetSort(bson.M{"_id": -1}) sorts the `primitiveID`
	// field in descending order.
	// Because a createdAt time stamp is the first 4 bytes
	// of primitiveID, then underneath code should sort in
	// descending order.
	findOptions.SetSort(bson.M{"_id": -1})
	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var article articleSchema
		var a newArticleForm
		cursor.Decode(&article)
		a.ID = article.ID.Hex()
		a.DateTime = article.ID.Timestamp().Local()
		a.Description = article.Description
		a.Title = article.Title
		a.Markdown = article.Markdown
		data = append(data, a)
	}
	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}

	homeTemplate, err = template.ParseFiles("../views/index.gohtml")
	if err != nil {
		log.Fatal(err)
	}

	if err := homeTemplate.Execute(w, data); err != nil {
		log.Fatal(err)
	}
}

func main() {
	// need to find a way to share the mongo db connection with
	// all my guys that need a db connection.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb+srv://dbUser:secret-password@cluster0.ehlxz.mongodb.net/admin?retryWrites=true&w=majority"))
	defer client.Disconnect(ctx)
	r := mux.NewRouter()
	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/articles/new", new).Methods("GET")
	r.HandleFunc("/articles", create).Methods("POST")
	r.HandleFunc("/articles/{id}/edit", edit).Methods("GET")
	r.HandleFunc("/articles/{id}/update", update).Methods("POST")
	r.HandleFunc("/articles/{id}", show).Methods("GET")
	r.HandleFunc("/articles/{id}/delete", deleteArticle).Methods("POST")
	fmt.Println("Starting application on port :4000")
	http.ListenAndServe(":4000", r)
}
