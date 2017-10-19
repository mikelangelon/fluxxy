package main

import (
	"fmt"
	"net/http"
	"rest"
	"game"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Fluxx service!")
}

func main() {
	game.Games = make(map[int64]*game.Game)

	v1 := rest.V1Handler()
	http.Handle("/v1/", v1)

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
