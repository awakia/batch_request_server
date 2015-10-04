package main

import (
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Request defines each request
type Request struct {
	Method      string
	RelativeURL string
	Body        string
	Name        string
}

func batchRequests(requests []*Request, endPoint *url.URL) {
	client := &http.Client{}
	responses := make([]*http.Response, len(requests))
	// TODO: change to use go rutine
	for i, request := range requests {
		log.Println("Resuest:", request.Method, request.RelativeURL)
		url, _ := endPoint.Parse(request.RelativeURL)
		req, _ := http.NewRequest(request.Method, url.String(), strings.NewReader(request.Body))
		resp, _ := client.Do(req)
		responses[i] = resp
		// TODO: convert response to result structure
	}
}

func main() {
	endPoint, _ := url.Parse("http://localhost:3000")
	requests := []*Request{
		&Request{
			Method:      "GET",
			RelativeURL: "api/v1/projects",
		},
	}
	batchRequests(requests, endPoint)
}
