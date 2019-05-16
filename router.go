package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/asdine/storm"
	"github.com/julienschmidt/httprouter"
)

var StormDb *storm.DB
var err error

func init() {
	StormDb, err = storm.Open(".LoginDetail")
	if err != nil {
		log.Fatalln(err)
	}
	//todo also open sql database if need be
}

type User struct {
	Name     string
	Password string
	IsAdmin  bool
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := r.Cookie("logs")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
}

func LoginGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := r.Cookie("logs")
	if err == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	err = templates.ExecuteTemplate(w, "login.html", "")
	if err != nil {
		log.Fatalln(err)
	}
}

func (user *User) LoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	fmt.Println(err)
	user.Name = r.PostFormValue("User")
	user.Password = r.PostFormValue("Pass")
}
