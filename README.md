# Looking for contributors to help make mockify better

# mockify
Simple API mocking. No more waiting on backend teams to deliver services. Simply map the API call with a response and continue building great software.

**Update**
* Added `/list` endpoint (GET) to describe the current state of the `ResponseMapping`
* Added `/add` endpoint (POST)
* Added `/delete` endpoint (POST) to delete an existing mock by key (`/helloworld/foo|GET`)
* Added postman collection & environment
* Added functionality for `requestHeader`

## tl;dr

```
docker run -it -p 0.0.0.0:8001:8001 -v ~/Desktop/routes-cart.json:/app/routes.json  -e MOCKIFY_PORT=8001 brianmoran/mockify
curl localhost:8001/list
```

## Getting Started
These instructions will help you get started mocking your APIs with Docker.

1. Create a mapping file (JSON or YAML) anywhere you like. Mockify will check for the environment variable `MOCKIFY_ROUTES`. If the environment variable does not exist, Mockify will default to `./config/routes.yaml`
    * See [The Configuration Readme](https://github.com/brianmoran/mockify/tree/master/config) for more information and examples
1. _(Optional)_ Set the following environment variables `MOCKIFY_PORT` and `MOCKIFY_ROUTES`
1. Build the app inside a docker container by running `docker-compose up`. The docker container uses only 7MB of memory!
1. Start the docker container using a specific port and you can override routes.json as well
```
MOCKIFY_PORT=8002 docker-compose up # set a specific port
MOCKIFY_ROUTES=/app/routes-other.json docker-compose up # set a different routes file within the mockify folder as docker will copy it
```
or non-dockerized
```
go get github.com/gorilla/mux
go get github.com/json-iterator/go
go get gopkg.in/yaml.v2
go build -o main ./app/cmd/mockify.go
export MOCKIFY_PORT=8001
export MOCKIFY_ROUTES=~/Desktop/routes.json
./main
```

### Examples 

Use Postman, cURL, or your own microservice to connect to the mock API
```
curl http://localhost:8001/helloworld/bar
{"foo":{"key1":1,"key2":true,"key3":[{"bar":true,"baz":[1,2,"3"],"foo":"foo"}]}}
```

You can even mock errors
```
curl http://localhost:8001/helloworld/nonexisting
{"message": "Something bad happened but you knew that right?"}
```

---

Here is a postman collection that includes all internal calls as well as the tests that you see in the example configuration: [https://www.getpostman.com/collections/2daab06a399baa2c8576](https://www.getpostman.com/collections/2daab06a399baa2c8576)
