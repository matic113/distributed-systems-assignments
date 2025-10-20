package controllers

import (
	"encoding/json"
	"net/http"
	"sec_2/models"
	"strconv"
)

type userController struct {
}

func (uc *userController) post(w http.ResponseWriter, r *http.Request) {
	u, _ := uc.parseRequest(r)
	u, _ = models.AddUser(u)
	encodeResponseAsJSON(u, w)
}

func (uc *userController) get(id int, w http.ResponseWriter) {
	u, _ := models.GetUserByID(id)
	encodeResponseAsJSON(u, w)
}

func (uc *userController) parseRequest(r *http.Request) (models.User, error) {
	dec := json.NewDecoder(r.Body)
	var u models.User
	err := dec.Decode(&u)
	if err != nil {
		return models.User{}, err
	}
	return u, nil
}

func (uc userController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/users" {
		switch r.Method {
		case http.MethodGet:
			user_id, _ := strconv.Atoi(r.URL.Query()["id"][0])
			uc.get(user_id, w)
		case http.MethodPost:
			uc.post(w, r)
		}
	} else {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func newUserController() *userController {
	return &userController{}
}
