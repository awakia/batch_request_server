package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/unrolled/render"
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
	Code    int          `json:"code"`
	Headers *http.Header `json:"headers"`
	Body    string       `json:"body"`
}

// NewResponse creates response instance from http response
func NewResponse(resp *http.Response) *Response {
	// defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	return &Response{
		resp.StatusCode,
		&resp.Header,
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

func parseRequest(r *http.Request) (*url.URL, bool, []*Request) {
	var endPoint *url.URL
	endPointStr := r.FormValue("end_point")
	if endPointStr != "" {
		endPoint, _ = url.Parse(endPointStr)
	}
	includeHeadersStr := r.FormValue("include_headers")
	includeHeaders := includeHeadersStr != "false"
	batch := r.FormValue("batch")
	var requests []*Request
	log.Println(batch)
	err := json.Unmarshal([]byte(batch), &requests)
	if err != nil {
		log.Println("JSON unmarshal err:", err)
	}
	log.Println(requests[0])
	log.Println(requests[1])
	return endPoint, includeHeaders, requests
}

func serve() {
	ren := render.New()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		endPoint, includeHeaders, requests := parseRequest(r)
		log.Printf("end_point: %v, include_headers: %v, requests: %v\n", endPoint, includeHeaders, requests)
		responses := batchRequests(requests, endPoint)
		if !includeHeaders {
			for _, response := range responses {
				response.Headers = nil
			}
		}
		ren.JSON(w, http.StatusOK, responses)
	})
	http.ListenAndServe(":8080", nil)
}

func main() {
	// endPoint, _ := url.Parse("http://localhost:3000")
	// requests := []*Request{
	// 	&Request{
	// 		Method:      "GET",
	// 		RelativeURL: "api/v1/projects",
	// 	},
	// }
	// batchRequests(requests, endPoint)
	serve()
}
