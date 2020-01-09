package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"

	_ "github.com/lib/pq"
	"github.com/nebyubeyene/Intern-Seek-Version-1/delivery/http/handler"
	"github.com/nebyubeyene/Intern-Seek-Version-1/user/repository"
	userRep "github.com/nebyubeyene/Intern-Seek-Version-1/user/repository"
	"github.com/nebyubeyene/Intern-Seek-Version-1/user/service"
	userServ "github.com/nebyubeyene/Intern-Seek-Version-1/user/service"
)

var tmpl = template.Must(template.ParseGlob("../../ui/templates/*"))

func indexHandler(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "index.html", nil)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func main() {

	dbconn, err := sql.Open("postgres", "user=postgres dbname=interndb password='P@$$w0rDd' sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		panic(err)
	}

	userRepo := userRep.NewUserRepositoryImpl(dbconn)
	userServi := userServ.NewUserServiceImpl(userRepo)

	compRepo := repository.NewCompanyRepositoryImpl(dbconn)
	compServ := service.NewCompanyService(compRepo)

	//userHandler := handler.NewUserHandler( userServi)

	compHandler := handler.NewCompanyHandler(compServ, userServi)

	router := httprouter.New()

	router.GET("/v1/company", compHandler.GetCompanies)
	router.POST("/v1/company/:id", compHandler.PostCompany)
	router.PUT("/v1/company/update/:id", compHandler.PutCompany)
	router.DELETE("/v1/company/delete/:id", compHandler.DeleteCompany)

	http.ListenAndServe(":8181", router)

	// mux := http.NewServeMux()
	// fs := http.FileServer(http.Dir("../../ui/assets"))
	// mux.Handle("/assets/", http.StripPrefix("/assets", fs))
	// mux.HandleFunc("/", indexHandler)
	// mux.HandleFunc("/login", loginHandler)
	// mux.HandleFunc("/signup", userHandler.SignUp)
	// http.ListenAndServe(":8080", mux)
}
