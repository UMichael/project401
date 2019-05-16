package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func main() {
	var lecturer User
	router := httprouter.New()
	router.GET("/", HomePage)
	router.GET("/login", LoginGet)
	router.POST("/login", lecturer.LoginPost)
	http.ListenAndServe(":8080", router)
}
