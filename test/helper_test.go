package test

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
	"wan-system/internal/entity"
)

func ClearAll() {
	if err := db.Exec("DELETE FROM wan_traffics").Error; err != nil {
		panic(err)
	}
	if err := db.Exec("DELETE FROM wan_capacities").Error; err != nil {
		panic(err)
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		panic(err)
	}
}

func ClearAuth() {
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		panic(err)
	}
}

func GetUser(id string) *entity.User {
	user := new(entity.User)
	if err := db.Where("id = ?", id).First(user).Error; err != nil {
		return nil
	}
	return user
}

func LoginHelper(t *testing.T, id string, password string) string {
	requestBody := map[string]string{
		"id":       id,
		"password": password,
	}
	bodyJson, _ := json.Marshal(requestBody)
	request := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")

	response, _ := app.Test(request)
	bytes, _ := io.ReadAll(response.Body)

	loginResp := new(struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	})
	json.Unmarshal(bytes, loginResp)
	return loginResp.Data.Token
}
