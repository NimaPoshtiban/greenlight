package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

func (app *application) healthcheckHandler(w http.ResponseWriter,r *http.Request){
	status := map[string]string{
		"status" : "available",
		"environment": app.config.env,
		"version" : version,
	}
	response,_ := json.MarshalIndent(status," ","\t")
	fmt.Fprintln(w,string(response))
}