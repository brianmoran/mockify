# mockify
Simple API mocking. No more waiting on backend teams to deliver services. Simply map the API call with a response and 
continue building great software.

**Update**
* Added */list* endpoint (GET) to describe the current state of the `ResponseMapping`
* Added `/add` endpoint (POST)
* Added postman collection & environment

## tl;dr
Pull and run the [image](https://hub.docker.com/r/brianmoran/mockify/) from Dockerhub
```
docker pull brianmoran/mockify
docker run -it -p 0.0.0.0:8001:8001 -v ~/Desktop/routes-cart.json:/app/routes.json  -e MOCKIFY_PORT=8001 mockify
```

## Getting Started
These instructions will help you get started mocking your API's.
1. Create a mapping file (json) anywhere you'd like. Mockify will check for the variable **MOCKIFY_ROUTES**. If the environment variable does not exist, then Mockify will default to **app/routes.json**

*Example configuration file*
```
[
  {
    "route": "/helloworld/{key}",
    "methods": ["get"],
    "responses": [
      {
        "uri": "/helloworld/foo",
        "method": "get",
        "statusCode": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "body": {
          "message": "Welcome to Mockify!"
        }
      },
      {
        "uri": "/helloworld/bar",
        "method": "get",
        "statusCode": 200,
        "headers": {
          "Content-Type": "application/json"
        },
        "body": {
          "foo": {
            "key1": 1,
            "key2": true,
            "key3": [
              {
                "foo": "foo",
                "bar": true,
                "baz": [1, 2, "3"]
              }
            ]
          }
        }
      }
    ]
  }
]
2. (Optional) Export the following variables **MOCKIFY_PORT** and **MOCKIFY_ROUTES**

2. Build the app inside a docker container by running `docker-compose up`. The docker container uses only 7MB of memory!
2. Start the docker container using a specific port and you can override routes.json as well
```
MOCKIFY_PORT=8002 docker-compose up # set a specific port
MOCKIFY_ROUTES=/app/routes-other.js # set a different routes file within the mockify folder as docker will copy it
```
or non-dockerized
```
go build -o main ./app/cmd/mockify.go
export MOCKIFY_PORT=8001
export MOCKIFY_ROUTES=~/Desktop/routes.json
./main
```
2. Use Postman, cURL, or your own microservice to connect to the mock API
```
curl -X GET \
  http://localhost:8001/helloworld/foo
```
```
{
    "message": "Hello world!",
    "misc": [
        "foo",
        "bar",
        "baz"
    ]
}
```
3. You can even mock errors
```
curl -X GET \
  http://localhost:7001/helloworld/bar
{
    "message": "Something bad happened but you knew that right?"
}
```
