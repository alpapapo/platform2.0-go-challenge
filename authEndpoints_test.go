package main

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

func TestCreateUserWhenMethodNotAllowed(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/user/register", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusMethodNotAllowed, response.Code)

	if body := response.Body.String(); body != "" {
		t.Errorf("Expected an empty body Got %s", body)
	}
}

func TestCreateUserWhenUserCreated(t *testing.T) {
	payload := []byte(`{"email":"panos@testidis.com","username":"tester", "password": "%(secret)"}`)
	request, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response := executeRequest(request)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestCreateUserWhenBadRequest(t *testing.T) {
	payload := []byte(`{"email":"takis@testidis.com","username":"tester"}`)
	request, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response := executeRequest(request)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestAuthenticateUserWhenUserNotFound(t *testing.T) {

	payload := []byte(`{"email":"pouthenas@nouser.com","password":"pouthenas"}`)
	request, err := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response := executeRequest(request)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestAuthenticateUserWhenUserExistsButWrongPassword(t *testing.T) {
	payload := []byte(`{"email":"takis@testidis.com","password":"wrong"}`)
	request, err := http.NewRequest("POST", "/api/user/login", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}
	response := executeRequest(request)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestRefreshTokenWhenNotValidToken(t *testing.T) {
	jwt := "notValidToken"
	request_refresh_token, _ := http.NewRequest("POST", "/api/token/refresh", nil)
	request_refresh_token.Header = map[string][]string{
		"Authorization": {fmt.Sprintf("Bearer %s", jwt)},
	}

	response := executeRequest(request_refresh_token)
	checkResponseCode(t, http.StatusForbidden, response.Code)
}