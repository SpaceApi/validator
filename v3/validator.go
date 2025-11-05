package v3

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	spaceapivalidator "github.com/spaceapi-community/go-spaceapi-validator"
	"goji.io"
	"goji.io/pat"
	"golang.org/x/time/rate"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultUserAgent = "SpaceApi Validator/1.3.0 (https://github.com/SpaceApi/validator)"
)

type serverInfo struct {
	Description string `json:"description"`
	Usage       string `json:"usage"`
	Version     string `json:"version"`
}

type jsonValidationResponse struct {
	Valid           bool          `json:"valid"`
	Message         string        `json:"message,omitempty"`
	CheckedVersions []string      `json:"checkedVersions"`
	ValidatedJson   interface{}   `json:"validatedJson,omitempty"`
	SchemaErrors    []schemaError `json:"schemaErrors,omitempty"`
}

type urlValidationRequest struct {
	URL                string `json:"url"`
	ExtendedValidation bool   `json:"extendedValidation,omitempty"`
	WithValidatedJson  bool   `json:"withValidatedJson,omitempty"`
}

type urlValidationParameters struct {
	AddressFamily AddressFamily `json:"addressFamily"`
	UserAgent     string        `json:"userAgent"`
}

type urlValidationResult struct {
	Parameters        urlValidationParameters `json:"parameters"`
	OriginalURL       string                  `json:"originalUrl"`
	RedirectedURL     string                  `json:"redirectedUrl"`
	Reachable         bool                    `json:"reachable"`
	HTTPStatus        int                     `json:"httpStatus"`
	IsHTTPS           bool                    `json:"isHttps"`
	HTTPSRedirect     bool                    `json:"httpsRedirect"`
	PermanentRedirect bool                    `json:"permanentRedirect"`
	CORS              bool                    `json:"cors"`
	ContentType       bool                    `json:"contentType"`
	InvalidCert       bool                    `json:"invalidCert"`
	Validation        jsonValidationResponse  `json:"validation"`
}

type urlValidationResponse struct {
	Results []urlValidationResult `json:"results"`
	Reachable       bool          `json:"reachable"`
	Valid           bool          `json:"valid"`
}

type schemaError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// GetSubMux returns the versions subrouter
func GetSubMux() *goji.Mux {
	v3 := goji.SubMux()
	v3.HandleFunc(pat.Get("/"), getInfo)
	v3.HandleFunc(pat.Post("/validateJSON"), postValidateJSON)
	v3.Handle(
		pat.Post("/validateURL"),
		limit(
			http.HandlerFunc(postValidateURL),
			rate.NewLimiter(200, 500), // (rate, burst)
		),
	)

	return v3
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

func getInfo(writer http.ResponseWriter, _ *http.Request) {
	serverInfo := serverInfo{
		Description: "Space API Validator API",
		Usage:       "To validate a JSON document, send a POST request to /v3/validateJSON. To validate an endpoint URL, send a POST request to /v3/validateURL. See /swagger for an API documentation.",
		Version:     "1.3.0",
	}

	err := json.NewEncoder(writer).Encode(serverInfo)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func validateJSON(document []byte, withValidatedJson bool) (*jsonValidationResponse, error) {

	res, err := spaceapivalidator.Validate(string(document))
	if err != nil {
		return nil, err
	}

	var raw map[string]interface{}
	if err := json.Unmarshal(document, &raw); err != nil {
		return nil, err
	}

	resp := jsonValidationResponse{
		Valid: res.Valid,
	}
	if withValidatedJson {
		resp.ValidatedJson = raw
	}

	for _, schema := range res.Schemas {
		resp.CheckedVersions = append(resp.CheckedVersions, schema.Version)
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

	return &resp, nil
}

func postValidateJSON(writer http.ResponseWriter, request *http.Request) {
	if request.Body == nil {
		http.Error(writer, "body can't be empty", http.StatusBadRequest)
		return
	}
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := validateJSON(body, true)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(resp)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getTestCases(extended bool) []urlValidationParameters {
	cases := []urlValidationParameters{}
	if (extended) {
		cases = append(cases, urlValidationParameters {
			AddressFamily: IPv6,
			UserAgent: defaultUserAgent,
		})
		cases = append(cases, urlValidationParameters {
			AddressFamily: IPv4,
			UserAgent: defaultUserAgent,
		})
		cases = append(cases, urlValidationParameters {
			AddressFamily: Any,
			UserAgent: "curl/8.16.0",
		})
		cases = append(cases, urlValidationParameters {
			AddressFamily: Any,
			UserAgent: "Go-http-client/1.1",
		})
	} else {
		cases = append(cases, urlValidationParameters {
			AddressFamily: Any,
			UserAgent: defaultUserAgent,
		})
	}
	return cases
}

func validateURL(params urlValidationParameters, req urlValidationRequest) (*urlValidationResult, error) {
	result := urlValidationResult{
		Parameters: params,
	}

	u, err := url.ParseRequestURI(req.URL)
	if err != nil {
		return nil, err
	}
	u.Fragment = ""
	result.OriginalURL = u.String()
	result.IsHTTPS = u.Scheme == "https"

	header, body, err := fetchURL(&result, &params, u, false)
	if err != nil {
		return nil, err
	}
	result.Reachable = true

	if header != nil {
		checkHeader(&result, header)
	}

	if body == "" {
		return &result, nil
	}

	validation, err := validateJSON([]byte(body), req.WithValidatedJson)
	if err != nil {
		return nil, err
	}

	result.Validation = *validation
	return &result, nil
}

func postValidateURL(writer http.ResponseWriter, request *http.Request) {
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
	testCases := getTestCases(valReq.ExtendedValidation)

	valRes.Results = make([]urlValidationResult, len(testCases))
	for i, testCase := range testCases {
		result, err := validateURL(testCase, valReq)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		valRes.Results[i] = *result
	}

	// Set "convenience flags" in the outer response
	valRes.Valid = true
	valRes.Reachable = false
	for _, result := range valRes.Results {
		if result.Reachable {
			valRes.Reachable = true
			if !result.Validation.Valid {
				valRes.Valid = false
			}
		}
	}

	writer.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(writer).Encode(valRes)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func checkHeader(result *urlValidationResult, header http.Header) {
	acao := header.Get("Access-Control-Allow-Origin")
	if acao == "*" {
		result.CORS = true
	}

	if strings.HasPrefix(header.Get("Content-Type"), "application/json") {
		result.ContentType = true
	}
}

func fetchURL(result *urlValidationResult, params *urlValidationParameters, url *url.URL, skipVerify bool) (http.Header, string, error) {
	dialer := net.Dialer{}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerify},
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, params.AddressFamily.TcpStack(), addr)
		},
	}

	result.PermanentRedirect = true
	client := http.Client{
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return errors.New("Too many redirects")
			}
			result.RedirectedURL = req.URL.String()
			if req.Response.StatusCode != http.StatusMovedPermanently && req.Response.StatusCode != http.StatusPermanentRedirect {
				result.PermanentRedirect = false
			}
			if req.URL.Scheme == "https" {
				result.HTTPSRedirect = true
			}
			return nil
		},
		Transport: tr,
	}
	defer client.CloseIdleConnections()

	req, err := http.NewRequest("GET", url.String(), nil)
	if err != nil {
		return nil, "", err
	}

	req.Header.Add("Origin", "https://validator.spaceapi.io")
	req.Header.Add("User-Agent", params.UserAgent)
	response, err := client.Do(req)
	if err != nil {
		if skipVerify == false {
			return fetchURL(result, params, url, true)
		}

		return nil, "", nil
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	result.Reachable = true
	result.HTTPStatus = response.StatusCode
	if response.StatusCode >= 400 {
		_, _ = io.Copy(ioutil.Discard, response.Body)
		return nil, "", nil
	}

	bodyArray, _ := ioutil.ReadAll(response.Body)
	result.InvalidCert = !(result.IsHTTPS || result.HTTPSRedirect) || skipVerify
	return response.Header, string(bodyArray), nil
}
