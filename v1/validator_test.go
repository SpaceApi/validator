package v1

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var validSpace = `{ 
	"data": { 
		"api": "0.13",
		"space": "my cool space",
		"logo": "https://example.com/logo.png",
		"url": "https://example.com",
		"location": {
			"address": "Ulmer Strasse 255, 70327 Stuttgart, Germany",
    		"lon": 9.236,
    		"lat": 48.777
		},
		"state": {
			"open": false
		},
		"contact": {
		},
		"issue_report_channels": [
			"email"
		]
	}
}`

var invalidType = `{ "data": "asd" }`
var invalidSpace = `{ "data": {} }`

func TestValidateWithValid(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/validate", strings.NewReader(validSpace))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(validate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := validationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Valid != true {
		t.Errorf("handler returned wrong response: got %v want %v",
			resp.Valid, true)
	}
}

func TestValidateWithInvalidSpace(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/validate", strings.NewReader(invalidSpace))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(validate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := validationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Valid != false {
		t.Errorf("handler returned wrong response: got %v want %v",
			resp.Valid, false)
	}

	if resp.Message == "" {
		t.Errorf("message should not be empty")
	}
}

func TestValidateWithInvalidType(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/validate", strings.NewReader(invalidType))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(validate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestValidateWithInvalidJson(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/validate", strings.NewReader("foo"))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(validate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestValidateWithEmptyBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1/validate", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(validate)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestServerInfo(t *testing.T) {
	req, err := http.NewRequest("POST", "/v1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(info)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := serverInfo{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}
}
