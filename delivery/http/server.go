package main

import (
	"html/template"
	"net/http"
)

var temp = template.Must(template.ParseGlob("../../ui/templates/*"))

func indexHandler(w http.ResponseWriter, r *http.Request) {

	temp.ExecuteTemplate(w, "index.html", nil)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "login.html", nil)
}

func signupHandler(w http.ResponseWriter,r *http.Request){
	temp.ExecuteTemplate(w,"signup.html",nil)
}

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("../../ui/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/signup",signupHandler)
	http.ListenAndServe(":8080", mux)
}
