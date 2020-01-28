package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/abdimussa87/intern-seek-client-application/entity"
	"github.com/dgrijalva/jwt-go"
)

type SignInHandler struct {
	tmpl *template.Template
}
type Error struct {
	Name string
}

// NewSignInHandler initializes and returns new SignInHandler
func NewSignInHandler(T *template.Template) *SignInHandler {
	return &SignInHandler{tmpl: T}
}

func (sih SignInHandler) SignIn(w http.ResponseWriter, r *http.Request) {

	if _, err := r.Cookie("token"); err == nil {
		c, _ := r.Cookie("token")
		tknStr := c.Value
		claims := &entity.Claims{}
		_, err = jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if claims.Role == "company" {
			http.Redirect(w, r, "/company", http.StatusSeeOther)
		} else if claims.Role == "intern" {
			http.Redirect(w, r, "/intern", http.StatusSeeOther)
		}
	}

	if r.Method == http.MethodPost {
		url := "http://localhost:8181/v1/signin"
		user := entity.User{}
		user.Username = r.FormValue("username")

		user.Password = r.FormValue("password")
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

		defer resp.Body.Close()
		//fmt.Println(resp.StatusCode)

		if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
			claims := &entity.Claims{}
			for _, cookie := range resp.Cookies() {
				fmt.Println("Cookie with name", cookie.Name)
				http.SetCookie(w, cookie)
				//trying to get the role type from the cookie using claims
				_, err = jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte("secret"), nil
				})

				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
			//Checking user role type
			userRole := &entity.UserRole{}

			userRole.Role = claims.Role

			if userRole.Role == "company" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
		} else {
			eror := Error{Name: "Invalid username or password"}
			sih.tmpl.ExecuteTemplate(w, "login.html", &eror)
		}

	} else if r.Method == http.MethodGet {
		sih.tmpl.ExecuteTemplate(w, "login.html", nil)
	}
}
