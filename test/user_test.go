package test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wan-system/internal/model"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		ID:       "fadel",
		Password: "password123",
		Name:     "Fadel Muhammad",
		Email:    "fadel@example.com",
	}
	bodyJson, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.UserResponse])
	json.Unmarshal(bytes, responseBody)

	assert.Equal(t, http.StatusCreated, response.StatusCode)
	assert.Equal(t, requestBody.ID, responseBody.Data.ID)
	assert.Equal(t, requestBody.Name, responseBody.Data.Name)
}

func TestRegisterError(t *testing.T) {
	ClearAll()
	requestBody := model.RegisterUserRequest{
		ID:    "",
		Email: "not-an-email",
	}
	bodyJson, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

func TestLogin(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		ID:       "fadel",
		Password: "password123",
	}
	bodyJson, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.UserResponse])
	json.Unmarshal(bytes, responseBody)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotEmpty(t, responseBody.Data.Token)
}

func TestLoginWrongPassword(t *testing.T) {
	ClearAll()
	TestRegister(t)

	requestBody := model.LoginUserRequest{
		ID:       "fadel",
		Password: "wrongpassword",
	}
	bodyJson, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusUnauthorized, response.StatusCode)
}

func TestGetCurrentUser(t *testing.T) {
	ClearAll()
	TestRegister(t)

	// Get token directly from DB or login
	requestBody := model.LoginUserRequest{
		ID:       "fadel",
		Password: "password123",
	}
	bodyJson, _ := json.Marshal(requestBody)
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	response, _ := app.Test(request)

	bytes, _ := io.ReadAll(response.Body)
	loginResp := new(model.WebResponse[model.UserResponse])
	json.Unmarshal(bytes, loginResp)
	token := loginResp.Data.Token

	// Final request
	request = httptest.NewRequest(http.MethodGet, "/api/auth/current", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, _ = io.ReadAll(response.Body)
	currResp := new(model.WebResponse[model.UserResponse])
	json.Unmarshal(bytes, currResp)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "fadel", currResp.Data.ID)
}

func TestLogout(t *testing.T) {
	ClearAll()
	TestRegister(t)

	// Login to get token
	requestBody := model.LoginUserRequest{
		ID:       "fadel",
		Password: "password123",
	}
	bodyJson, _ := json.Marshal(requestBody)
	request := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	response, _ := app.Test(request)
	bytes, _ := io.ReadAll(response.Body)
	loginResp := new(model.WebResponse[model.UserResponse])
	json.Unmarshal(bytes, loginResp)
	token := loginResp.Data.Token

	// Logout
	request = httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

	// Verify token cleared in DB
	user := GetUser("fadel")
	assert.Empty(t, user.Token)
}
