package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Route struct {
	Route     string `json:"route"`
	Methods   []string
	Responses []Response `json:"responses"`
}

type Response struct {
	URI string `json:"uri"`
	Method string `json:"method"`
	StatusCode int      `json:"statusCode"`
	Headers map[string]string `json:"headers"`
	Body   map[string]interface{} `json:"body"`
}

var ResponseMapping = make(map[string]Response)

func loadRoutes(f string) []Route {
	log.Infof("Looking for routes.json file: %s", f)
	jsonFile, err := os.Open(f)
	if err != nil {
		log.Errorf("Unable to open file routes.json: [%s]", err)
		os.Exit(1)
	}
	defer jsonFile.Close()

	var routes []Route
	if err := json.NewDecoder(jsonFile).Decode(&routes); err != nil {
		log.Errorf("Unable to decode json object![%s]", err)
		os.Exit(2)
	}
	return routes
}

func (route *Route) createResponses() {
	log.Infof("%+v", route)
	for _, response := range route.Responses {
		key := fmt.Sprintf("%s|%s", response.URI, strings.ToUpper(response.Method))
		ResponseMapping[key] = response
	}
}

func (route *Route) routeHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("REQUEST: %+v %+v", r.Method, r.RequestURI)

	log.Infof("ResponseMapping: %+v", ResponseMapping)
	key := fmt.Sprintf("%s|%s", r.RequestURI, r.Method)
	response, ok := ResponseMapping[key]
	if !ok {
		log.Errorf("Response not mapped for method %s and URI %s", r.Method, r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "404 Response not mapped for method %s and URI %s", r.Method, r.RequestURI)
		return
	}

	log.Infof("RESPONSE: %+v", response)

	//write headers
	headers := response.Headers
	for k, v := range headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(response.StatusCode)

	if err := json.NewEncoder(w).Encode(response.Body); err != nil {
		log.Errorf("Unable to marshal body: %v", response.Body)
		os.Exit(5)
	}
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	jsonB, _ := json.Marshal(ResponseMapping)
	if _, err := w.Write(jsonB); err != nil {
		w.WriteHeader(500)
		log.Errorf("unable to list response mapping")
		w.Write([]byte("unable to list response mapping"))
	}
}

func NewMockify() {
	port, ok := os.LookupEnv("MOCKIFY_PORT")
	if !ok {
		port = "8001"
		log.Error(fmt.Sprintf("MOCKIFY_PORT not set; using default [%s]!", port))
	}
	var routes []Route
	routesFile, ok := os.LookupEnv("MOCKIFY_ROUTES")
	if !ok {
		log.Info("MOCKIFY_ROUTES not set.")
		path, err := os.Getwd()
		if err != nil {
			log.Errorf("unable to get working directory: [%s]", err)
			return
		}
		routes = loadRoutes(path + "/config/routes.json")
	} else {
		routes = loadRoutes(routesFile)
	}

	router := setupMockifyRouter(routes)

	log.Infof("%+v", ResponseMapping)
	log.Info("Ready on port " + port + "!")
	err := http.ListenAndServe("0.0.0.0:"+port, router)
	log.Error(err)
	os.Exit(6)
}

func setupMockifyRouter(routes []Route) *mux.Router {
	router := mux.NewRouter()

	//add builtin routes
	router.HandleFunc("/list", listHandler).Methods("GET")

	for _, route := range routes {
		route.createResponses()
		router.HandleFunc(route.Route, route.routeHandler).Methods(route.Methods...)
	}
	return router
}

func main() {
	NewMockify()
}
