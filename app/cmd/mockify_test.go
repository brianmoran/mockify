package main

import (
	"testing"
	"os"
)

func TestLoadConfig(t *testing.T) {
	// Move up (2) dirs
	os.Chdir("../../")
	path, _ := os.Getwd()
	config := loadConfig(path)
	if config.Port == "" {
		t.Error("A port is required")
		t.Fail()
	}
	if len(config.Routes) == 0 {
		t.Error("At least 1 route is required")
		t.Fail()
	} else {
		for _, route := range config.Routes {
			if route.ResponsePath == "" {
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