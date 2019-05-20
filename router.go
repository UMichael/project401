package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

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
type CourseDetails struct {
	ID         int `storm:"increment"`
	Name       string
	CourseCode string
	Password   string
}
type LecturerDetails struct {
	ID         int `storm:"increment"`
	Password   string
	Name       string
	CourseCode string `storm:"unique"`
	//Students   []StudentDetails `storm:"inline"`
}

var StudentDb, studentOfferingCourses, lecturerDB *storm.DB

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
	StudentDb, err = storm.Open(".StudentDetails")
	var data []Database
	StormDb.All(&data)
	fmt.Println(data)
	lecturerDB, err = storm.Open(".lecturerDb")
	studentOfferingCourses, err = storm.Open(".courses")
	//courseDetails, err = storm.Open(".courseDetails")

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
	i := 1
	name := "name"
	for {
		var student StudentDetails
		name = r.FormValue("name_" + strconv.Itoa(i))
		if name == "" {
			break
		}
		student.Name = name
		student.Matric = r.FormValue("matric_" + strconv.Itoa(i))
		student.Level = r.FormValue("level_" + strconv.Itoa(i))
		student.FingersChecked = r.FormValue("finger_" + strconv.Itoa(i))
		err := StudentDb.Save(&student)
		if err != nil {
			fmt.Println("error saving students", err)
		}
		fmt.Println(student)
		i++
	}
	http.Redirect(w, r, "/", 302)
}
func (user *User) Students(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var students []StudentDetails
	err := StudentDb.AllByIndex("Level", &students)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(students)
	templates = templates.Lookup("table.html")
	err = templates.Execute(w, students)
	if err != nil {
		log.Fatalln(err)
	}
}

func (user *User) AddCoursePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	if user.Details.IsAdmin != true {
		http.Redirect(w, r, "/login", 302)
		return
	}
	var course LecturerDetails
	course.Password = r.FormValue("Pass")
	course.Name = r.FormValue("name")
	course.CourseCode = r.FormValue("course")
	fmt.Println(studentOfferingCourses, lecturerDB.Save(&course))
	//todo fix error if user has registered
	http.Redirect(w, r, "/", 302)
}
func (user *User) AddCourse(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if user.Details.IsAdmin != true {
		http.Redirect(w, r, "/login", 302)
	}
	templates = templates.Lookup("AddCourse.html")
	err := templates.Execute(w, user)
	if err != nil {
		log.Fatalln(err)
	}
}

func (user *User) LecturerLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()
	var lecturer LecturerDetails
	Code := r.FormValue("code")
	Password := r.FormValue("Pass")
	fmt.Println(Code, Password)
	err := lecturerDB.One("Password", Password, &lecturer)
	if err != nil || lecturer.CourseCode != Code {
		user.Login.BadDetail = true
		fmt.Println(err)
		http.Redirect(w, r, "/lecturer", 302)
		return
	}
	http.SetCookie(
		w,
		&http.Cookie{
			Path: "/",
			Name: "lecturer",
		})
	http.Redirect(w, r, "/", 302)
}

func (user *User) LecturerLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	templates = templates.Lookup("lecturerLogin.html")
	err := templates.Execute(w, user)
	if err != nil {
		log.Fatalln(err)
	}
}
