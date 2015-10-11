package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

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

func showError(err error) {
	if err != nil {
		log.Println(err)
	}
}

// NewResponse creates response instance from http response
func NewResponse(resp *http.Response) *Response {
	if resp == nil {
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	showError(err)
	return &Response{
		resp.StatusCode,
		&resp.Header,
		string(body),
	}
}

func batchRequests(requests []*Request, endPoint *url.URL) []*Response {
	responses := make([]*Response, len(requests))
	var wg sync.WaitGroup
	// TODO: change to use go rutine
	for i, request := range requests {
		wg.Add(1)
		go func(i int, request *Request) {
			client := &http.Client{
				Timeout: 10 * time.Second,
			}
			log.Println("Resuest:", request.Method, request.RelativeURL)
			url, err := endPoint.Parse(request.RelativeURL)
			showError(err)
			req, err := http.NewRequest(request.Method, url.String(), strings.NewReader(request.Body))
			showError(err)
			resp, err := client.Do(req)
			showError(err)
			log.Println(resp)
			responses[i] = NewResponse(resp)
			wg.Done()
		}(i, request)
	}
	wg.Wait()
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
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)

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
