package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/asdine/storm"
	"github.com/julienschmidt/httprouter"
)

var StormDb *storm.DB
var err error

type User struct {
	Details Database
	Login   LoginDetails
}

type LoginDetails struct {
	BadDetail bool
	Success   bool
}

type Database struct {
	ID       int    `storm:"increment"`
	Name     string `storm:"unique"`
	Password string
	IsAdmin  bool
}

type StudentDetails struct {
	ID             int `storm:"increment"`
	Name           string
	Matric         string `storm:"unique"`
	Level          string `storm:"index"`
	FingersChecked string
}

func init() {
	StormDb, err = storm.Open(".LoginDetail")
	if err != nil {
		log.Fatalln(err)
	}
	StormDb.Save(&Database{
		IsAdmin:  true,
		Name:     "admin",
		Password: "admin",
	})
	var data []Database
	StormDb.All(&data)
	fmt.Println(data)

	//todo also open sql database if need be
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func (user *User) HomePage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := r.Cookie("data")
	if err != nil {
		http.Redirect(w, r, "/login", 302)
		return
	}
	templates = templates.Lookup("home.html")
	err = templates.Execute(w, user)
	if err != nil {
		log.Fatalln(err)
	}

}

func LoginGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := r.Cookie("logs")
	if err == nil {
		http.Redirect(w, r, "/", 302)
		return
	}
	var data Database
	templates = templates.Lookup("login.html")
	templates.Execute(w, data)
}

func (user *User) LoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	err := r.ParseForm()
	fmt.Println(err)
	user.Details.Name = r.PostFormValue("User")
	user.Details.Password = r.PostFormValue("Pass")
	fmt.Println(user)
	var NewData Database
	err = StormDb.One("Name", user.Details.Name, &NewData)
	if err != nil {
		data := LoginDetails{
			BadDetail: true,
		}
		templates = templates.Lookup("login.html")
		templates.Execute(w, &data)
		return
	}
	cookies := http.Cookie{
		Name:     "data",
		Value:    user.Details.Name,
		Path:     "/",
		HttpOnly: true,
	}
	user.Details = NewData
	http.SetCookie(w, &cookies)
	http.Redirect(w, r, "/", 302)
}

func addAdmin(r *http.Request) error {
	r.ParseForm()

	Admin := User{
		Details: Database{
			IsAdmin:  true,
			Name:     r.FormValue("Name"),
			Password: r.FormValue("Pass"),
		},
	}
	return StormDb.Save(&Admin)
}

func (user *User) ChangePassword(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	name := r.FormValue("name")
	oldPass := r.FormValue("oldPass")
	newPass := r.FormValue("newPass")
	var userDB Database
	if err := StormDb.One("Name", name, &userDB); err != nil {
		user.Login.BadDetail = true
		fmt.Println(err)
		http.Redirect(w, r, "/changePass", 302)
		return
	}
	fmt.Println(user)
	if userDB.Password != oldPass {
		user.Login.BadDetail = true
		http.Redirect(w, r, "/login", 302)
		return
	}
	user.Details.Password = newPass
	user.Login.Success = true
	fmt.Println(user)
	StormDb.Save(&user)
	http.Redirect(w, r, "/changePass", 302)
}

func (user *User) ChangePasswordGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	templates = templates.Lookup("changeDetails.html")
	err := templates.Execute(w, user)
	if err != nil {
		log.Fatalln(err)
	}
	user.Login.Success = false
	user.Login.BadDetail = false
}

func (user *User) AdminEnroll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	data, err := goquery.NewDocumentFromReader(r.Body)
	if err != nil {
		log.Fatalln(err)
	}
	data.Find("tbody#tablebody").Each(func(i int, tbodyHtml *goquery.Selection) {
		tbodyHtml.Find("tr").Each(func(j int, trHtml *goquery.Selection) {
			var student StudentDetails
			trHtml.Find("td").Each(func(k int, tdHtml *goquery.Selection) {
				if k == 0 {
					student.Matric = tdHtml.Text()
				} else if k == 1 {
					student.Name = tdHtml.Text()
				} else if k == 2 {
					student.Level = tdHtml.Text()
				} else if k == 3 {
					student.FingersChecked = tdHtml.Text()
				}
			})
			StormDb.Save(&student)
			fmt.Println(student)
		})
	})
	http.Redirect(w, r, "/", 302)
}
