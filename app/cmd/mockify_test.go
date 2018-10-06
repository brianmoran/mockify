package main

import (
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
