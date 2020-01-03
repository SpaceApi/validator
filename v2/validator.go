package v2

import (
	"crypto/tls"
	"encoding/json"
	spaceapivalidator "github.com/spaceapi-community/go-spaceapi-validator"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/time/rate"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type serverInfo struct {
	Description string `json:"description"`
	Usage       string `json:"usage"`
	Version     string `json:"version"`
}

type urlValidationRequest struct {
	URL string `json:"url"`
}

type urlValidationResponse struct {
	Valid        bool          `json:"valid"`
	Message      string        `json:"message,omitempty"`
	IsHTTPS      bool          `json:"isHttps"`
	HTTPSForward bool          `json:"httpsForward"`
	Reachable    bool          `json:"reachable"`
	Cors         bool          `json:"cors"`
	ContentType  bool          `json:"contentType"`
	CertValid    bool          `json:"certValid"`
	ValidatedJson interface{} `json:"validatedJson,omitempty"`
	SchemaErrors []schemaError `json:"schemaErrors,omitempty"`
}

type schemaError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type jsonValidationResponse struct {
	Valid        bool          `json:"valid"`
	Message      string        `json:"message"`
	ValidatedJson interface{} `json:"validatedJson,omitempty"`
	SchemaErrors []schemaError `json:"schemaErrors,omitempty"`
}

// GetSubMux returns the versions subrouter
func GetSubMux() *goji.Mux {
	v2 := goji.SubMux()
	v2.HandleFunc(pat.Get("/"), info)
	v2.HandleFunc(pat.Post("/validateJSON"), validateJSON)
	v2.Handle(
		pat.Post("/validateURL"),
		limit(
			http.HandlerFunc(validateURL),
			rate.NewLimiter(10, 25),
		),
	)

	return v2
}

func limit(next http.Handler, limiter *rate.Limiter) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func info(writer http.ResponseWriter, _ *http.Request) {
	serverInfo := serverInfo{
		Description: "Space API Validator API",
		Usage:       "Send a POST request in JSON format to /v2/validateJSON. See https://github.com/SpaceApi/validator for more information.",
		Version:     "1.0.0",
	}

	err := json.NewEncoder(writer).Encode(serverInfo)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func validateURL(writer http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		http.Error(writer, "body can't be empty", http.StatusBadRequest)
		return
	}

	var valReq urlValidationRequest
	var valRes urlValidationResponse

	err := json.NewDecoder(request.Body).Decode(&valReq)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := url.ParseRequestURI(valReq.URL)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	valRes.IsHTTPS = u.Scheme == "https"

	header, body, err := fetchURL(&valRes, u, false)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if header != nil {
		checkHeader(&valRes, header)
	}

	if body == "" {
		writer.Header().Add("Content-Type", "application/json")
		err = json.NewEncoder(writer).Encode(valRes)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}

		return
	}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(body), &raw); err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	valRes.ValidatedJson = raw

	res, err := spaceapivalidator.Validate(body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	valRes.Valid = res.Valid
	var errMsg string
	for _, validatorError := range res.Errors {
		errMsg = errMsg + validatorError.Context + ": " + validatorError.Description + "\n"
		valRes.SchemaErrors = append(valRes.SchemaErrors, schemaError{
			Field:   validatorError.Context,
			Message: validatorError.Description,
		})
	}
	valRes.Message = errMsg

	writer.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(valRes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func checkHeader(response *urlValidationResponse, header http.Header) {
	acao := header.Get("Access-Control-Allow-Origin")
	if acao == "*" || acao == "https://validator.spaceapi.io" {
		response.Cors = true
	}

	if strings.HasPrefix(header.Get("Content-Type"), "application/json") {
		response.ContentType = true
	}
}

func fetchURL(validationResponse *urlValidationResponse, url *url.URL, skipVerify bool) (http.Header, string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	client := &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.Scheme == "https" {
				validationResponse.HTTPSForward = true
			}
			return nil
		},
		Transport: tr,
	}

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		validationResponse.Reachable = false
		return nil, "", err
	}

	req.Header.Add("Origin", "https://validator.spaceapi.io")
	response, err := client.Do(req)
	if err != nil {
		if skipVerify == false {
			return fetchURL(validationResponse, url, true)
		}

		validationResponse.Reachable = false
		return nil, "", nil
	}

	if response.StatusCode >= 400 {
		validationResponse.Reachable = false
		return nil, "", nil
	}

	bodyArray, _ := ioutil.ReadAll(response.Body)
	validationResponse.Reachable = true
	validationResponse.CertValid = (validationResponse.IsHTTPS || validationResponse.HTTPSForward) && !skipVerify
	return response.Header, string(bodyArray), nil
}

func validateJSON(writer http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		http.Error(writer, "body can't be empty", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := spaceapivalidator.Validate(string(body))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	resp := jsonValidationResponse{
		Valid: res.Valid,
		ValidatedJson: raw,
	}

	var errMsg string
	for _, validatorError := range res.Errors {
		errMsg = errMsg + validatorError.Context + ": " + validatorError.Description + "\n"
		resp.SchemaErrors = append(resp.SchemaErrors, schemaError{
			Field:   validatorError.Context,
			Message: validatorError.Description,
		})
	}
	resp.Message = errMsg

	writer.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(resp)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
