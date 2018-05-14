package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/sampalm/33_mongoDB/01_update-user-controller/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var defaultDB = "goweb"

type UserController struct {
	session *mgo.Session
}

func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func (uc *UserController) GetUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-type", "application/json")
	id := p.ByName("id")
	// Check if ID is an ObjectID hex representation
	if !bson.IsObjectIdHex(id) {
		je := models.JsonErr{
			Type:   "bad request",
			Status: http.StatusBadRequest,
			Error:  "Invalid ID",
		}
		json.NewEncoder(w).Encode(je)
		return
	}

	// fetch user
	u := models.User{}
	err := uc.session.DB(defaultDB).C("users").FindId(id).One(&u)
	if err != nil {
		je := models.JsonErr{
			Type:   "database query",
			Status: 500,
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(je)
		return
	}
	// Set headers
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

func (uc *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "application/json")
	var profiles []models.User
	err := uc.session.DB(defaultDB).C("users").Find(bson.M{}).All(&profiles)
	if err != nil {
		je := models.JsonErr{
			Type:   "database query all",
			Status: 500,
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(je)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(profiles)
}

func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-type", "application/json")
	// decode post body
	u := models.User{}
	json.NewDecoder(r.Body).Decode(&u)

	// insert decoded data into mongodb
	u.ID = bson.NewObjectId().Hex()
	_, err := uc.session.DB(defaultDB).C("users").UpsertId(u.ID, u) //select db -> collection/table -> insert json object/data
	if err != nil {
		je := models.JsonErr{
			Type:   "database insert",
			Status: 500,
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(je)
		return
	}
	w.WriteHeader(http.StatusCreated) // 201
	json.NewEncoder(w).Encode(u)
}

func (uc *UserController) DeleteUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	w.Header().Set("Content-type", "application/json")
	id := p.ByName("id")
	// Check if ID is an ObjectID hex representation
	if !bson.IsObjectIdHex(id) {
		je := models.JsonErr{
			Type:   "bad request",
			Status: http.StatusBadRequest,
			Error:  "Invalid ID",
		}
		json.NewEncoder(w).Encode(je)
		return
	}

	// Delete user from database
	err := uc.session.DB(defaultDB).C("users").RemoveId(id)
	if err != nil {
		je := models.JsonErr{
			Type:   "database delete",
			Status: 500,
			Error:  err.Error(),
		}
		json.NewEncoder(w).Encode(je)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "User %s deleted successfully.\n", id)
}
