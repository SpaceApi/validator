package v1

import (
	"encoding/json"
	spaceapivalidator "github.com/spaceapi-community/go-spaceapi-validator"
	"goji.io"
	"goji.io/pat"
	"net/http"
)

type serverInfo struct {
	Description string `json:"description"`
	Usage       string `json:"usage"`
	Version     string `json:"version"`
}

type validationRequest struct {
	Data interface{} `json:"data"`
}

type validationResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message"`
}

// GetSubMux returns the versions subrouter
func GetSubMux() *goji.Mux {
	v1 := goji.SubMux()
	v1.HandleFunc(pat.Get("/"), info)
	v1.HandleFunc(pat.Post("/validate"), validate)

	v1.HandleFunc(pat.Get("/validate"), func(writer http.ResponseWriter, _ *http.Request) {
		writer.WriteHeader(405)
	})
	v1.HandleFunc(pat.Get("/validate/"), forwardToValidate)
	v1.HandleFunc(pat.Post("/validate/"), forwardToValidate)

	return v1
}

func forwardToValidate(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/v1/validate", 302)
}

func info(writer http.ResponseWriter, request *http.Request) {
	serverInfo := serverInfo{
		Description: "Space API Validator API",
		Usage:       "Send a POST request in JSON format to /v1/validate. See https://github.com/SpaceApi/validator for more information.",
		Version:     "2.0.0",
	}

	err := json.NewEncoder(writer).Encode(serverInfo)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func validate(writer http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		http.Error(writer, "body has to be provided", http.StatusBadRequest)
		return
	}

	var req validationRequest
	err := json.NewDecoder(request.Body).Decode(&req)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	jsonString, err := json.Marshal(req.Data)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := spaceapivalidator.Validate(string(jsonString))
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	var errMsg string
	for _, validatorError := range res.Errors {
		errMsg = errMsg + validatorError.Context + ": " + validatorError.Description + "\n"
	}

	resp := validationResponse{
		Valid:   res.Valid,
		Message: errMsg,
	}

	writer.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(resp)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}
