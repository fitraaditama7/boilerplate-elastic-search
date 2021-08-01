package main

import (
	"boilerplate-elastic-search/internal/pkg/storage/elasticsearch"
	"boilerplate-elastic-search/internal/post"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	elastic, err := elasticsearch.New([]string{"http://0.0.0.0:9200"})
	if err != nil {
		log.Fatalln(err)
	}

	if err := elastic.CreateIndex("post"); err != nil {
		log.Fatalln(err)
	}

	storage, err := elasticsearch.NewPostStorage(*elastic)
	if err != nil {
		log.Fatalln(err)
	}

	postAPI := post.New(storage)

	router := httprouter.New()
	router.HandlerFunc("POST", "/api/v1/posts", postAPI.Create)
	router.HandlerFunc("PATCH", "/api/v1/posts/:id", postAPI.Update)
	router.HandlerFunc("DELETE", "/api/v1/posts/:id", postAPI.Delete)
	router.HandlerFunc("GET", "/api/v1/posts/:id", postAPI.Find)
	router.HandlerFunc("GET", "/api/v1/search/:keyword", postAPI.Search)

	log.Fatalln(http.ListenAndServe(":4000", router))
}
