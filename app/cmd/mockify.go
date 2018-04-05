package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"os"
	log "github.com/sirupsen/logrus"
)

type Route struct {
	Path     string `json:"path"`
	Methods   []string `json:"methods"`
	ResponsePath string `json:"responsePath"`
}

type Config struct {
	Port string `json:"port"`
	Routes []Route `json:"routes"`
}

type Response struct {
	StatusCode int `json:"statusCode"`
	Body interface{} `json:"body"`
	Headers map[string]string
}

func main() {
	wd, _ := os.Getwd()
	log.Infof("Looking for routes.json file in  %s/app directory\n", wd)
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

	router := mux.NewRouter()
	for _, route := range config.Routes {
		router.HandleFunc(route.Path, func(w http.ResponseWriter, r *http.Request) {
			rawData, err := ioutil.ReadFile(route.ResponsePath)
			if err != nil {
				log.Errorf("Unable to open response file %s", route.ResponsePath)
				os.Exit(3)
			}
			var response = Response{}
			json.Unmarshal(rawData, &response)
			log.Infof("%+v\n", response)

			for k, v := range response.Headers {
				w.Header().Add(k, v)
			}

			w.WriteHeader(response.StatusCode)

			var body []byte
			body, err = json.Marshal(response.Body)
			if err!= nil {
				log.Errorf("Unable to marshall body: %s", response.Body)
				os.Exit(5)
			}
			w.Write(body)
		}).Methods(route.Methods...)
	}

	err = http.ListenAndServe("0.0.0.0:"+config.Port, router)
	log.Error(err)
	os.Exit(6)
}