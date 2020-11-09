package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/gorilla/mux"

	jsoniter "github.com/json-iterator/go"
)

type Route struct {
	Route     string `yaml:"route" json:"route"`
	Methods   []string
	Responses []Response `yaml:"responses" json:"responses"`
}

type Response struct {
	URI           string            `yaml:"uri" json:"uri"`
	Method        string            `yaml:"method" json:"method"`
	RequestHeader string            `yaml:"requestHeader" json:"requestHeader"`
	RequestBody   string            `yaml:"requestBody" json:"requestBody"`
	StatusCode    int               `yaml:"statusCode" json:"statusCode"`
	Headers       map[string]string `yaml:"headers" json:"headers"`
	Body          interface{}       `yaml:"body" json:"body"`
}

// For more how the different response mappings are used, read the README-file in the config folder
var (
	RequestBodyResponseMappings    = make(map[string]Response)
	RequestHeaderResponseMappings  = make(map[string]Response)
	LowestPriorityResponseMappings = make(map[string]Response)
	Router                         = mux.NewRouter()
)

const (
	invalidHeaderChars = "[:\\s]*"
)

func printError(format string, v ...interface{}) {
	log.SetPrefix("[ERROR] ")
	log.Printf(format, v...)
	log.SetPrefix("")
}

func printRegisteredResponse(key string, response Response) {
	log.SetPrefix("[RESPONSE] ")
	log.Printf("[%s] -> %+v", key, response)
	log.SetPrefix("")
}

func loadRoutes(f string) []Route {
	log.Printf("Looking for routes in file: %s", f)

	routes := make([]Route, 0)

	yamlFile, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &routes)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return routes
}

func (route *Route) createResponses() {
	headerRegEx := regexp.MustCompile(invalidHeaderChars)

	for _, response := range route.Responses {
		upperMethod := strings.ToUpper(response.Method)
		upperURI := strings.ToUpper(response.URI)

		addToLowestPriority := true

		if response.RequestBody != "" {
			upperBody := strings.ToUpper(response.RequestBody)
			key := fmt.Sprintf("%s|%s|%s|BODY", upperURI, upperMethod, upperBody)
			printRegisteredResponse(key, response)
			RequestBodyResponseMappings[key] = response
			addToLowestPriority = false
		}

		// Add to second priority if we have a requestHeader
		// A route can have both a requestHeader and a requestBody, and these will be added to both mapping slices
		if response.RequestHeader != "" {
			upperHeader := strings.ToUpper(response.RequestHeader)
			upperHeader = strings.TrimSpace(upperHeader)
			upperHeader = headerRegEx.ReplaceAllString(upperHeader, "")
			key := fmt.Sprintf("%s|%s|%s|HEADER", upperURI, upperMethod, upperHeader)
			printRegisteredResponse(key, response)
			RequestHeaderResponseMappings[key] = response
			addToLowestPriority = false
		}

		if addToLowestPriority {
			key := fmt.Sprintf("%s|%s|LOWEST", upperURI, upperMethod)
			printRegisteredResponse(key, response)
			LowestPriorityResponseMappings[key] = response
		}
	}
}

func getResponse(method, uri, body string, requestHeaders http.Header) *Response {
	// Check if we have a match in the highest priority splice
	for _, response := range RequestBodyResponseMappings {
		if uri == response.URI && method == response.Method && strings.Contains(body, response.RequestBody) {
			return &response
		}
	}

	// Check if we have a match in the second highest priority splice
	for _, response := range RequestHeaderResponseMappings {
		if uri == response.URI && method == response.Method && response.RequestHeader != "" {
			suppliedMatchHeader := strings.Split(response.RequestHeader, ":")
			if len(suppliedMatchHeader) > 2 {
				log.Fatalf(`Wrongfully use of requestHeader! Should only be "key: value", you had: '%v'!`, suppliedMatchHeader)
			}
			suppliedMatchHeader[0] = strings.TrimSpace(suppliedMatchHeader[0])
			suppliedMatchHeader[1] = strings.TrimSpace(suppliedMatchHeader[1])
			found := requestHeaders.Get(suppliedMatchHeader[0])
			if found != "" && found == suppliedMatchHeader[1] {
				return &response
			}
		}
	}

	// Check if we have a match in the lowest priority splice
	for _, response := range LowestPriorityResponseMappings {
		if uri == response.URI && method == response.Method {
			return &response
		}
	}

	return nil
}

func (route *Route) routeHandler(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		bodyBytes = []byte("")
	}
	log.Printf("REQUEST: %+v %+v [%+v] [%+v]", r.Method, r.RequestURI, string(bodyBytes), r.Header)

	response := getResponse(r.Method, r.RequestURI, string(bodyBytes), r.Header)
	if response == nil {
		printError("Response not mapped for method %s and URI %s", r.Method, r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
		_, _ = fmt.Fprintf(w, "404 Response not mapped for method %s and URI %s", r.Method, r.RequestURI)
		return
	}

	log.Printf("RESPONSE: %+v", response)

	//write headers
	for k, v := range response.Headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(response.StatusCode)
	isJson := false
	_, ok := response.Headers["Content-Type"]
	if ok && response.Headers["Content-Type"] == "application/json" {
		isJson = true
	}
	if !isJson {
		_, _ = w.Write([]byte(response.Body.(string)))
	} else {
		var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary
		jsonB, err := jsonx.Marshal(response.Body)
		if err != nil {
			printError("Response could not be converted to JSON")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = fmt.Fprint(w, "500 Response could not be converted to JSON")
			return
		}
		_, _ = w.Write(jsonB)
	}
}

func listHandler(w http.ResponseWriter, _ *http.Request) {
	var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary

	w.Header().Add("Content-Type", "application/json")

	jsonB, err := jsonx.Marshal(RequestBodyResponseMappings)
	if err != nil {
		fmt.Println("Error", err.Error())
		w.WriteHeader(500)
		printError("unable to list response mapping")
		_, _ = w.Write([]byte("unable to list response mapping"))
	}
	if _, err := w.Write(jsonB); err != nil {
		w.WriteHeader(500)
		printError("unable to list response mapping")
		_, _ = w.Write([]byte("unable to list response mapping"))
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	var errString string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errString = "unable to parse request body"
		printError(errString)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(errString))
	}

	route := Route{}
	err = json.Unmarshal(body, &route)
	if err != nil {
		errString = "unable to unmarshal body"
		printError(errString)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(errString))
	}

	route.createResponses()
	Router.HandleFunc(route.Route, route.routeHandler).Methods(route.Methods...)

}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	var errString string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errString = "unable to parse request body"
		printError(errString)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(errString))
	}

	if body == nil {
		errString = "body is empty"
		printError(errString)
		w.WriteHeader(500)
		_, _ = w.Write([]byte(errString))
	}
	key := string(body)

	_, ok := RequestBodyResponseMappings[key]
	if !ok {
		log.Printf("key: %s doesn't exist", key)
		w.WriteHeader(200)
		_, _ = w.Write([]byte("nothing to delete"))
		return
	}

	delete(RequestBodyResponseMappings, key)
	w.WriteHeader(200)
	_, _ = w.Write([]byte("mock deleted"))
}

// NewMockify sets up a new instance/http server
func NewMockify() {
	port, ok := os.LookupEnv("MOCKIFY_PORT")
	if !ok {
		port = "8001"
		printError("MOCKIFY_PORT not set; using default [%s]!", port)
	}
	var routes []Route
	routesFile, ok := os.LookupEnv("MOCKIFY_ROUTES")
	if !ok {
		log.Print("MOCKIFY_ROUTES not set.")
		os.Exit(1)
	} else {
		routes = loadRoutes(routesFile)
	}

	setupMockifyRouter(routes)

	log.Printf("Ready on port %s!", port)
	if err := http.ListenAndServe("0.0.0.0:"+port, Router); err != nil {
		log.Fatal(err)
	}
	os.Exit(6)
}

func setupMockifyRouter(routes []Route) {
	//add builtin routes
	Router.HandleFunc("/list", listHandler).Methods(http.MethodGet)
	Router.HandleFunc("/add", addHandler).Methods(http.MethodPost)
	Router.HandleFunc("/delete", deleteHandler).Methods(http.MethodPost)

	for _, route := range routes {
		route.createResponses()
		Router.HandleFunc(route.Route, route.routeHandler).Methods(route.Methods...)
	}
}

func main() {
	path, exist := os.LookupEnv("MOCKIFY_ROUTES")
	if exist {
		log.Printf(fmt.Sprintf("MOCKIFY_ROUTES set. [%s]", path))
	} else {
		log.Printf(fmt.Sprintf("MOCKIFY_ROUTES not set. Default ./config/routes.yaml"))
		_ = os.Setenv("MOCKIFY_ROUTES", "./config/routes.yaml")
	}
	NewMockify()
}
