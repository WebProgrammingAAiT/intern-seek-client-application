package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/abdimussa87/intern-seek-client-application/entity"
	"github.com/dgrijalva/jwt-go"
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
		url := "http://localhost:8181/v1/signup"
		url2 := "http://localhost:8181/v1/userrole"
		user := entity.User{}
		userRole := entity.UserRole{}
		user.Username = r.FormValue("username")

		user.Password = r.FormValue("password")

		user.Email = r.FormValue("email")

		user.Phone = r.FormValue("phone")

		user.Name = r.FormValue("fullname")
		userRole.Role = r.FormValue("role")

		if user.Username == "" || user.Password == "" || user.Email == "" || user.Phone == "" || user.Name == "" || userRole.Role == "" {
			eror := Error{Name: "Please enter fields correctly"}
			suh.tmpl.ExecuteTemplate(w, "signup.html", eror)
			return

		}

		retypedPassword := r.FormValue("retype-password")

		if retypedPassword != user.Password {
			eror := Error{Name: "Retype Password correctly"}
			suh.tmpl.ExecuteTemplate(w, "signup.html", eror)
			return
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

		//getting userId from response
		responseMap := map[string]interface{}{}
		err = json.NewDecoder(resp.Body).Decode(&responseMap)

		if err != nil {
			json.NewEncoder(w).Encode("Invalid response")
			return
		}
		// header = strings.TrimSpace(header)
		// fmt.Println(header)
		// if header == "" {
		// 	//Token is missing
		// 	fmt.Println(header)
		// 	fmt.Println("Insid eror")
		// 	w.WriteHeader(http.StatusForbidden)
		// 	json.NewEncoder(w).Encode("Missing auth token")
		// 	return
		// }
		tokenStr := fmt.Sprint(responseMap["token"])
		claims := &entity.Claims{}
		_, err = jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		userRole.UserId = claims.UserID
		//storing user role
		output2, err := json.MarshalIndent(userRole, "", "\t\t")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		req2, err := http.NewRequest("POST", url2, bytes.NewBuffer(output2))
		if err != nil {
			panic(err)
		}
		req2.Header.Set("Content-Type", "application/json")
		client2 := &http.Client{}
		resp2, err := client2.Do(req2)
		if err != nil {
			panic(err)
		}

		err = json.NewDecoder(resp2.Body).Decode(&userRole)
		if err != nil {
			fmt.Printf(err.Error())
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		defer resp2.Body.Close()

		expirationTime := time.Now().Add(15 * time.Minute)
		// Create the JWT claims, which includes the username and expiry time
		claims2 := &entity.Claims{
			UserID: userRole.UserId,
			Role:   userRole.Role,
			Name:   user.Name,
			StandardClaims: jwt.StandardClaims{

				ExpiresAt: expirationTime.Unix(),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims2)

		tokenString, err := token.SignedString([]byte("secret"))
		if err != nil {

			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		// for _, cookie := range resp.Cookies() {

		// 	http.SetCookie(w, cookie)
		// }

		// defer resp.Body.Close()
		// fmt.Println(resp.StatusCode)
		if resp.StatusCode >= 200 && resp.StatusCode <= 299 && resp2.StatusCode >= 200 && resp2.StatusCode <= 299 {
			if userRole.Role == "company" {
				http.Redirect(w, r, "/company", http.StatusSeeOther)
			} else {
				http.Redirect(w, r, "/intern", http.StatusSeeOther)
			}
		} else {

			eror := Error{Name: "The Username already exists"}
			suh.tmpl.ExecuteTemplate(w, "signup.html", eror)

		}

	} else if r.Method == http.MethodGet {
		suh.tmpl.ExecuteTemplate(w, "signup.html", nil)
	}

}
