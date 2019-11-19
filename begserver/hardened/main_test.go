package main

import (
  "io/ioutil"
  "net/http"
  "net/http/httptest"
  "strings"
  "testing"
)

func TestHandleRoot(t *testing.T) {
  mux := http.NewServeMux()
  mux.HandleFunc("/", handleRoot)
  writer := httptest.NewRecorder()

  request, _ := http.NewRequest("GET", "/", nil)
  mux.ServeHTTP(writer, request)

  if writer.Code != 200 {
    t.Errorf("Response code %v expected 200", writer.Code)
  }

  if response, err := ioutil.ReadAll(writer.Body); err != nil {
    t.Fail()
  } else {
    if !strings.Contains(string(response), "Hello, you've requested: /") {
      t.Errorf("Response does not contain `Hello, you've requested: /`: %v", string(response))
    }

    if clacks := writer.Header().Get("X-Clacks-Overhead"); clacks != "GNU Terry Pratchett" {
      t.Errorf("Require X-Clacks-Overhead: GNU Terry Pratchett. Never forget.")
    }

    if contopts := writer.Header().Get("X-Content-Type-Options"); contopts != "nosniff" {
      t.Errorf("Require X-Content-Type-Options: nosniff")
    }
  }
}
