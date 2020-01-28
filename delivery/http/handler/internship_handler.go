package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/abdimussa87/intern-seek-client-application/entity"
	"github.com/dgrijalva/jwt-go"
)

type InternshipHandler struct {
	tmpl *template.Template
}

// NewInternshipHandler initializes and returns new InternshipHandler
func NewIntenshipHandler(T *template.Template) *InternshipHandler {
	return &InternshipHandler{tmpl: T}
}

// AddInternship adds internship
func (ih InternshipHandler) AddInternship(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		c, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println("Login to post")
			return
		}
		tknStr := c.Value
		claims := &entity.Claims{}
		_, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID := uint(claims.UserID)

		fmt.Println(userID)
		url := "http://localhost:8181/v1/internship"
		intern := entity.Internship{}

		intern.Name = r.FormValue("name")
		date := r.FormValue("deadline")
		fmt.Println(date)
		intern.ClosingDate, _ = time.Parse("00/00/0000", date)
		intern.Description = r.FormValue("description")
		intern.RequiredAcademicLevel = r.FormValue("requiredAcademic")
		intern.CompanyID = userID

		if intern.Description == "" || intern.RequiredAcademicLevel == "" || intern.Name == "" {
			eror := Error{Name: "Please enter fields correctly"}
			ih.tmpl.ExecuteTemplate(w, "company.post.new.internship.layout", eror)

		}

		output, err := json.MarshalIndent(intern, "", "\t\t")

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(output))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
		fmt.Println(resp.StatusCode)
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {

			http.Redirect(w, r, "/company", http.StatusSeeOther)

		}

	} else if r.Method == http.MethodGet {
		ih.tmpl.ExecuteTemplate(w, "company.post.new.internship.layout", nil)
	}

}
func (ih InternshipHandler) RetrieveInternship(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		c, err := r.Cookie("token")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		tknStr := c.Value
		claims := &entity.Claims{}
		_, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userID := uint(claims.UserID)

		fmt.Println(userID)
		url := "http://localhost:8181/v1/internships"
		interns := []entity.Internship{}

		client := &http.Client{}
		response, err := client.Get(url)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(body, &interns)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		ih.tmpl.ExecuteTemplate(w, "company.manage.job.layout", interns)

	}

}
