package main

import (
	"database/sql"
	"html/template"
	"net/http"

	"github.com/abdimussa87/Intern-Seek-Version-1/delivery/http/handler"
	userRep "github.com/abdimussa87/Intern-Seek-Version-1/user/repository"
	userServ "github.com/abdimussa87/Intern-Seek-Version-1/user/service"
	_ "github.com/lib/pq"
)

var tmpl = template.Must(template.ParseGlob("../../ui/templates/*"))

func indexHandler(w http.ResponseWriter, r *http.Request) {

	tmpl.ExecuteTemplate(w, "index.html", nil)

}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.html", nil)
}

func main() {

	dbconn, err := sql.Open("postgres", "user=app_admin dbname=interndb password='P@$$w0rdD2' sslmode=disable")

	if err != nil {
		panic(err)
	}

	defer dbconn.Close()

	if err := dbconn.Ping(); err != nil {
		panic(err)
	}

	userRepo := userRep.NewUserRepositoryImpl(dbconn)
	userServi := userServ.NewUserServiceImpl(userRepo)

	userHandler := handler.NewUserHandler(tmpl, userServi)

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("../../ui/assets"))
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/signup", userHandler.SignUp)
	http.ListenAndServe(":8080", mux)
}
