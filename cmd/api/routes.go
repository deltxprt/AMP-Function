package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthCheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/ampstatus", app.listInstancesHandler)

	router.HandlerFunc(http.MethodGet, "/v1/ampstatus/:instance", app.ListInstanceHandler)

	return router
}
