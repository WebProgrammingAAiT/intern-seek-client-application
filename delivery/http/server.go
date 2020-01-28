package main

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/abdimussa87/intern-seek-client-application/delivery/http/handler"
	"github.com/abdimussa87/intern-seek-client-application/entity"
	"github.com/dgrijalva/jwt-go"
)

var temp = template.Must(template.ParseGlob("../../ui/templates/*"))

func internDescHandler(w http.ResponseWriter, r *http.Request) {

	temp.ExecuteTemplate(w, "internship.desc.layout", nil)

}

type Claims struct {
	Username string
	jwt.StandardClaims
}

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	temp.ExecuteTemplate(w, "login.html", nil)
// }

// func signupHandler(w http.ResponseWriter, r *http.Request) {
// 	temp.ExecuteTemplate(w, "signup.html", nil)
// }

func companyPostHandler(w http.ResponseWriter, r *http.Request) {
	compDetail := entity.CompanyDetail{City: "ADDis", Country: "Ethiopia", Description: "This a good company", FocusArea: "Software"}
	compUser := entity.User{
		Username: "ADf",
		Email:    "A@g.com",
		Name:     "Habene",
		Phone:    "0912545658",
	}
	company := entity.Company{
		CompDetail: compDetail,
		CompUser:   compUser,
	}
	temp.ExecuteTemplate(w, "company.profile.layout", company)
}

func internAppliedHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "intern.applied.layout", nil)
}

func main() {

	searchHandler := handler.NewSearchHandler(temp)
	signInHandler := handler.NewSignInHandler(temp)
	signUpHandler := handler.NewSignUpHandler(temp)
	companyProfileHandler := handler.NewCompanyProfileHandler(temp)
	companyNewInternshipHandler := handler.NewIntenshipHandler(temp)
	internProfileHandler := handler.NewInternProfileHandler(temp)
	indexHandler := handler.NewIndexHandler(temp)
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("../../ui/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler.GetIndex)
	mux.HandleFunc("/login", signInHandler.SignIn)
	mux.HandleFunc("/signup", signUpHandler.SignUp)
	mux.HandleFunc("/company/manage", companyNewInternshipHandler.CompanyRetrieveInternship)
	mux.HandleFunc("/company/new-internship", companyNewInternshipHandler.AddInternship)
	mux.Handle("/company", isAuthorizedCompany(companyProfileHandler.CompanyProfile))
	mux.Handle("/intern", isAuthorizedIntern(internProfileHandler.InternProfile))
	mux.HandleFunc("/intern/applied", internAppliedHandler)
	mux.HandleFunc("/internship/desc", internDescHandler)
	mux.HandleFunc("/search", searchHandler.Search)
	mux.HandleFunc("/logout", handler.Logout)
	http.ListenAndServe(":8080", mux)
}

//Middleware for checking authorization for viewing a page
func isAuthorizedCompany(endpoint func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if cookie, err := r.Cookie("token"); err == nil {
			claims := &entity.Claims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if token.Valid {
				if claims.Role == "company" {
					endpoint(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode("Unauthorized")
					return

				}
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Unauthorized")
			return
		}
	})
}

//Middleware for checking authorization for viewing a page
func isAuthorizedIntern(endpoint func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if cookie, err := r.Cookie("token"); err == nil {
			claims := &entity.Claims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte("secret"), nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if token.Valid {
				if claims.Role == "intern" {
					endpoint(w, r)
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					json.NewEncoder(w).Encode("Unauthorized")
					return

				}
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode("Unauthorized")
			return
		}
	})
}
