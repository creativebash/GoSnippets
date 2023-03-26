package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCreateUser(t *testing.T) {
    // Create a test server
    ts := httptest.NewServer(http.HandlerFunc(handleCreateUser))
    defer ts.Close()

    // Create a new user
    newUser := &User{
        Username:  "testuser",
        Email:     "testuser@example.com",
        Firstname: "Test",
        Lastname:  "User",
        Sex:       "male",
    }

    // Marshal the new user to json
    jsonValue, _ := json.Marshal(newUser)

    // Create a new http request with the json payload in the request body
    req, err := http.NewRequest("POST", ts.URL, bytes.NewBuffer(jsonValue))
    if err != nil {
        t.Fatal(err)
    }
    req.Header.Set("Content-Type", "application/json")

    // Send the request to the test server
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        t.Fatal(err)
    }

    // Check the response status code
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
    }
}
