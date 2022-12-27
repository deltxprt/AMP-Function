package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/api/healthcheck", healthCheckHandler)

	router.HandlerFunc(http.MethodGet, "/api/ampstatus", ampInfoHandler)

	return router
}
