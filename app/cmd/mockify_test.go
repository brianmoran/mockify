package main

import (
	"net/http/httptest"
	"testing"
)

func TestLoadRoutes(t *testing.T) {
	routes := loadRoutes("../../config/routes.json")
	if len(routes) == 0 {
		t.Error("at least 1 route is required")
		t.Fail()
	} else {
		for _, route := range routes {
			if route.Route == "" {
				t.Error("missing a route")
				t.Fail()
			}
			if len(route.Methods) == 0 {
				t.Error("route needs at least 1 supported request method")
				t.Fail()
			}
		}
	}
}

func TestSimpleServer(t *testing.T) {
	config := loadRoutes("../../config/routes.json")
	setupMockifyRouter(config)

	req := httptest.NewRequest("GET", "/helloworld/foo", nil)
	rec := httptest.NewRecorder()

	Router.ServeHTTP(rec, req)

	wantBody := "{\"message\":\"Welcome to Mockify!\"}\n" // json.NewEncoder adds a trailing \n
	gotBody := rec.Body.String()
	if gotBody != wantBody {
		t.Errorf("expected body [%s]; got [%s]", wantBody, gotBody)
		t.Fail()
	}

	wantStatusCode := 200
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
