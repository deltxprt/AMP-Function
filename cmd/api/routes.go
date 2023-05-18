package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/status", app.listInstancesHandler)

	router.HandlerFunc(http.MethodGet, "/v1/status/:instance", app.ListInstanceHandler)

	router.HandlerFunc(http.MethodPost, "/v1/action/:action", app.actionInstanceHandler)

	return router
}
