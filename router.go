package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/asdine/storm"
	"github.com/tealeg/xlsx"

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
	Courses        string
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
	if err != nil {
		log.Fatalln(err)
	}
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
	_, err := r.Cookie("admin")
	if err != nil {
		_, err = r.Cookie("lecturer")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
	}

	templates = templates.Lookup("home.html")
	err = templates.Execute(w, user)
	if err != nil {
		log.Fatalln(err)
	}

}

func (user *User) LoginGet(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	if user.Details.Name != "admin" {
		http.SetCookie(
			w,
			&http.Cookie{
				Path: "/",
				Name: "lecturer",
			})
	} else {
		http.SetCookie(
			w,
			&http.Cookie{
				Path: "/",
				Name: "admin",
			})
	}
	user.Details = NewData
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
	_, err := r.Cookie("admin")
	if err != nil {
		_, err = r.Cookie("lecturer")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
	}
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
	if user.Details.IsAdmin != true {
		http.Redirect(w, r, "/login", 302)
	}
	_, err := r.Cookie("admin")
	if err != nil {
		_, err = r.Cookie("lecturer")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
	}
	r.ParseMultipartForm((1 << 10) * 24)
	level := r.FormValue("level")
	i := 1
	name := "name"
	file, _, err := r.FormFile("file")
	switch err {
	case nil:
		defer file.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, file); err != nil {
			return
		}
		xlFile, _ := xlsx.OpenBinary(buf.Bytes())
		var found bool
		for _, sheet := range xlFile.Sheets {
			for _, row := range sheet.Rows {
				var student StudentDetails
				for _, cell := range row.Cells {
					text := cell.String()
					if text != "1" && found != true {

					} else if len(text) > 0 {

						if text[:1] == "1" && len(text) == 9 {
							student.Matric = text
							student.Level = level
						} else if _, err := strconv.Atoi(text); err != nil {
							student.Name = text

						}
						//fmt.Println(text)
						found = true
					}
				}
				if student.Name == "" {
					continue
				} else {
					StudentDb.Save(&student)
					fmt.Println(student)
				}
			}
		}

	case http.ErrMissingFile:
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
	default:
		log.Println(err)
	}
	http.Redirect(w, r, "/", 302)
}
func (user *User) Students(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := r.Cookie("admin")
	if err != nil {
		_, err = r.Cookie("lecturer")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
	}
	var students []StudentDetails
	err = StudentDb.AllByIndex("Level", &students)
	if err != nil {
		log.Fatalln(err)
	}
	templates = templates.Lookup("table.html")
	err = templates.Execute(w, students)
	if err != nil {
		log.Fatalln(err)
	}
}

func (user *User) AddCoursePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	_, err := r.Cookie("admin")
	if err != nil {
		_, err = r.Cookie("lecturer")
		if err != nil {
			http.Redirect(w, r, "/login", 302)
			return
		}
	}
	r.ParseForm()
	if user.Details.IsAdmin != true {
		http.Redirect(w, r, "/login", 302)
		return
	}

	StormDb.Save(&Database{
		IsAdmin:  true,
		Name:     r.FormValue("course"),
		Password: r.FormValue("Pass"),
	})

	//fmt.Println(studentOfferingCourses, lecturerDB.Save(&course))
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

// func (user *User) LecturerLoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	r.ParseForm()
// 	var lecturer LecturerDetails
// 	Code := r.FormValue("code")
// 	Password := r.FormValue("Pass")
// 	fmt.Println(Code, Password)
// 	err := StormDb.One("Password", Password, &lecturer)
// 	if err != nil || lecturer.CourseCode != Code {
// 		user.Login.BadDetail = true
// 		fmt.Println(err)
// 		http.Redirect(w, r, "/lecturer", 302)
// 		return
// 	}
// 	http.SetCookie(
// 		w,
// 		&http.Cookie{
// 			Path: "/",
// 			Name: "lecturer",
// 		})
// 	http.Redirect(w, r, "/", 302)
// }

// func (user *User) LecturerLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	templates = templates.Lookup("lecturerLogin.html")
// 	err := templates.Execute(w, user)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// }
