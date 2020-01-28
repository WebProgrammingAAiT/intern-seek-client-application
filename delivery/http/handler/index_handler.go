package handler

import (
	"html/template"
	"net/http"

	"github.com/abdimussa87/intern-seek-client-application/entity"
	"github.com/dgrijalva/jwt-go"
)

type IndexHandler struct {
	tmpl *template.Template
}

//NewIndexHandler initializes and returns new IndexHandler
func NewIndexHandler(T *template.Template) *IndexHandler {
	return &IndexHandler{tmpl: T}
}

func (indexHandler IndexHandler) GetIndex(w http.ResponseWriter, r *http.Request) {

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
				indexHandler.tmpl.ExecuteTemplate(w, "index.company.layout", nil)
			} else if claims.Role == "intern" {

				indexHandler.tmpl.ExecuteTemplate(w, "index.intern.layout", nil)
			}
		}

	} else {
		indexHandler.tmpl.ExecuteTemplate(w, "login.html", nil)
	}
}
