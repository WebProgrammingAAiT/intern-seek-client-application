package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/abdimussa87/intern-seek-client-application/entity"
	"github.com/dgrijalva/jwt-go"
)

type InternProfileHandler struct {
	tmpl *template.Template
}

//NewInternProfileHandler initializes and returns new InternProfileHandler
func NewInternProfileHandler(T *template.Template) *InternProfileHandler {
	return &InternProfileHandler{tmpl: T}
}

func (internph InternProfileHandler) InternProfile(w http.ResponseWriter, r *http.Request) {
	internDetail := &entity.PersonalDetails{}
	User := entity.User{}
	intern := entity.Intern{}
	//fieldsOfIntern := make([]entity.Field, 100)

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

	// url3 := "http://localhost:8181/v1/field"

	// req3, err := http.NewRequest("GET", url3, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// req3.Header.Set("Content-Type", "application/json")
	// client3 := &http.Client{}
	// response3, err := client3.Do(req3)
	// if err != nil {
	// 	panic(err)
	// }
	// defer response3.Body.Close()
	// body3, _ := ioutil.ReadAll(response3.Body)

	// err = json.Unmarshal(body3, &intern.AvailableField)

	if r.Method == http.MethodPost {
		r.ParseForm()
		internDetail.UserID = claims.UserID
		//internDetail.Fields = r.Form["Field"]
		internDetail.AcademicLevel = r.FormValue("AcademicLevel")
		fields := r.Form["Field"]

		//fmt.Println(fieldsOfIntern)
		//internDetail.Fields = fieldsOfIntern
		//fieldIds := r.Form["fieldid"]
		//for i,fid := range fieldIds{
		//	j,_ := strconv.Atoi(fid)
		//	fmt.Println(j)
		//	fmt.Println(i)
		//	//internDetail.Fields[i].ID=uint(j)
		//}

		id, err := strconv.Atoi((r.FormValue("internid")))
		if err == nil {
			internDetail.ID = (uint(id))
		}

		if internDetail.ID == 0 {
			fmt.Printf("Intern id equals %d ", internDetail.ID)
			fmt.Println("Please")
			fmt.Printf("Interndetail uid %d ", internDetail.UserID)
			fmt.Printf("User id %d ", User.ID)
			//Used to get userId from cookie

			url := "http://localhost:8181/v1/intern"
			url2 := fmt.Sprintf("http://localhost:8181/v1/user/update/%s", strconv.Itoa(int(claims.UserID)))

			s := 0

			for range fields {
				s++
			}

			internDetail.Fields = make([]entity.Field, s) ////

			k := 0
			for _, fi := range fields {
				splittedString := strings.Split(fi, "-")
				fieldName := splittedString[0]
				fieldId := splittedString[1]
				fmt.Println(fieldName)
				fmt.Println(fieldId)
				internDetail.Fields[k].Name = fieldName
				j, _ := strconv.Atoi(fieldId)
				internDetail.Fields[k].ID = uint(j)
				k++

			}
			// fmt.Println(internDetail.UserID)
			// fmt.Println(internDetail.AcademicLevel)
			// fmt.Println(internDetail.Fields)

			output, err := json.MarshalIndent(&internDetail, "", "\t\t")
			//fmt.Println(output)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			fmt.Printf("here1")
			req, err := http.NewRequest("POST", url, bytes.NewBuffer(output))
			if err != nil {
				fmt.Println("here2")
				panic(err)
			}
			fmt.Println("here3")
			req.Header.Set("Content-Type", "application/json")
			//req.Header.Set("token", tknStr)
			client := &http.Client{}
			resp, err := client.Do(req)
			fmt.Println("here4")
			if err != nil {
				panic(err)
			}
			fmt.Println("here5")
			err = json.NewDecoder(resp.Body).Decode(internDetail)
			fmt.Println("here6")
			if err != nil {
				fmt.Printf(err.Error())
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			fmt.Printf("here7")
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
			url3 := "http://localhost:8181/v1/field"

			req3, err := http.NewRequest("GET", url3, nil)
			if err != nil {
				panic(err)
			}

			req3.Header.Set("Content-Type", "application/json")
			client3 := &http.Client{}
			response3, err := client3.Do(req3)
			if err != nil {
				panic(err)
			}
			defer response3.Body.Close()
			body3, _ := ioutil.ReadAll(response3.Body)

			err = json.Unmarshal(body3, &intern.AvailableField)
			//fmt.Println(intern.AvailableField[0])
			if err != nil {
				fmt.Println(err)
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			intern.InternDetail = *internDetail
			intern.InternUser = User
			internph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
		} else {
			//fmt.Println("In else")
			//Used to get userId from cookie
			//fmt.Printf("Compdetail id %d", int(internDetail.ID))
			url := fmt.Sprintf("http://localhost:8181/v1/intern/update/%s", strconv.Itoa(int(internDetail.ID)))
			url2 := fmt.Sprintf("http://localhost:8181/v1/user/update/%s", strconv.Itoa(int(claims.UserID)))

			s := 0

			for range fields {
				s++
			}

			f := make([]entity.Field, s)
			internDetail.Fields = f ////

			k := 0
			for _, fi := range fields {
				splittedString := strings.Split(fi, "-")
				fieldName := splittedString[0]
				fieldId := splittedString[1]
				fmt.Println(fieldName)
				fmt.Println(fieldId)
				internDetail.Fields[k].Name = fieldName
				j, _ := strconv.Atoi(fieldId)
				internDetail.Fields[k].ID = uint(j)
				k++

			}
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
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			defer resp.Body.Close()
			body, _ := ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(body, &internDetail)
			//fmt.Println(intern.AvailableField[0])
			if err != nil {
				fmt.Println(err)
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			intern.InternDetail = *internDetail

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
			url3 := "http://localhost:8181/v1/field"

			req3, err := http.NewRequest("GET", url3, nil)
			if err != nil {
				panic(err)
			}

			req3.Header.Set("Content-Type", "application/json")
			client3 := &http.Client{}
			response3, err := client3.Do(req3)
			if err != nil {
				panic(err)
			}
			defer response3.Body.Close()
			body3, _ := ioutil.ReadAll(response3.Body)

			err = json.Unmarshal(body3, &intern.AvailableField)
			//fmt.Println(intern.AvailableField[0])
			if err != nil {
				fmt.Println(err)
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			intern.InternDetail = *internDetail
			intern.InternUser = User
			http.Redirect(w, r, "/intern", http.StatusSeeOther)
			// internph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
		}

	} else if r.Method == http.MethodGet {

		url := fmt.Sprintf("http://localhost:8181/v1/users/%s", strconv.Itoa(int(User.ID)))

		url3 := "http://localhost:8181/v1/field"

		req3, err := http.NewRequest("GET", url3, nil)
		if err != nil {
			panic(err)
		}

		req3.Header.Set("Content-Type", "application/json")
		client3 := &http.Client{}
		response3, err := client3.Do(req3)
		if err != nil {
			panic(err)
		}
		defer response3.Body.Close()
		body3, _ := ioutil.ReadAll(response3.Body)

		err = json.Unmarshal(body3, &intern.AvailableField)
		//internDetail.Fields = make([]entity.Field, 5)
		//fmt.Println(intern.AvailableField[0])

		// req3, err := http.NewRequest("GET", url3, nil)
		// if err != nil {
		// 	panic(err)
		// }

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
		url2 := fmt.Sprintf("http://localhost:8181/v1/internbyuser/%s", strconv.Itoa(int(User.ID)))

		response2, err := client.Get(url2)
		if err != nil {
			intern.InternUser = User
			intern.InternDetail = *internDetail
			fmt.Println("got inside an eror")
			internph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
			return
		}
		defer response2.Body.Close()

		body, _ = ioutil.ReadAll(response2.Body)

		err = json.Unmarshal(body, internDetail)

		if err != nil {

			intern.InternUser = User
			intern.InternDetail = *internDetail
			fmt.Println("got inside an eror")
			internph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
			return
		}
		intern.InternUser = User
		intern.InternDetail = *internDetail

		internph.tmpl.ExecuteTemplate(w, "intern.profile.layout", intern)
	}
}
