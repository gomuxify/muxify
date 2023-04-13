package main

import (
	"fmt"
	"net/http"

	"github.com/gomuxify/muxify"
)

func main() {
	r := muxify.NewRouter()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Posting /")
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Getting /")
	})

	http.ListenAndServe(":8080", r)
}
