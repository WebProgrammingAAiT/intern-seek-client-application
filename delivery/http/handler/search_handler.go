package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/abdimussa87/intern-seek-client-application/entity"
)

type SearchHandler struct {
	tmpl *template.Template
}

//NewSearchHandler initializes and returns new SearchHandler
func NewSearchHandler(T *template.Template) *SearchHandler {
	return &SearchHandler{tmpl: T}
}

func (searchHandler SearchHandler) Search(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		fmt.Println("INside search handler")
		name := r.FormValue("field")
		field := entity.Field{}
		field.Name = name
		url := fmt.Sprintf("http://localhost:8181/v1/field/%s/internship", name)
		output, err := json.MarshalIndent(field, "", "\t\t")

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		req, err := http.NewRequest("GET", url, bytes.NewBuffer(output))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {

			eror := Error{Name: "Search not found"}
			searchHandler.tmpl.ExecuteTemplate(w, "index.intern.layout", eror)
			return

		}

		defer resp.Body.Close()
		internships := []entity.Internship{}
		err = json.NewDecoder(resp.Body).Decode(&internships)
		fmt.Println("INside search handler")
		if err != nil {
			fmt.Printf(err.Error())

			searchHandler.tmpl.ExecuteTemplate(w, "index.intern.layout", nil)
			fmt.Println("IN error")
			return

		}

		searchHandler.tmpl.ExecuteTemplate(w, "index.intern.layout", internships)
	}
}
