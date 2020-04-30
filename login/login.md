## Login Form

Let's add a simple login form to our application.

We will handle:

/login
/secret

Access to secret without being logged in, will redirect us to the /login window. Access to anything else will work as normal.

As always, we are building on top of the latest project (in this case, templates)

Before we begin, let's re-test with OWASP:

```
FAIL-NEW: X-Frame-Options Header Not Set [10020] x 4
    http://172.17.0.1/ (200 OK)
    http://172.17.0.1/robots.txt (200 OK)
    http://172.17.0.1/sitemap.xml (200 OK)
    http://172.17.0.1 (200 OK)
FAIL-NEW: Content Security Policy (CSP) Header Not Set [10038] x 4
    http://172.17.0.1/ (200 OK)
    http://172.17.0.1/robots.txt (200 OK)
    http://172.17.0.1/sitemap.xml (200 OK)
    http://172.17.0.1 (200 OK)
FAIL-NEW: 2     FAIL-INPROG: 0  WARN-NEW: 0     WARN-INPROG: 0  INFO: 0 IGNORE: 0       PASS: 48
```

Oops. Ok, we can fix this. They're both headers, and we've already started to test for those. We will convert our test function to be table driven, and add our two new headers:

```
func testSecurityHeaders(t *testing.T, writer *httptest.ResponseRecorder) {
	cases := []struct {
		header   string
		expected string
	}{
		{
			header:   "X-Clacks-Overhead",
			expected: "GNU Terry Pratchett",
		},
		{
			header:   "X-Content-Type-Options",
			expected: "nosniff",
		},
		{
			header:   "X-Frame-Options",
			expected: "deny",
		},
		{
			header:   "Content-Security-Policy",
			expected: "default-src 'self'; frame-ancestors 'none';",
		},
	}

	for _, tc := range cases {
		t.Logf("Security - HTTP Header - %s", tc.header)
		got := writer.Header().Get(tc.header)
		if got != tc.expected {
			t.Errorf("[FAIL] Require %s to return '%s'. Got '%v'", tc.header, tc.expected, got)
		}
	}
}
```

The test driven pattern asks for Red, Green, Refactor - so:

Red:
```
--- FAIL: TestHandleRoot (0.00s)
  main_test.go:35: Security - HTTP Header - X-Clacks-Overhead
  main_test.go:35: Security - HTTP Header - X-Content-Type-Options
  main_test.go:35: Security - HTTP Header - X-Frame-Options
  main_test.go:38: [FAIL] Require X-Frame-Options to return 'deny'. Got ''
  main_test.go:35: Security - HTTP Header - Content-Security-Policy
  main_test.go:38: [FAIL] Require Content-Security-Policy to return 'default-src 'self'; frame-ancestors 'none';'. Got ''
--- FAIL: TestSetupHttpHandlers (0.01s)
  main_test.go:35: Security - HTTP Header - X-Clacks-Overhead
  main_test.go:35: Security - HTTP Header - X-Content-Type-Options
  main_test.go:35: Security - HTTP Header - X-Frame-Options
  main_test.go:38: [FAIL] Require X-Frame-Options to return 'deny'. Got ''
  main_test.go:35: Security - HTTP Header - Content-Security-Policy
  main_test.go:38: [FAIL] Require Content-Security-Policy to return 'default-src 'self'; frame-ancestors 'none';'. Got ''
FAIL
exit status 1
```

Green:
```
func securityHeaders(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
    w.Header().Set("X-Content-Type-Options", "nosniff")
    w.Header().Set("X-Frame-Options", "deny")
    w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none';")
    h.ServeHTTP(w, r)
  })
}
```

Refactor:

We're not going to refactor this at the moment. We could refactor this to a table driven design for ```w.Header().Set()``` if we so choose, but it is probably not necessary at this time.

We're going to place two links on the base of our template so that Zap! can find the pages, and we will start to fill the information in later. First, we need to be sure that the Zap! tests pass:

```
FAIL-NEW: 0     FAIL-INPROG: 0  WARN-NEW: 0     WARN-INPROG: 0  INFO: 0 IGNORE: 0       PASS: 50
```

Excellent, let's get started. Here are the tests we need:

* Access to /secret redirects to /login when not logged in
* Access to /secret shows "You found it!" when logged in
* GET /login shows a username and password form
* POST /login authenticates a user and logs them in
* OWASP Zap! Continues to pass
* Passwords are "stored" securely

We'll write some simple tests for this - can we receive and set a cookie, can we login, etc.

```
func TestSecret(t *testing.T) {
	cases := []struct {
		name       string
		path       string
		code       int
		redirect   string
		method     string
		expected   string
		data       *strings.Reader
		usecookie  bool
		savecookie bool
	}{
		{
			name:     "Secret is protected",
			path:     "/secret",
			code:     302,
			redirect: "/login",
			method:   "GET",
		},
		{
			name:     "Login Form is displayed correctly",
			path:     "/login",
			code:     200,
			method:   "GET",
			expected: "Login Form",
		},
		{
			name:     "Login Form when submitted with invalid username, returns the login form",
			path:     "/login",
			code:     200,
			method:   "POST",
			data:     strings.NewReader("username=wronguser&password=testpass"),
			expected: "Login Form",
		},
		{
			name:     "Login Form when submitted with invalid password, returns the login form",
			path:     "/login",
			code:     200,
			method:   "POST",
			data:     strings.NewReader("username=testuser&password=wrongpass"),
			expected: "Login Form",
		},
		{
			name:     "Login Form when submitted with completely invalid data, returns the login form",
			path:     "/login",
			code:     200,
			method:   "POST",
			data:     strings.NewReader("username=wronguser&password=wrongpass"),
			expected: "Login Form",
		},
		{
			name:       "Login Form when submitted with a correct password, returns a correct result",
			path:       "/login",
			code:       200,
			method:     "POST",
			data:       strings.NewReader("username=testuser&password=testpass"),
			expected:   "Hello Test",
			savecookie: true,
		},
		{
			name:      "Secret permits secure cookie",
			path:      "/secret",
			code:      200,
			expected:  "You found it!",
			usecookie: true,
		},
	}

	var cookie string
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mux := setupHttpHandlers()
			writer := httptest.NewRecorder()

			request, _ := http.NewRequest(tc.method, tc.path, nil)
			if tc.method == "POST" {
				request, _ = http.NewRequest(tc.method, tc.path, tc.data)
				request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
			}

			if tc.usecookie == true {
				t.Logf("Using cookie %v", cookie)
				request.Header.Set("Cookie", cookie)
			}
			mux.ServeHTTP(writer, request)

			testSecurityHeaders(t, writer)

			if writer.Code != tc.code {
				t.Errorf("Response code %v expected %v", writer.Code, tc.code)
			}

			if tc.redirect != "" {
				if writer.HeaderMap.Get("Location") != tc.redirect {
					t.Errorf("Expected redirect location %v, got %v", tc.redirect, writer.HeaderMap.Get("Location"))
				}
			}

			if response, err := ioutil.ReadAll(writer.Body); err != nil {
				t.Fail()
			} else {
				if !strings.Contains(string(response), tc.expected) {
					t.Errorf("Response does not contain `%v`. Got %v", tc.expected, string(response))
				}
			}

			if tc.savecookie == true {
				cookie = writer.Header().Get("Set-Cookie")
			}
		})
	}
}
```

That's a lot. Let's review the code necessary to make this work. We're going to use <a href="github.com/gorilla/sessions">github.com/gorilla/sessions</a> as well, so be sure to import that first.

Pay special attention to the NewCookieStore comments, which I have included at the top.

```
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
	var data PageData
	if r.Method == http.MethodPost {
		if r.FormValue("username") == "testuser" && r.FormValue("password") == "testpass" {
			data.User = "Test"
			session, _ := store.Get(r, "my-go-session")
			session.Values["Authenticated"] = "Yes"
			session.Save(r, w)
		}
	}
	tmpl.Execute(w, data)
}

func setupHttpHandlers() *mux.Router {
	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(securityHeaders(http.StripPrefix("/static", fs)))

	r.Path("/login").Handler(securityHeaders(http.HandlerFunc(handleLogin)))
	r.Path("/secret").Handler(securityHeaders(http.HandlerFunc(handleSecret)))
	r.PathPrefix("/").Handler(securityHeaders(http.HandlerFunc(handleRoot)))
	return r
}
```


We're almost done, let's review the outstanding tests

* OWASP Zap! Continues to pass
* Passwords are "stored" securely

We'll tackle the password storage first, and then wrap up with Owasp:

To begin, we will import golang.org/x/crypto/bcrypt, hash the passwords that we receive, and compare that way. We will replace one line in our test program, and add a quick password hash generator.

```
		if r.FormValue("username") == "testuser" && r.FormValue("password") == "testpass" {
```

Becomes

```
        password, _ := bcrypt.GenerateFromPassword([]byte("testpass"), bcrypt.DefaultCost)
        if r.FormValue("username") == "testuser" && bcrypt.CompareHashAndPassword(password, []byte(r.FormValue("password"))) == nil {
```

Obviously this should not be used exactly like this in production - you will want to store your passwords somewhere securely (perhaps a database), but this shows how simply we can hash our passwords. It is also important to follow the logging best practices - whilst the error returned should be "login failed", you should separately log that the password itself was incorrect (bcrypt.CompareHashAndPassword will return not nil in the event of an error). Whilst now would be a good time to include logging, this is already getting to be a long sprint for learning, so we will endeavour to re-visit logging at a later time.

Last test outstanding:

* OWASP Zap! Continues to pass

```
WARN-NEW: Absence of Anti-CSRF Tokens [10202] x 2
        http://172.17.0.1/login (200 OK)
        http://172.17.0.1/login (200 OK)
FAIL-NEW: 0     FAIL-INPROG: 0  WARN-NEW: 1     WARN-INPROG: 0  INFO: 0 IGNORE: 0       PASS: 49
```

Ah. CSRF Tokens.

We mitigate this through the use of github.com/gorilla/csrf. I'm going to drop out to a second page for this - check out logincsrf.md for more information.
