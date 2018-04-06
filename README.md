# mockify
API mocks for the minimalist

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
3. Build the app inside a docker container (7MB)
4. Start the docker container `docker run -it -p 0.0.0.0:7001:7001 mockify`
4. Use Postman or your own microservice to connect to the mock API
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

### This is *very* basic and I am looking for suggestions and improvements