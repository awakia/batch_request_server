package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// Request defines each request
type Request struct {
	Method      string `json:"method"`
	RelativeURL string `json:"relative_url"`
	Body        string `json:"body"`
	Name        string `json:"name"`
}

// Response defines corresponding response
type Response struct {
	Code    int                 `json:"code"`
	Headers map[string][]string `json:"headers"`
	Body    string              `json:"body"`
}

// NewResponse creates response instance from http response
func NewResponse(resp *http.Response) *Response {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return &Response{
		resp.StatusCode,
		resp.Header,
		string(body),
	}
}

func batchRequests(requests []*Request, endPoint *url.URL) []*Response {
	client := &http.Client{}
	responses := make([]*Response, len(requests))
	// TODO: change to use go rutine
	for i, request := range requests {
		log.Println("Resuest:", request.Method, request.RelativeURL)
		url, _ := endPoint.Parse(request.RelativeURL)
		req, _ := http.NewRequest(request.Method, url.String(), strings.NewReader(request.Body))
		resp, _ := client.Do(req)
		responses[i] = NewResponse(resp)
	}
	return responses
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
