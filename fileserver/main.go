package main

import (
    "net/http"
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

func setupHttpHandlers() *http.ServeMux {
  mux := http.NewServeMux()
  mux.Handle("/", securityHeaders(http.HandlerFunc(handleRoot)))

  fs := http.FileServer(http.Dir("static/"))
  mux.Handle("/static/", securityHeaders(http.StripPrefix("/static", fs)))
  return mux
}

func main() {
  mux := setupHttpHandlers()
  http.ListenAndServe(":80", mux)
}
