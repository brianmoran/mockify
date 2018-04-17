# mockify
API mocks for the minimalist. No more waiting on backend teams to deliver services. Simply map the API call with a response and continue building great software.

## Getting Started
These instructions will help you get started mocking your API's.
1. Modify routes.json to suit your needs. The path is loaded into [gorilla/mux](https://github.com/gorilla/mux) so you can use param matching.

*app/routes.json*
```
{
  "port": "7001",
  "routes": [
    {
      "path": "/helloworld/{key}",
      "methods": ["GET", "POST"],
      "responsePath": "app/response/helloworld/helloworld.json"
    },
    ...
  ]
}
```
2. Add response files/folders to the *app/responses* directory
Each response file should consist of a *list* of responses. The key is built off of the *method* and *URI*.
2. Build the app inside a docker container using the provided shell script `docker_install.sh`. The docker container uses only 7MB of memory!
2. Start the docker container using a specified port
```
docker run -it -p 0.0.0.0:8001:8001 -e PORT=8001 mockify
```
or non-dockerized
```
go build -o main ./app/cmd/mockify.go
export PORT=8001
./main
```
2. Use Postman, cURL, or your own microservice to connect to the mock API
```
curl -X GET \
  http://localhost:7001/helloworld/foo
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
### This is *very* basic and I am looking for suggestions and improvements