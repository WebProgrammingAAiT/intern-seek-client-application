package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/MahletH/Intern-Seek-Version-1/intern-seek-client-application/entity"

)

type CompanyProfileHandler struct {
	tmpl *template.Template
}

//NewCompanyProfileHandler initializes and returns new CompanyProfileHandler
func NewCompanyProfileHandler(T *template.Template) *CompanyProfileHandler {
	return &CompanyProfileHandler{tmpl: T}
}

func (cph CompanyProfileHandler) CompanyProfile(w http.ResponseWriter, r *http.Request) {
	compDetail := &entity.CompanyDetail{}
	User := entity.User{}
	company := entity.Company{}

	//getting userId
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

	User.ID = claims.UserID
	User.Name = r.FormValue("Name")
	User.Username = r.FormValue("Username")
	User.Phone = r.FormValue("Phone")
	User.Email = r.FormValue("Email")

	if r.Method == http.MethodPost {
		compDetail.UserID = claims.UserID
		compDetail.FocusArea = r.FormValue("FocusArea")
		compDetail.Description = r.FormValue("Description")
		compDetail.Country = r.FormValue("Country")
		compDetail.City = r.FormValue("City")
		id, err := strconv.Atoi((r.FormValue("compid")))
		if err == nil {
			compDetail.ID = uint(id)
		}

		if compDetail.ID == 0 {
			fmt.Printf("Company id equals %d", compDetail.ID)
			fmt.Printf("Please")
			fmt.Printf("Compdetail uid %d", compDetail.UserID)
			fmt.Printf("User id %d", User.ID)
			//Used to get userId from cookie

			url := "http://localhost:8181/v1/company"
			url2 := fmt.Sprintf("http://localhost:8181/v1/user/update/%s", strconv.Itoa(int(claims.UserID)))

			fmt.Println(compDetail.UserID)

			output, err := json.MarshalIndent(&compDetail, "", "\t\t")

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

			err = json.NewDecoder(resp.Body).Decode(compDetail)
			if err != nil {
				fmt.Printf(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			defer resp.Body.Close()
			fmt.Printf("User before sending id equal %d", int(User.ID))
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
			company.CompDetail = *compDetail
			company.CompUser = User
			cph.tmpl.ExecuteTemplate(w, "company.profile.layout", company)
		} else {
			fmt.Println("In else")
			//Used to get userId from cookie
			fmt.Printf("Compdetail id %d", int(compDetail.ID))
			url := fmt.Sprintf("http://localhost:8181/v1/company/update/%s", strconv.Itoa(int(compDetail.ID)))
			url2 := fmt.Sprintf("http://localhost:8181/v1/user/update/%s", strconv.Itoa(int(claims.UserID)))

			output, err := json.MarshalIndent(&compDetail, "", "\t\t")

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
			company.CompDetail = *compDetail
			company.CompUser = User
			cph.tmpl.ExecuteTemplate(w, "company.profile.layout", company)
		}

	} else if r.Method == http.MethodGet {

		url := fmt.Sprintf("http://localhost:8181/v1/users/%s", strconv.Itoa(int(User.ID)))

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
		url2 := fmt.Sprintf("http://localhost:8181/v1/companybyuserid/%s", strconv.Itoa(int(User.ID)))

		response2, err := client.Get(url2)
		if err != nil {
			company.CompUser = User
			company.CompDetail = *compDetail
			fmt.Println("got inside an eror")
			cph.tmpl.ExecuteTemplate(w, "company.profile.layout", company)
			return
		}
		defer response2.Body.Close()

		body, _ = ioutil.ReadAll(response2.Body)

		err = json.Unmarshal(body, compDetail)

		if err != nil {

			company.CompUser = User
			company.CompDetail = *compDetail
			cph.tmpl.ExecuteTemplate(w, "company.profile.layout", company)
			return
		}
		company.CompUser = User
		company.CompDetail = *compDetail

		cph.tmpl.ExecuteTemplate(w, "company.profile.layout", company)
	}
}
