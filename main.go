package main

import (
	"html/template"
	"net/http"

	"github.com/Beemnet/internseek/handler"
	"github.com/Beemnet/internseek/internship/repository"
	"github.com/Beemnet/internseek/internship/service"
	"github.com/julienschmidt/httprouter"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var tmpl = template.Must(template.ParseGlob("C:/Users/123/go/src/github.com/Beemnet/internseek/ui/templates/*"))

func index(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "internship.new.layout", nil)
}

func main() {

	dbconn, err := gorm.Open("postgres", "user=postgres password=CaputDraconis dbname=interndb sslmode=disable")
	if err != nil {
		panic(err)
	}

	defer dbconn.Close()
	/*
		errs := dbconn.CreateTable(&entity.Comment{}, &entity.Role{}).GetErrors()
		if errs := nil{
			panic(err)
		}
	*/
	internshipRepo := repository.NewInternshipGormRepo(dbconn)
	internshipSrv := service.NewInternshipService(internshipRepo)
	internshipHandler := handler.NewInternshipHandler(internshipSrv)

	router := httprouter.New()

	fs := http.FileServer(http.Dir("ui/assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", fs))
	http.HandleFunc("/", index)
	router.GET("/internship/", internshipHandler.GetInternships)
	router.POST("/internship/", internshipHandler.PostInternship)

	http.ListenAndServe("localhost:8181", router)

}
