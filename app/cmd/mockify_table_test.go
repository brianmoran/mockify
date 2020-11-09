package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockifyTestStructure struct {
	name                 string
	description          string
	configFileName       string
	requestBody          *bytes.Buffer
	requestMethod        string
	requestPath          string
	expectedStatusCode   int
	expectedResponseBody string
	expectedContentType  string
	setup                func(*http.Request, *httptest.ResponseRecorder)
}

type Tests struct {
	Tests []mockifyTestStructure
	T     *testing.T
}

func TestMockify(t *testing.T) {
	tests := []mockifyTestStructure{
		{
			name:                 "RequestBody_shouldChooseCorrectRoute",
			description:          "Verify that the correct route where the 'def' requestBody has been specified is chosen",
			requestBody:          bytes.NewBufferString(`{"type":"def"}`),
			requestMethod:        "POST",
			requestPath:          "/api/mcp",
			expectedStatusCode:   http.StatusCreated,
			expectedContentType:  "application/json",
			expectedResponseBody: `{"foo":{"key1":1,"key2":true,"key3":[{"bar":true,"baz":[1,2,"3"],"foo":"foo"}]}}`,
		},
		{
			name:                 "RequestBody_requestBodyShouldHaveHigherPriorityThanRequestHeader",
			description:          "Verify that the correct route where the 'def' requestBody has been specified is chosen, even when a requestHeader is supplied",
			requestBody:          bytes.NewBufferString(`{"type":"def"}`),
			requestMethod:        "POST",
			requestPath:          "/api/mcp",
			expectedStatusCode:   http.StatusCreated,
			expectedContentType:  "application/json",
			expectedResponseBody: `{"foo":{"key1":1,"key2":true,"key3":[{"bar":true,"baz":[1,2,"3"],"foo":"foo"}]}}`,
			setup: func(request *http.Request, recorder *httptest.ResponseRecorder) {
				request.Header.Set("Authorization", "foo-bar")
			},
		},
		{
			name:                 "RequestBody_shouldWorkWithBothBodyAndHeader",
			description:          "Verify that a route with both requestBody and requestHeader works",
			requestBody:          bytes.NewBufferString(`{"type":"body-have-higher-priority-over-header"}`),
			requestMethod:        "POST",
			requestPath:          "/api/mcp",
			expectedStatusCode:   http.StatusFound,
			expectedContentType:  "application/json",
			expectedResponseBody: `{"win":{"key1":1,"key2":true}}`,
			setup: func(request *http.Request, recorder *httptest.ResponseRecorder) {
				request.Header.Set("Authorization", "foo-bar")
			},
		},
		{
			name:                 "RequestHeader_shouldChooseCorrectRoute",
			description:          "Verify that the correct route where the 'Authorization: foo-bar' requestHeader has been specified is chosen",
			requestMethod:        "POST",
			requestPath:          "/api/mcp",
			expectedStatusCode:   http.StatusFound,
			expectedContentType:  "application/json",
			expectedResponseBody: `{"win":{"key1":1,"key2":true}}`,
			setup: func(request *http.Request, recorder *httptest.ResponseRecorder) {
				request.Header.Set("Authorization", "foo-bar")
				request.Header.Set("Foo", "Bar")
			},
		},
	}

	tt := Tests{
		Tests: tests,
		T:     t,
	}
	tt.runTests()
}

func TestMockifyDifferentConfigFiles(t *testing.T) {
	tests := []mockifyTestStructure{
		{
			name:                 "DefaultConfigFileInYAML",
			description:          "Test with the default YAML config file",
			requestMethod:        "GET",
			requestPath:          "/helloworld/foo",
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Welcome to Mockify!"}`,
			expectedContentType:  "application/json",
		},
		{
			name:                 "ConfigFileInJSON",
			description:          "Test with a configuration file in JSON instead of default YAML",
			configFileName:       "../../config/routes.json",
			requestMethod:        "GET",
			requestPath:          "/helloworld/foo",
			expectedStatusCode:   200,
			expectedResponseBody: `{"message":"Welcome to Mockify!"}`,
			expectedContentType:  "application/json",
		},
	}

	tt := Tests{
		Tests: tests,
		T:     t,
	}
	tt.runTests()
}

func (impl Tests) runTests() {
	for _, test := range impl.Tests {
		impl.T.Run(test.name, func(t *testing.T) {
			t.Log(test.description)

			var config []Route
			if test.configFileName == "" {
				config = loadRoutes("../../config/routes.yaml")
			} else {
				config = loadRoutes(test.configFileName)
			}
			setupMockifyRouter(config)

			req := httptest.NewRequest(test.requestMethod, test.requestPath, nil)
			if test.requestBody != nil {
				req = httptest.NewRequest(test.requestMethod, test.requestPath, test.requestBody)
			}
			rec := httptest.NewRecorder()

			if test.setup != nil {
				test.setup(req, rec)
			}

			Router.ServeHTTP(rec, req)

			gotBody := rec.Body.String()
			if gotBody != test.expectedResponseBody {
				t.Errorf(`expected body "%s"; got "%s""`, test.expectedResponseBody, gotBody)
				t.Fail()
			}

			gotStatusCode := rec.Result().StatusCode
			if gotStatusCode != test.expectedStatusCode {
				t.Errorf(`expected status code "%d"; got "%d"`, test.expectedStatusCode, gotStatusCode)
				t.Fail()
			}

			gotContentType := rec.Header().Get("Content-Type")
			if gotContentType != test.expectedContentType {
				t.Errorf(`expected content type "%s"; got "%s"`, test.expectedContentType, gotContentType)
				t.Fail()
			}
		})
	}
}
