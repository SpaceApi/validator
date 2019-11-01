package main

import (
	"github.com/rs/cors"
	"github.com/spaceapi/validator/v1"
	"github.com/spaceapi/validator/v2"
	"goji.io"
	"goji.io/pat"
	"log"
	"net/http"
)

func main() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	root := goji.NewMux()
	root.Use(c.Handler)

	root.HandleFunc(pat.Get("/"), versionRedirect)
	root.HandleFunc(pat.Get("/openapi.json"), openApi)

	root.Handle(pat.New("/v1/*"), v1.GetValidatorV1Mux())
	root.Handle(pat.New("/v2/*"), v2.GetValidatorV2Mux())

	log.Println("starting validator...")
	log.Fatal(http.ListenAndServe(":8080", root))
}

func versionRedirect(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/v1", 302)
}

func openApi(writer http.ResponseWriter, request *http.Request) {
	_, _ = writer.Write([]byte(openapi))
}
