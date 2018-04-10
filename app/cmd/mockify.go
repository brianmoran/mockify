package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
)

type Route struct {
	Path         string   `json:"path"`
	Methods      []string `json:"methods"`
	ResponsePath string   `json:"responsePath"`
}

type Config struct {
	Port   string  `json:"port"`
	Routes []Route `json:"routes"`
}

type Response struct {
	Method     string      `json:"method"`
	URI        string      `json:"uri"`
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
	Headers    map[string]string
}

type Responses struct {
	Responses []Response `json:"responses"`
}

var responseMapping = make(map[string]Response)

func loadConfig(path string) Config {
	log.Infof("Looking for routes.json file in  %s/app directory", path)
	jsonFile, err := ioutil.ReadFile("app/routes.json")
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

func (route *Route) prefetchResponses() {
	rawData, err := ioutil.ReadFile(route.ResponsePath)
	if err != nil {
		log.Errorf("Unable to open response file %s", route.ResponsePath)
		os.Exit(3)
	}

	var responses Responses
	json.Unmarshal(rawData, &responses)

	for _, r := range responses.Responses {
		key := r.Method + r.URI
		responseMapping[key] = r
	}
}

func (route *Route) routeHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("REQUEST: %+v %+v", r.Method, r.RequestURI)

	log.Infof("ResponseMapping: %+v", responseMapping)
	key := r.Method + r.RequestURI
	response, ok := responseMapping[key]
	if !ok {
		log.Errorf("Response not mapped for method %s and URI %s", r.Method, r.RequestURI)
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
	w.Write(body)
}

func NewMockify() {
	path, _ := os.Getwd()
	config := loadConfig(path)

	router := mux.NewRouter()
	for _, route := range config.Routes {
		route.prefetchResponses()
		router.HandleFunc(route.Path, route.routeHandler).Methods(route.Methods...)
	}

	err := http.ListenAndServe("0.0.0.0:"+config.Port, router)
	log.Error(err)
	os.Exit(6)
}

func main() {
	NewMockify()
}
