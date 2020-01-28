package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
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
		compDetail := &entity.CompanyDetail{}
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

		client := &http.Client{}
		url2 := fmt.Sprintf("http://localhost:8181/v1/companybyuserid/%s", strconv.Itoa(int(userID)))

		response2, err := client.Get(url2)
		if err != nil {
			fmt.Println("No company found")
			return
		}
		defer response2.Body.Close()

		body, _ := ioutil.ReadAll(response2.Body)

		err = json.Unmarshal(body, compDetail)

		if err != nil {

			fmt.Printf("Unable to parse")
			return
		}
		compID := compDetail.ID

		fmt.Println(compID)
		url := "http://localhost:8181/v1/internship"
		intern := entity.Internship{}

		intern.Name = r.FormValue("name")
		date := r.FormValue("deadline")
		fmt.Println(date)
		intern.ClosingDate, _ = time.Parse("00/00/0000", date)
		intern.Description = r.FormValue("description")
		intern.RequiredAcademicLevel = r.FormValue("requiredAcademic")

		intern.CompanyID = compID

		if intern.Description == "" || intern.RequiredAcademicLevel == "" || intern.Name == "" {
			eror := Error{Name: "Please enter fields correctly"}
			ih.tmpl.ExecuteTemplate(w, "company.post.new.internship.layout", eror)

		}

		fields := r.Form["Field"]
		s := 0

		for range fields {
			s++
		}

		intern.FieldsReq = make([]entity.Field, s) ////

		k := 0
		for _, fi := range fields {
			splittedString := strings.Split(fi, "-")
			fieldName := splittedString[0]
			fieldId := splittedString[1]
			fmt.Println(fieldName)
			fmt.Println(fieldId)
			intern.FieldsReq[k].Name = fieldName
			j, _ := strconv.Atoi(fieldId)
			intern.FieldsReq[k].ID = uint(j)
			k++

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
		AvailableField := []entity.Field{}

		url3 := "http://localhost:8181/v1/field"

		req3, err := http.NewRequest("GET", url3, nil)
		if err != nil {
			panic(err)
		}

		req3.Header.Set("Content-Type", "application/json")
		client3 := &http.Client{}
		response3, err := client3.Do(req3)
		if err != nil {
			panic(err)
		}
		defer response3.Body.Close()
		body3, _ := ioutil.ReadAll(response3.Body)

		err = json.Unmarshal(body3, &AvailableField)
		fmt.Println(AvailableField)
		ih.tmpl.ExecuteTemplate(w, "company.post.new.internship.layout", AvailableField)
	}

}
func (ih InternshipHandler) RetrieveInternship(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
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
		fmt.Print(interns)

		ih.tmpl.ExecuteTemplate(w, "company.manage.job.layout", interns)

	}

}
func (ih InternshipHandler) CompanyRetrieveInternship(w http.ResponseWriter, r *http.Request) {

	interns := []entity.Internship{}
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
	url := fmt.Sprintf("http://localhost:8181/v1/companyInternship/%s/internships", strconv.Itoa(int(userID)))

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

	fmt.Print(interns)
	ih.tmpl.ExecuteTemplate(w, "company.manage.job.layout", interns)

}
