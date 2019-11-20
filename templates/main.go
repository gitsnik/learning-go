package main

import (
    "net/http"
    "html/template"
    "github.com/gorilla/mux"
)

type PageData struct {
  RequestPath string
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
  tmpl := template.Must(template.ParseFiles("html/layout.html"))
  data := PageData{
    RequestPath: r.URL.Path,
  }
  tmpl.Execute(w, data)
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
