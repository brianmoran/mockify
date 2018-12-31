package main

type Routes []struct {
	Route     string `json:"route"`
	Responses []struct {
		URI string `json:"uri"`
		Get struct {
			StatusCode int `json:"statusCode"`
			Response   map[string] interface {} `json:"response"`
		} `json:"get"`
	} `json:"responses"`
}