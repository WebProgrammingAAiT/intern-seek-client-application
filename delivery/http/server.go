package main

import (
	"html/template"
	"net/http"

	"github.com/MahletH/Intern-Seek-Version-1/intern-seek-client-application/delivery/http/handler"
	"github.com/MahletH/Intern-Seek-Version-1/intern-seek-client-application/entity"

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

func indexHandler(w http.ResponseWriter, r *http.Request) {

	temp.ExecuteTemplate(w, "index.layout", nil)
}

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	temp.ExecuteTemplate(w, "login.html", nil)
// }

// func signupHandler(w http.ResponseWriter, r *http.Request) {
// 	temp.ExecuteTemplate(w, "signup.html", nil)
// }

// func companyManageHandler(w http.ResponseWriter, r *http.Request) {

// 	temp.ExecuteTemplate(w, "company.manage.job.layout", nil)

// }
func companyNewInternshipHandler(w http.ResponseWriter, r *http.Request) {

	temp.ExecuteTemplate(w, "company.post.new.internship.layout", nil)

}
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
func internProfileHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "intern.profile.layout", nil)
}
func main() {

	signInHandler := handler.NewSignInHandler(temp)
	signUpHandler := handler.NewSignUpHandler(temp)
	companyProfileHandler := handler.NewCompanyProfileHandler(temp)
	companyNewInternshipHandler := handler.NewIntenshipHandler(temp)
	companyManageHandler := handler.NewIntenshipHandler(temp)

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("../../ui/assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", signInHandler.SignIn)
	mux.HandleFunc("/signup", signUpHandler.SignUp)
	mux.HandleFunc("/company/manage", companyManageHandler.RetrieveInternship)
	mux.HandleFunc("/company", companyProfileHandler.CompanyProfile)
	mux.HandleFunc("/company/new-internship", companyNewInternshipHandler.AddInternship)

	mux.HandleFunc("/intern", internProfileHandler)
	mux.HandleFunc("/intern/applied", internAppliedHandler)
	mux.HandleFunc("/internship/desc", internDescHandler)
	http.ListenAndServe(":8080", mux)
}
