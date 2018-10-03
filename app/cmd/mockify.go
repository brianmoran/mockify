package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Routes []Route `json:"routes"`
}

type Route struct {
	Path      string           `json:"path"`
	Methods   []string         `json:"methods"`
	Responses []ResponseConfig `json:"responses"`
}

type ResponseConfig struct {
	Methods []string `json:"methods"`
	URI     string   `json:"uri"`
	Get     Response `json:"get"`
	Post    Response `json:"post"`
	Put     Response `json:"put"`
	Delete  Response `json:"delete"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
	Headers    map[string]string
}

type ResponseContent struct {
	Message map[string]string
}

var responseMapping = make(map[string]Response)

func loadRoutes(f string) Config {
	log.Infof("Looking for routes.json file: %s", f)
	jsonFile, err := ioutil.ReadFile(f)
	if err != nil {
		log.Errorf("Unable to parse file routes.json")
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(jsonFile, &config)
	if err != nil {
		log.Errorf("Unable to unmarshall json objects!")
		os.Exit(2)
	}
	return config
}

func (route *Route) createResponses() {
	log.Infof("%+v", route)
	for _, response := range route.Responses {

		for _, method := range response.Methods {
			key := method + response.URI
			switch method {
			case "GET":
				responseMapping[key] = response.Get
			case "POST":
				responseMapping[key] = response.Post
			case "PUT":
				responseMapping[key] = response.Put
			case "DELETE":
				responseMapping[key] = response.Delete
			}
		}
	}
}

func (route *Route) routeHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("REQUEST: %+v %+v", r.Method, r.RequestURI)

	log.Infof("ResponseMapping: %+v", responseMapping)
	key := r.Method + r.RequestURI
	response, ok := responseMapping[key]
	if !ok {
		log.Errorf("Response not mapped for method %s and URI %s", r.Method, r.RequestURI)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(fmt.Sprintf("404 Response not mapped for method %s and URI %s", r.Method, r.RequestURI)))
		return
	}

	log.Infof("RESPONSE: %+v", response)

	for k, v := range response.Headers {
		w.Header().Add(k, v)
	}

	w.WriteHeader(response.StatusCode)

	var body []byte
	body, err := json.Marshal(response.Body)
	if err != nil {
		log.Errorf("Unable to marshall body: %s", response.Body)
		os.Exit(5)
	}

	output := string(body)
	for k, v := range mux.Vars(r) {
		kr := fmt.Sprintf("{%s}", k)
		output = strings.Replace(output, kr, v, -1)

		log.Infof("Replace: %+v by %+v", kr, v)
	}

	w.Write([]byte(output))
}

func NewMockify() {
	port, ok := os.LookupEnv("MOCKIFY_PORT")
	if !ok {
		log.Error("MOCKIFY_PORT not set!")
		port = "8001"
	}
	var config Config
	routesFile, ok := os.LookupEnv("MOCKIFY_ROUTES")
	if !ok {
		log.Info("MOCKIFY_ROUTES not set.")
		path, _ := os.Getwd()
		config = loadRoutes(path + "/config/routes.json")
	} else {
		config = loadRoutes(routesFile)
	}

	router := mux.NewRouter()
	for _, route := range config.Routes {
		route.createResponses()
		router.HandleFunc(route.Path, route.routeHandler).Methods(route.Methods...)
	}

	log.Infof("%+v", responseMapping)
	log.Info("Ready on port " + port + "!")
	err := http.ListenAndServe("0.0.0.0:"+port, router)
	log.Error(err)
	os.Exit(6)
}

func main() {
	NewMockify()
}
