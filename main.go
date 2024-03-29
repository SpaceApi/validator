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
	root.HandleFunc(pat.Get("/openapi.json"), openAPI)

	root.HandleFunc(pat.Get("/v1"), func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/v1/", 302)
	})
	root.HandleFunc(pat.Get("/v2"), func(writer http.ResponseWriter, request *http.Request) {
		http.Redirect(writer, request, "/v2/", 302)
	})

	root.Handle(pat.New("/v1/*"), v1.GetSubMux())
	root.Handle(pat.New("/v2/*"), v2.GetSubMux())

	log.Println("starting validator on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", root))
}

func versionRedirect(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/v1/", 302)
}

func openAPI(writer http.ResponseWriter, _ *http.Request) {
	_, _ = writer.Write([]byte(openapi))
}
