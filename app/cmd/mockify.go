package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	URI         string            `yaml:"uri" json:"uri"`
	Method      string            `yaml:"method" json:"method"`
	RequestBody string            `yaml:"requestBody" json:"requestBody"`
	StatusCode  int               `yaml:"statusCode" json:"statusCode"`
	Headers     map[string]string `yaml:"headers" json:"headers"`
	Body        interface{}       `yaml:"body" json:"body"`
}

var ResponseMapping = make(map[string]Response)
var Router = mux.NewRouter()

func printError(format string, v ...interface{}) {
	log.SetPrefix("[ERROR] ")
	log.Printf(format, v...)
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
	log.Printf("%+v", route)
	for _, response := range route.Responses {
		key := fmt.Sprintf("%s|%s|%s", response.URI, strings.ToUpper(response.Method), strings.ToUpper(response.RequestBody))
		ResponseMapping[key] = response
	}
}

func getResponse(method, uri, body string) *Response {
	for _, response := range ResponseMapping {
		if uri == response.URI && method == response.Method {
			if response.RequestBody == "" {
				return &response
			} else if strings.Contains(body, response.RequestBody) {
				return &response
			}
		}
	}
	return nil
}

func (route *Route) routeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("REQUEST: %+v %+v", r.Method, r.RequestURI)
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		bodyBytes = []byte("")
	}

	response := getResponse(r.Method, r.RequestURI, string(bodyBytes))
	if response == nil {
		printError("Response not mapped for method %s and URI %s", r.Method, r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Response not mapped for method %s and URI %s", r.Method, r.RequestURI)
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
		w.Write([]byte(response.Body.(string)))
	} else {
		var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary
		jsonB, err := jsonx.Marshal(response.Body)
		if err != nil {
			printError("Response could not be converted to JSON")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "500 Response could not be converted to JSON")
			return
		}
		w.Write(jsonB)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	var jsonx = jsoniter.ConfigCompatibleWithStandardLibrary

	w.Header().Add("Content-Type", "application/json")

	//jsonB, err := json.Marshal(ResponseMapping)
	jsonB, err := jsonx.Marshal(ResponseMapping)
	if err != nil {
		fmt.Println("Error", err.Error())
		w.WriteHeader(500)
		printError("unable to list response mapping")
		w.Write([]byte("unable to list response mapping"))
	}
	if _, err := w.Write(jsonB); err != nil {
		w.WriteHeader(500)
		printError("unable to list response mapping")
		w.Write([]byte("unable to list response mapping"))
	}
}

func addHandler(w http.ResponseWriter, r *http.Request) {
	var errString string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errString = "unable to parse request body"
		printError(errString)
		w.WriteHeader(500)
		w.Write([]byte(errString))
	}

	route := Route{}
	err = json.Unmarshal(body, &route)
	if err != nil {
		errString = "unable to unmarshal body"
		printError(errString)
		w.WriteHeader(500)
		w.Write([]byte(errString))
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
		w.Write([]byte(errString))
	}

	if body == nil {
		errString = "body is empty"
		printError(errString)
		w.WriteHeader(500)
		w.Write([]byte(errString))
	}
	key := string(body)

	_, ok := ResponseMapping[key]
	if !ok {
		log.Printf("key: %s doesn't exist", key)
		w.WriteHeader(200)
		w.Write([]byte("nothing to delete"))
		return
	}

	delete(ResponseMapping, key)
	w.WriteHeader(200)
	w.Write([]byte("mock deleted"))
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

	log.Printf("%+v", ResponseMapping)
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
		os.Setenv("MOCKIFY_ROUTES", "./config/routes.yaml")
	}
	NewMockify()
}
