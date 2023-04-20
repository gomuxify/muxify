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

	r.Get("/organisations/:organisation/people/:people", func(w http.ResponseWriter, r *http.Request) {
		params := muxify.Params(r)
		fmt.Fprintf(w, "Organisation: %s\n", params["organisation"])
		fmt.Fprintf(w, "Person: %s\n", params["people"])
	})

	http.ListenAndServe(":8080", r)
}
