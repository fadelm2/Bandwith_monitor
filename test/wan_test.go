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

func TestCreateCapacity(t *testing.T) {
	ClearAll()
	TestRegister(t)
	token := LoginHelper(t, "fadel", "password123")

	requestBody := model.WanCapacityRequest{
		WanID:            "WAN-TEST-1",
		CapacityMbps:     1000.0,
		ThresholdPercent: 80.0,
		Description:      "Test Connection",
	}
	bodyJson, _ := json.Marshal(requestBody)

	request := httptest.NewRequest(http.MethodPost, "/internal/capacity", strings.NewReader(string(bodyJson)))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.WanCapacityResponse])
	json.Unmarshal(bytes, responseBody)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "WAN-TEST-1", responseBody.Data.WanID)
}

func TestListCapacity(t *testing.T) {
	ClearAll()
	TestRegister(t)
	token := LoginHelper(t, "fadel", "password123")

	// Create one first
	TestCreateCapacity(t)

	request := httptest.NewRequest(http.MethodGet, "/internal/capacity", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.WanCapacityResponse])
	json.Unmarshal(bytes, responseBody)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.GreaterOrEqual(t, len(responseBody.Data), 1)
}

func TestGetCapacity(t *testing.T) {
	ClearAll()
	TestRegister(t)
	token := LoginHelper(t, "fadel", "password123")
	TestCreateCapacity(t)

	request := httptest.NewRequest(http.MethodGet, "/internal/capacity/WAN-TEST-1", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[model.WanCapacityResponse])
	json.Unmarshal(bytes, responseBody)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Equal(t, "WAN-TEST-1", responseBody.Data.WanID)
}

func TestSearchTraffic(t *testing.T) {
	ClearAll()
	TestRegister(t)
	token := LoginHelper(t, "fadel", "password123")

	// Search traffic (may be empty initially)
	request := httptest.NewRequest(http.MethodGet, "/internal/traffic?page=1&size=10", nil)
	request.Header.Set("Authorization", "Bearer "+token)

	response, err := app.Test(request)
	assert.Nil(t, err)

	bytes, _ := io.ReadAll(response.Body)
	responseBody := new(model.WebResponse[[]model.WanTrafficResponse])
	json.Unmarshal(bytes, responseBody)

	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.NotNil(t, responseBody.Paging)
}
