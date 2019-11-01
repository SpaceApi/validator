package v2

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var validSpace = `{
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
}`

var invalidSpace = `{ "data": "asd" }`

func forgeValidateJsonRequest(t *testing.T, body io.Reader) (error, *httptest.ResponseRecorder) {
	req, err := http.NewRequest("POST", "/v2/validateJson", body)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(validateJson)
	handler.ServeHTTP(rr, req)
	return err, rr
}

func forgeValidateUrlRequest(t *testing.T, body io.Reader) (error, *httptest.ResponseRecorder) {
	req, err := http.NewRequest("POST", "/v2/validateUrl", body)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(validateUrl)
	handler.ServeHTTP(rr, req)
	return err, rr
}

//// VALIDATE JSON ////

func TestValidateJsonWithValid(t *testing.T) {
	err, rr := forgeValidateJsonRequest(t, strings.NewReader(validSpace))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := jsonValidationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Valid != true {
		t.Errorf("handler returned wrong response: got %v want %v",
			resp.Valid, true)
	}
}

func TestValidateJsonWithInvalid(t *testing.T) {
	err, rr := forgeValidateJsonRequest(t, strings.NewReader(invalidSpace))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := jsonValidationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Valid != false {
		t.Errorf("handler returned wrong response: got %v want %v",
			resp.Valid, false)
	}
}

func TestValidateJsonWithInvalidJson(t *testing.T) {
	_, rr := forgeValidateJsonRequest(t, strings.NewReader("foo"))

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestValidateJsonWithEmptyBody(t *testing.T) {
	_, rr := forgeValidateJsonRequest(t, nil)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

//// VALIDATE URL ////

func TestValidateUrlWithValid(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(validSpace))
		}))
	defer ts.Close()

	err, rr := forgeValidateUrlRequest(t, strings.NewReader(`{ "url": "`+ts.URL+`" }`))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := urlValidationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Valid != true {
		t.Errorf("handler returned wrong response: got %v want %v",
			resp.Valid, true)
	}
}

func TestValidateUrlWithUnreachablePath(t *testing.T) {
	err, rr := forgeValidateUrlRequest(t, strings.NewReader(`{ "url": "https://example.com/status.json" }`))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := urlValidationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Valid != false {
		t.Errorf("handler returned wrong response: got %v want %v",
			resp.Valid, true)
	}

	if resp.Reachable != false {
		t.Errorf("handler returned wrong reachability: got %v want %v",
			resp.Reachable, false)
	}
}

func TestValidateUrlWithUnreachableServer(t *testing.T) {
	err, rr := forgeValidateUrlRequest(t, strings.NewReader(`{ "url": "http://localhost:666/status.json" }`))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := urlValidationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Valid != false {
		t.Errorf("handler returned wrong response: got %v want %v",
			resp.Valid, true)
	}

	if resp.Reachable != false {
		t.Errorf("handler returned wrong reachability: got %v want %v",
			resp.Reachable, false)
	}
}

func TestValidateUrlWithEmptyBody(t *testing.T) {
	_, rr := forgeValidateUrlRequest(t, nil)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestValidateUrlWithEmptyStringBody(t *testing.T) {
	_, rr := forgeValidateUrlRequest(t, strings.NewReader(""))

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestValidateUrlWithInvalidBody(t *testing.T) {
	_, rr := forgeValidateUrlRequest(t, strings.NewReader(`{}`))

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

func TestValidateUrlCors(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Content-Type", "application/json")
			_, _ = w.Write([]byte(validSpace))
		}))
	defer ts.Close()

	err, rr := forgeValidateUrlRequest(t, strings.NewReader(`{ "url": "`+ts.URL+`" }`))

	resp := urlValidationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Cors != true {
		t.Errorf("cors check failed: got %v want %v",
			resp.Cors, true)
	}

	if resp.ContentType != true {
		t.Errorf("content type check failed: got %v want %v",
			resp.ContentType, true)
	}
}

func TestValidateUrlCorsFalse(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(validSpace))
		}))
	defer ts.Close()

	err, rr := forgeValidateUrlRequest(t, strings.NewReader(`{ "url": "`+ts.URL+`" }`))

	resp := urlValidationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Cors != false {
		t.Errorf("cors check failed: got %v want %v",
			resp.Cors, false)
	}
	if resp.ContentType != false {
		t.Errorf("content type check failed: got %v want %v",
			resp.ContentType, false)
	}
}

func TestValidateUrlInvalidSpace(t *testing.T) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(invalidSpace))
		}))
	defer ts.Close()

	err, rr := forgeValidateUrlRequest(t, strings.NewReader(`{ "url": "`+ts.URL+`" }`))

	resp := urlValidationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.Valid != false {
		t.Errorf("cors check failed: got %v want %v",
			resp.Valid, false)
	}

	if resp.Message == "" {
		t.Errorf("message should not be empty")
	}
}

func TestValidateUrlInvalidTls(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(validSpace))
		}))
	defer ts.Close()

	err, rr := forgeValidateUrlRequest(t, strings.NewReader(`{ "url": "`+ts.URL+`" }`))

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	resp := urlValidationResponse{}
	err = json.NewDecoder(rr.Body).Decode(&resp)
	if err != nil {
		t.Fatal(err)
	}

	if resp.IsHttps != true {
		t.Errorf("https check failed: got %v want %v",
			resp.IsHttps, true)
	}

	if resp.CertValid != false {
		t.Errorf("cert check failed: got %v want %v",
			resp.CertValid, false)
	}
}

func TestServerInfo(t *testing.T) {
	req, err := http.NewRequest("POST", "/v2", nil)
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
