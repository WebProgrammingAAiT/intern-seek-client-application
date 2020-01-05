package handler

import (
	"html/template"
	"net/http"

	"github.com/abdimussa87/Intern-Seek-Version-1/entity"
	"github.com/abdimussa87/Intern-Seek-Version-1/user"
)

type UserHandler struct {
	t        *template.Template
	userServ user.UserService
}

func NewUserHandler(T *template.Template, US user.UserService) *UserHandler {
	return &UserHandler{t: T, userServ: US}
}

//SignUp handles requests coming at /signup
func (uh UserHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		user := entity.User{}
		user.Fullname = r.FormValue("fullname")
		user.UUID = r.FormValue("username")
		user.Email = r.FormValue("email")
		user.Phone = r.FormValue("phone")
		user.Password = r.FormValue("password")

		err := uh.userServ.StoreUser(user)
		if err != nil {
			panic(err)
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else if r.Method == http.MethodGet {
		uh.t.ExecuteTemplate(w, "signup.html", nil)
	}
}
