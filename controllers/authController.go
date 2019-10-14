package controllers

import (
	"encoding/json"
	c "github.com/GlobalWebIndex/platform2.0-go-challenge/common"
	"github.com/GlobalWebIndex/platform2.0-go-challenge/models"
	"net/http"
)

var CreateUser = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		c.Respond(w, http.StatusBadRequest, c.Message(false, "Invalid request"))
		return
	}

	resp := user.Create()//Create account
	if resp["status"] == true {
		c.Respond(w, http.StatusCreated, resp)
	} else {
		c.Respond(w, http.StatusBadRequest, resp)
	}
}

var AuthenticateUser = func(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		c.Respond(w, http.StatusBadRequest, c.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(user.Email, user.Password)
	if resp["status"] == true {
		c.Respond(w, http.StatusOK, resp)
	} else {
		c.Respond(w, http.StatusBadRequest, resp)
	}
}

var RefreshToken = func(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("user").(uint)
	user := models.GetUser(userId)
	user.SetToken()

	resp := c.Message(true, "Token Refreshed")
	resp["token"] = user.Token
	c.Respond(w, http.StatusOK, resp)
}


