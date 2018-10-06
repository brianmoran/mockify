package main

import (
	"net/http/httptest"
	"testing"
)

func TestLoadRoutes(t *testing.T) {
	config := loadRoutes("../../config/routes.json")
	if len(config.Routes) == 0 {
		t.Error("At least 1 route is required")
		t.Fail()
	} else {
		for _, route := range config.Routes {
			if route.Path == "" {
				t.Error("Route is missing a responsePath")
				t.Fail()
			}
			if route.Path == "" {
				t.Error("Route is missing a path")
				t.Fail()
			}
			if len(route.Methods) == 0 {
				t.Error("Route needs at least 1 supported request method")
				t.Fail()
			}
		}
	}
}

func TestSimpleServer(t *testing.T) {
	config := loadRoutes("../../config/routes.json")
	router := setupMockifyRouter(config)

	req := httptest.NewRequest("GET", "/helloworld/bar", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	wantBody := "{\"message\":\"[GET] Hello bar!\"}\n" // json.NewEncoder adds a trailing \n
	gotBody := rec.Body.String()
	if gotBody != wantBody {
		t.Errorf("expected body [%s]; got [%s]", wantBody, gotBody)
		t.Fail()
	}

	wantStatusCode := 400
	gotStatusCode := rec.Result().StatusCode
	if gotStatusCode != wantStatusCode {
		t.Errorf("expected statusCode [%d]; got [%d]", wantStatusCode, gotStatusCode)
		t.Fail()
	}

	wantContentType := "application/json"
	gotContentType := rec.HeaderMap.Get("Content-Type")
	if gotContentType != wantContentType {
		t.Errorf("expected content type [%s]; got [%s]", wantContentType, gotContentType)
		t.Fail()
	}
}
