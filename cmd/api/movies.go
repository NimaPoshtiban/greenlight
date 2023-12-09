package main

import (
	"fmt"
	"net/http"
)

func (app *application) createMovieHandler(w http.ResponseWriter,r *http.Request){
	fmt.Fprint(w,"Create movie using post")
}

func (app *application) showMovieHandler(w http.ResponseWriter,r *http.Request){

	id,err := app.readIdParam(r)

	if err != nil  {
		http.NotFound(w,r)
		return
	}

	fmt.Fprintf(w,"Showing movie by Id %d",id)
}