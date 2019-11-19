package main

import (
    "fmt"
    "net/http"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}

func main() {
    http.HandleFunc("/", handleRoot)

    http.ListenAndServe(":80", nil)
}
