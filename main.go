package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()
	router.GET("/projectinfo/v1/:git/:owner/:project", gitHandler)
	http.ListenAndServe(":8080", router)
}
