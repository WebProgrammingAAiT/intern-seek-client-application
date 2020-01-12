package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/abdimussa87/intern-seek-client-application/entity"
)

type SignUpHandler struct {
	tmpl *template.Template
}

// NewSignUpHandler initializes and returns new SignUpHandler
func NewSignUpHandler(T *template.Template) *SignUpHandler {
	return &SignUpHandler{tmpl: T}
}

func (suh SignUpHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	if _, err := r.Cookie("token"); err == nil {
		http.Redirect(w, r, "/company", http.StatusSeeOther)
	}

	if r.Method == http.MethodPost {
		url := "http://localhost:8181/v1/signup"
		user := entity.User{}
		user.Username = r.FormValue("username")

		user.Password = r.FormValue("password")

		user.Email = r.FormValue("email")

		user.Phone = r.FormValue("phone")

		user.Name = r.FormValue("fullname")

		if user.Username == "" || user.Password == "" || user.Email == "" || user.Phone == "" || user.Name == "" {
			eror := Error{Name: "Please enter fields correctly"}
			suh.tmpl.ExecuteTemplate(w, "signup.html", eror)

		}

		retypedPassword := r.FormValue("retype-password")

		if retypedPassword != user.Password {
			eror := Error{Name: "Retype Password correctly"}
			suh.tmpl.ExecuteTemplate(w, "signup.html", eror)

		}

		output, err := json.MarshalIndent(user, "", "\t\t")

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

		for _, cookie := range resp.Cookies() {

			http.SetCookie(w, cookie)
		}

		defer resp.Body.Close()
		fmt.Println(resp.StatusCode)
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {

			http.Redirect(w, r, "/company", http.StatusSeeOther)

		} else {

			eror := Error{Name: "The Username already exists"}
			suh.tmpl.ExecuteTemplate(w, "signup.html", eror)

		}

	} else if r.Method == http.MethodGet {
		suh.tmpl.ExecuteTemplate(w, "signup.html", nil)
	}

}
