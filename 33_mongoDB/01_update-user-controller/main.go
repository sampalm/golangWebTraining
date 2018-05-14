package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sampalm/33_mongoDB/01_update-user-controller/controllers"
	"gopkg.in/mgo.v2"
)

func main() {
	router := httprouter.New()
	uc := controllers.NewUserController(getSession())
	router.GET("/user/:id", uc.GetUser)
	router.GET("/users/", uc.GetAllUsers)
	router.POST("/user/", uc.CreateUser)
	router.DELETE("/user/:id", uc.DeleteUser)
	http.ListenAndServe(":8080", router)
}
func getSession() *mgo.Session {
	s, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}
	return s
}
