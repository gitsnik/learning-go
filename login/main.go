package main

import (
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"math/rand"
	"net/http"
)

type PageData struct {
	RequestPath string
	User        string
	CsrfField   template.HTML
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/layout.html"))
	data := PageData{
		RequestPath: r.URL.Path,
	}
	tmpl.Execute(w, data)
}

const pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateCode(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}
	return string(bytes)
}

// Note: Don't store your key in your source code. Pass it via an
// environmental variable, or flag (or both), and don't accidentally commit it
// alongside your code. Ensure your key is sufficiently random - i.e. use Go's
// crypto/rand or securecookie.GenerateRandomKey(32) and persist the result.
// var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY"))
var store = sessions.NewCookieStore([]byte("abcdefghijklmnopqrstuvwxyz123456"))

func handleSecret(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "my-go-session")
	if session.Values["Authenticated"] != "Yes" {
		http.Redirect(w, r, "/login", 302)
	} else {
		tmpl := template.Must(template.ParseFiles("html/secret.html"))
		tmpl.Execute(w, nil)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("html/login.html"))
	data := PageData{
		CsrfField: csrf.TemplateField(r),
	}
	if r.Method == http.MethodPost {
		password, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
		if r.FormValue("username") == "testuser" && bcrypt.CompareHashAndPassword(password, []byte(r.FormValue("password"))) == nil {
			data.User = "Test"
			session, _ := store.Get(r, "my-go-session")
			session.Values["Authenticated"] = "Yes"
			session.Save(r, w)
		}
	}
	tmpl.Execute(w, data)
}

func securityHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none';")
		h.ServeHTTP(w, r)
	})
}

func setupHttpHandlers() *mux.Router {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(securityHeaders(http.StripPrefix("/static", fs)))

	csrfSecure := false
	csrfMiddleware := csrf.Protect(
		[]byte(generateCode(32)),
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.Secure(csrfSecure),
		csrf.FieldName("CSRFToken"),
	)

	r.Path("/login").Handler(securityHeaders(csrfMiddleware(http.HandlerFunc(handleLogin))))
	r.Path("/secret").Handler(securityHeaders(csrfMiddleware(http.HandlerFunc(handleSecret))))
	r.PathPrefix("/").Handler(securityHeaders(http.HandlerFunc(handleRoot)))
	return r
}

func main() {
	mux := setupHttpHandlers()
	http.ListenAndServe(":80", mux)
}
