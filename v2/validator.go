package v2

import (
	"crypto/tls"
	"encoding/json"
	spaceapivalidator "github.com/spaceapi-community/go-spaceapi-validator"
	"goji.io"
	"goji.io/pat"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type serverInfo struct {
	Description string `json:"description"`
	Usage       string `json:"usage"`
	Version     string `json:"version"`
}

type urlValidationRequest struct {
	Url string `json:"url"`
}

type urlValidationResponse struct {
	Valid        bool   `json:"valid"`
	Message      string `json:"message,omitempty"`
	IsHttps      bool   `json:"isHttps"`
	HttpsForward bool   `json:"httpsForward"`
	Reachable    bool   `json:"reachable"`
	Cors         bool   `json:"cors"`
	ContentType  bool   `json:"contentType"`
	CertValid    bool   `json:"certValid"`
	SchemaErrors	[]schemaError	`json:"schemaErrors,omitempty"`
}

type schemaError struct {
	Field	string	`json:"field"`
	Message	string	`json:"message"`
}

type jsonValidationResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
	SchemaErrors	[]schemaError	`json:"schemaErrors,omitempty"`
}

func GetValidatorV2Mux() *goji.Mux {
	v2 := goji.SubMux()
	v2.HandleFunc(pat.Get("/"), info)
	v2.HandleFunc(pat.Post("/validateJson"), validateJson)
	v2.HandleFunc(pat.Post("/validateUrl"), validateUrl)

	return v2
}

func info(writer http.ResponseWriter, _ *http.Request) {
	serverInfo := serverInfo{
		Description: "Space API Validator API",
		Usage:       "Send a POST request in JSON format to /v2/validateJson. See https://github.com/SpaceApi/validator for more information.",
		Version:     "1.0.0",
	}

	err := json.NewEncoder(writer).Encode(serverInfo)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func validateUrl(writer http.ResponseWriter, request *http.Request) {
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

	u, err := url.ParseRequestURI(valReq.Url)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	valRes.IsHttps = u.Scheme == "https"

	header, body, err := fetchUrl(&valRes, u, false)
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

	if header.Get("Content-Type") == "application/json" {
		response.ContentType = true
	}
}

func fetchUrl(validationResponse *urlValidationResponse, url *url.URL, skipVerify bool) (http.Header, string, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
	}

	client := &http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if req.URL.Scheme == "https" {
				validationResponse.HttpsForward = true
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
			return fetchUrl(validationResponse, url, true)
		}

		validationResponse.Reachable = false
		return nil, "", nil
	}

	if response.StatusCode >= 400 {
		validationResponse.Reachable = false
		return nil, "", nil
	}

	bodyArray, err := ioutil.ReadAll(response.Body)
	validationResponse.Reachable = true
	validationResponse.CertValid = (validationResponse.IsHttps || validationResponse.HttpsForward) && !skipVerify
	return response.Header, string(bodyArray), nil
}

func validateJson(writer http.ResponseWriter, request *http.Request) {
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

	resp := jsonValidationResponse{
		Valid:   res.Valid,
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
