package main

import (
	"fmt"
	"net/http"
	"rest"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fluxx service!")
}

func main() {
	v1 := rest.V1Handler()
	http.Handle("/v1/", v1)

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
