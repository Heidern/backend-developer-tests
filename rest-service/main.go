package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/stackpath/backend-developer-tests/rest-service/pkg/controllers"

	custommw "github.com/stackpath/backend-developer-tests/rest-service/pkg/middleware"
)

func main() {
	fmt.Println("SP// Backend Developer Test - RESTful Service")
	fmt.Println()

	logger := httplog.NewLogger("rest-api", httplog.Options{
		JSON: true,
	})

	// TODO: Add RESTful web service here

	r := chi.NewRouter()

	r.Use(httplog.RequestLogger(logger))
	r.Use(middleware.RequestID)
	r.Use(custommw.RequestIDHeader)

	r.Mount("/people", (&controllers.PersonController{}).Router())

	http.ListenAndServe(":2022", r)
}
