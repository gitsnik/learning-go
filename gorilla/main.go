package main

import (
    "net/http"
    "github.com/gorilla/mux"
)

func handleRoot(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello, you've requested: " + r.URL.Path + "\n"))
}

func securityHeaders(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    h.ServeHTTP(w, r)
  })
}

func setupHttpHandlers() *mux.Router {
  r := mux.NewRouter()

  fs := http.FileServer(http.Dir("static/"))
  r.PathPrefix("/static/").Handler(securityHeaders(http.StripPrefix("/static", fs)))

  r.PathPrefix("/").Handler(securityHeaders(http.HandlerFunc(handleRoot)))
  return r
}

func main() {
  mux := setupHttpHandlers()
  http.ListenAndServe(":80", mux)
}
