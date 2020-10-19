package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/agile-engine-image-search/controllers"
	"github.com/agile-engine-image-search/model"
	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/search/${searchTerm}", controllers.SearchImages)
	http.ListenAndServe(":8080", myRouter)
}

var authToken *string
var pictures []model.Picture = []model.Picture{}
var picturesFullDataRemaining int

func getAuthToken() {
	request := authTokenRequest()

	if request == nil {
		return
	}

	performAuthTokenRequest(request)
}

func performAuthTokenRequest(request *http.Request) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Do(request)

	if err != nil {
		log.Fatalf("error performing Auth request: %v", err)
		return
	}

	authToken = getTokenFromBody(response)
}

func getTokenFromBody(response *http.Response) *string {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("error reading body: %v", err)
		return nil
	}

	var authToken model.AuthToken

	err = json.Unmarshal(body, &authToken)

	if err != nil {
		log.Fatalf("error unmarshalling body: %v", err)
		return nil
	}

	return &authToken.Token
}

func authTokenRequest() *http.Request {
	requestBody, err := json.Marshal(map[string]string{
		"apiKey": "23567b218376f79d9415",
	})

	if err != nil {
		log.Fatalf("error marshalling request body %v", err)
		return nil
	}

	request, err := http.NewRequest("POST", "http://interview.agileengine.com/auth", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		log.Fatalf("error creating Auth request %v", err)
		return nil
	}
	return request
}

func getPicturesResponseFromBody(response *http.Response) *model.PicturesResponse {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("error reading body: %v", err)
		return nil
	}

	var picturesResponse model.PicturesResponse

	err = json.Unmarshal(body, &picturesResponse)

	if err != nil {
		log.Fatalf("error unmarshalling body: %v", err)
		return nil
	}

	return &picturesResponse
}

func getPicturesRequest(pageNum int) *http.Request {
	request, err := http.NewRequest("GET", "http://interview.agileengine.com/images?page="+strconv.Itoa(pageNum), nil)
	request.Header.Set("Authorization", "Bearer "+*authToken)

	if err != nil {
		log.Fatalf("error creating pictures request %v", err)
		return nil
	}
	return request
}

func getPicturesResponse(request *http.Request) *model.PicturesResponse {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Do(request)

	if err != nil {
		log.Fatalf("error performing get Pictures request: %v", err)
		return nil
	}

	PicturesResponse := getPicturesResponseFromBody(response)

	if PicturesResponse == nil {
		return nil
	}
	return PicturesResponse
}

func downloadPictures() {
	getAuthToken()
	if authToken == nil {
		return
	}
	hasMorePictures := true
	page := 1
	for hasMorePictures {
		log.Print("Downloading page " + strconv.Itoa(page))
		request := getPicturesRequest(page)
		if request == nil {
			return
		}
		picturesResponse := getPicturesResponse(request)
		pictures = append(pictures, picturesResponse.Pictures...)
		page++
		hasMorePictures = picturesResponse.HasMore
	}
}

func getPictureFromBody(response *http.Response) *model.Picture {
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("error reading body: %v", err)
		return nil
	}

	var picture model.Picture

	err = json.Unmarshal(body, &picture)

	if err != nil {
		log.Fatalf("error unmarshalling body: %v", err)
		return nil
	}

	return &picture
}

func getPictureFullDataRequest(index int) *http.Request {
	request, err := http.NewRequest("GET", "http://interview.agileengine.com/images/"+pictures[index].Id, nil)
	request.Header.Set("Authorization", "Bearer "+*authToken)

	if err != nil {
		log.Fatalf("error creating picture full data request %v", err)
		return nil
	}
	return request
}

func addPictureFullData(request *http.Request, index int) {
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}

	response, err := client.Do(request)

	if err != nil {
		log.Fatalf("error performing get Pictures request: %v", err)
		return
	}

	picture := getPictureFromBody(response)

	if picture == nil {
		log.Fatalf("error reading picture from body: %v", err)
		return
	}
	pictures[index] = *picture
}

func addPicturesFullData() {
	log.Printf("Loading Full info from " + strconv.Itoa(len(pictures)) + " pictures")
	picturesFullDataRemaining = len(pictures)
	for i := range pictures {
		go func(i int) {
			request := getPictureFullDataRequest(i)
			if request == nil {
				return
			}
			addPictureFullData(request, i)
			picturesFullDataRemaining--
			log.Printf("Loading Full info from picture" + strconv.Itoa(i))
			if picturesFullDataRemaining == 0 {
				//services.saveImages(pictures)
			}
		}(i)
	}
}

func downloadPicturesAndStartServer() {
	downloadPictures()
	addPicturesFullData()
	log.Print("Starting server")
	handleRequests()
}

func main() {
	downloadPicturesAndStartServer()
}
