package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Json(response http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("error while mashalling object %v, trace: %+v", data, err)
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_, err = response.Write(bytes)
	if err != nil {
		log.Fatalf("error while writting bytes to response writer: %+v", err)
	}
}

func RouteParam(request *http.Request, name string) string {
	return mux.Vars(request)[name]
}

func SearchImages(response http.ResponseWriter, request *http.Request) {
	//searchTerm := RouteParam(request, "searchTerm")
	//images := services.searchImagesBySearchTerm(searchTerm)
	Json(response, http.StatusOK, nil)
}
