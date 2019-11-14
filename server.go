package main

import (
	"html/template"
	"net/http"
)

var temp = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {

	temp.Execute(w, nil)

}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", mux)
}
