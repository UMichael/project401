package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	var user User
	router := httprouter.New()
	router.GET("/", user.HomePage)
	router.GET("/addcourse", user.AddCourse)
	router.GET("/changePass", user.ChangePasswordGet)
	router.GET("/login", LoginGet)
	router.GET("/students", user.Students)
	router.GET("/lecturer", user.LecturerLogin)

	router.POST("/changePass", user.ChangePassword)
	router.POST("/enroll", user.AdminEnroll)
	router.POST("/addcourse", user.AddCoursePost)
	router.POST("/lecturer", user.LecturerLoginPost)
	router.POST("/login", user.LoginPost)
	router.ServeFiles("/bootstrap4/*filepath", http.Dir("./templates/bootstrap4"))
	router.ServeFiles("/js/*filepath", http.Dir("./templates/js"))
	http.ListenAndServe(":8080", router)
}
