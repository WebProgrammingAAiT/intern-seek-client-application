package handler

import (
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {

	expiredCookie := &http.Cookie{Name: "token", MaxAge: -1, Expires: time.Now().Add(-100 * time.Hour)}

	http.SetCookie(w, expiredCookie)

	http.Redirect(w, r, "/login", http.StatusSeeOther)

}
