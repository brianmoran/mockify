package main

import (
	"testing"
)

func TestLoadRoutes(t *testing.T) {
	routes := loadRoutes("../../config/routes.yaml")
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
