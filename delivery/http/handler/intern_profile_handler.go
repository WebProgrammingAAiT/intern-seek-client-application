package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"github.com/lensabillion/Project/entity"
	"github.com/dgrijalva/jwt-go"
)

type InternProfileHandler struct{
	tmpl *template.Template
}

//NewInternProfileHandler initializes and returns new InternProfileHandler

func NewInternProfileHandler(T *template.Template) *InternProfileHandler{
	return &InternProfileHandler{tmpl:T}
}

func (iph InternProfileHandler) InternProfile(w http.ResponseWriter,r *http.Request){
	internDetail := &entity.PersonalDetails{}
	User := entity.User{}
	intern := entity.Intern{}

	if r.Method == http.MethodPost{
		if internDetail.ID == 0{
			fmt.Printf("Intern id equals %d",internDetail.ID)
			//Used to get userId from cookie
			c,err := r.Cookie("token")
			if err != nil{
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			tknStr :=c.Value
			claims := &entity.Claims{}
			_,err=jwt.ParseWithClaims(tknStr,claims,func(token *jwt.Token)(interface{},error){
				return []byte("secret"),nil
			})
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			url := "http://localhost:8181/v1/intern"
			url2 := fmt.Sprintf("http://localhost:8181/v1/user/update/%s", strconv.Itoa(int(claims.UserID)))

			internDetail.UserID =claims.UserID
			internDetail.AcademicLevel=r.FormValue("AcademicLevel")
			internDetail.Field=r.FormValue("Field")

			User.Name = r.FormValue("Name")
			User.Username = r.FormValue("Username")
			User.Phone = r.FormValue("Phone")
			User.Email = r.FormValue("Email")
			fmt.Println(internDetail.UserID)
			output, err := json.MarshalIndent(&internDetail, "", "\t\t")

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
			_, err = client.Do(req)
			if err != nil {
				panic(err)
			}

			output2, err := json.MarshalIndent(&User, "", "\t\t")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			req2, err := http.NewRequest("PUT", url2, bytes.NewBuffer(output2))
			if err != nil {
				panic(err)
			}

			req2.Header.Set("Content-Type", "application/json")
			client2 := &http.Client{}
			_, err = client2.Do(req2)
			if err != nil {
				panic(err)
			}
			intern.InternDetail = *internDetail
			intern.InternUser = User
			iph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
		} else {
			fmt.Println("In else")
			//Used to get userId from cookie
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
			url := fmt.Sprintf("http://localhost:8181/v1/intern/update/%s", strconv.Itoa(int(internDetail.ID)))
			url2 := fmt.Sprintf("http://localhost:8181/v1/user/update/%s", strconv.Itoa(int(claims.UserID)))
			User.Name = r.FormValue("Name")
			User.ID = claims.UserID
			User.Username = r.FormValue("Username")
			User.Phone = r.FormValue("Phone")
			User.Email = r.FormValue("Email")

			internDetail.AcademicLevel = r.FormValue("AcademicLevel")
			internDetail.Field = r.FormValue("Field")


			output, err := json.MarshalIndent(&internDetail, "", "\t\t")

			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			req, err := http.NewRequest("PUT", url, bytes.NewBuffer(output))
			if err != nil {
				panic(err)
			}

			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{}
			_, err = client.Do(req)
			if err != nil {
				panic(err)
			}

			output2, err := json.MarshalIndent(&User, "", "\t\t")
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			req2, err := http.NewRequest("PUT", url2, bytes.NewBuffer(output2))
			if err != nil {
				panic(err)
			}

			req2.Header.Set("Content-Type", "application/json")
			client2 := &http.Client{}
			_, err = client2.Do(req2)
			if err != nil {
				panic(err)
			}
			intern.InternDetail = *internDetail
			intern.InternUser=User
			iph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
		}

	} else if r.Method == http.MethodGet {
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
		userID := strconv.Itoa(int(claims.UserID))

		fmt.Println(userID)
		url := fmt.Sprintf("http://localhost:8181/v1/users/%s", userID)

		client := &http.Client{}
		response, err := client.Get(url)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)

		err = json.Unmarshal(body, &User)
		fmt.Println(User.Name)
		if err != nil {
			fmt.Println(err)
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		url2 := fmt.Sprintf("http://localhost:8181/v1/internbyuserid/%s", userID)

		response2, err := client.Get(url2)
		if err != nil {
			intern.InternUser = User
			intern.InternDetail = *internDetail
			fmt.Println("got inside an error")
			iph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
			return
		}
		defer response2.Body.Close()

		body, _ = ioutil.ReadAll(response2.Body)

		err = json.Unmarshal(body, internDetail)

		if err != nil {
			// fmt.Println(err)
			// w.Header().Set("Content-Type", "application/json")
			// http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			intern.InternUser = User
			intern.InternDetail= *internDetail
			fmt.Println("got inside an error")
			iph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
			return
		}
		intern.InternUser = User
		intern.InternDetail = *internDetail

		iph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
	}
}
