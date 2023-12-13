package main

import (
	"github.com/julienschmidt/httprouter"
	_ "github.com/nimaposhtiban/greenlight/cmd/api/docs"
	"net/http"

	"github.com/swaggo/http-swagger/v2"
)

func swaggerHandler(res http.ResponseWriter, req *http.Request, p httprouter.Params) {
	httpSwagger.WrapHandler(res, req)
}

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.GET("/v1/swagger/*filepath", swaggerHandler)

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	return router
}
