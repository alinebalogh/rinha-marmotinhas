package main

import (
	"net/http"
	"rinhago/internal/api"
)

func main() {

	router := api.SetupRouter()

	http.ListenAndServe(":8000", router)
}
