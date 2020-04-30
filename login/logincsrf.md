## CSRF Tokens

When we last looked, we had a secure login form that was still requiring some work. Particularly, Zap! was reporting no Anti-CSRF Tokens:

```
WARN-NEW: Absence of Anti-CSRF Tokens [10202] x 2
        http://172.17.0.1/login (200 OK)
        http://172.17.0.1/login (200 OK)
FAIL-NEW: 0     FAIL-INPROG: 0  WARN-NEW: 1     WARN-INPROG: 0  INFO: 0 IGNORE: 0       PASS: 49
```

I will let you look into the process first, but just know that CSRF pretty much requires that you do a GET before a POST, as you will need to extract the CSRF token and deliver it along with your POST request. Fortunately, the gorilla/csrf module is particularly good at handling the edge cases.

Here's some of the supporting code for our program:

```
import (
    "math/rand"
)

const pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
// Generate a random code each time the program starts. This will cause issues with load balancing (in which case you might want to
// start using a central key value store to handle this part across your programs, and regenerate the key periodically
func generateCode(length int) string {
  bytes := make([]byte, length)
  for i := 0; i < length; i++ {
    bytes[i] = pool[rand.Intn(len(pool))]
  }
  return string(bytes)
}
```

The test is a bit tricker, but some of the key differences are here:

```
t.Run(tc.name, func(t *testing.T) {
  mux := http.NewServeMux()
  csrfKey := generateCode(32)
  p := csrf.Protect([]byte(csrfKey))(mux)

  var csrfToken string

  mux.Handle("/login", securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
          csrfToken = csrf.Token(r)
          handleLogin(w, r)
  })))

  ...

			if tc.method == "POST" {
				p.ServeHTTP(writer, request)

				request, _ = http.NewRequest(tc.method, tc.path, tc.data)
				request.Header.Set("X-CSRF-Token", csrfToken)
				request.Header.Set("Cookie", writer.Header().Get("Set-Cookie"))

				writer = httptest.NewRecorder()

				request, _ = http.NewRequest(tc.method, tc.path, tc.data)
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
			}
```

If we were to stop here we would have an interesting situation, because our tests are now actually more secure than our main program. Let's remedy that inside the setupHttpHandlers function, and add some extra handling of the CsrfField we will need for our form

```
type PageData struct {
    RequestPath string
    User        string
    CsrfField   template.Html
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    tmp := template.Must(template.ParseFiles("html/login.html"))
    data := PageData{
      CsrfField: csrf.TemplateField(r),
    }
...

func setupHttpHandlers() *mux.Router {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(securityHeaders(http.StripPrefix("/static", fs)))

    // Set to true when running behind TLS and dear god you should
    // be running behind TLS
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
```

Now we can add the CsrfField to our html/login.html form:

```
  Password: <input type='password' name='password' id='password' />
  {{ .CsrfField }}
  <br />
```

And that's it. Let's re-test against OWASP:

```
FAIL-NEW: 0     FAIL-INPROG: 0  WARN-NEW: 0     WARN-INPROG: 0  INFO: 0 IGNORE: 0       PASS: 50
```

Excellent
