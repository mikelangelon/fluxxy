package rest

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"encoding/json"
	"strconv"
	"log"
)

func V1Handler() *restful.Container {
	c := restful.NewContainer()
	c.Add(GameWebService)
	return c
}

var tailJSON = []byte{'\n'}
// ServeJSON writes the HTTP response body.
func ServeJSON(w http.ResponseWriter, statusCode int, src interface{}) {
	bytes, err := json.MarshalIndent(src, "", "\t")
	if err != nil {
		log.Print("goe rest: serialize response body: ", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h := w.Header()
	h.Set("Content-Type", "application/json;charset=UTF-8")
	h.Set("Content-Length", strconv.Itoa(len(bytes)+len(tailJSON)))
	w.WriteHeader(statusCode)

	if _, err := w.Write(bytes); err != nil {
		log.Print("goe rest: write response body: ", err)
	}
	if _, err := w.Write(tailJSON); err != nil {
		log.Print("goe rest: write response body: ", err)
	}
}