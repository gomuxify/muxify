package main

import (
	"fmt"
	"net/http"

	"github.com/gomuxify/muxify"
)

func main() {
	r := muxify.NewRouter()
	r.AllowTrailingSlash()
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Posting /")
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Getting /")
	})

	r.Get("/person/:name", func(w http.ResponseWriter, r *http.Request) {
		name := muxify.GetParam(r, "name")
		fmt.Fprintf(w, "Hello, %s\n", name)
	})

	http.ListenAndServe(":8080", r)
}
