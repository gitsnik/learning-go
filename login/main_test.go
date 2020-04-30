package main

import (
	"github.com/gorilla/csrf"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

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

func TestHandleRoot(t *testing.T) {
	mux := setupHttpHandlers()
	writer := httptest.NewRecorder()

	request, _ := http.NewRequest("GET", "/", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code %v expected 200", writer.Code)
	}

	testSecurityHeaders(t, writer)

	if response, err := ioutil.ReadAll(writer.Body); err != nil {
		t.Fail()
	} else if !strings.Contains(string(response), "Hello, you've requested: /") {
		t.Errorf("Response does not contain `Hello, you've requested: /`: %v", string(response))
	}
}

func TestSetupHttpHandlers(t *testing.T) {
	mux := setupHttpHandlers()
	writer := httptest.NewRecorder()

	request, _ := http.NewRequest("GET", "/static/test.txt", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != 200 {
		t.Errorf("Response code %v expected 200", writer.Code)
	}

	testSecurityHeaders(t, writer)

	if response, err := ioutil.ReadAll(writer.Body); err != nil {
		t.Fail()
	} else {
		if !strings.Contains(string(response), "IamGroot!") {
			t.Errorf("Response does not contain `IamGroot!`. Got %v", string(response))
		}
	}
}

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
			method:    "GET",
			expected:  "You found it!",
			usecookie: true,
		},
	}

	var cookie string
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mux := http.NewServeMux()
			csrfKey := generateCode(32)
			p := csrf.Protect([]byte(csrfKey))(mux)

			var csrfToken string

			mux.Handle("/login", securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				csrfToken = csrf.Token(r)
				handleLogin(w, r)
			})))

			mux.Handle("/secret", securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				csrfToken = csrf.Token(r)
				handleSecret(w, r)
			})))

			mux.Handle("/", securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				csrfToken = csrf.Token(r)
				handleRoot(w, r)
			})))

			writer := httptest.NewRecorder()

			request, _ := http.NewRequest("GET", tc.path, nil)
			if tc.method != "GET" {
				p.ServeHTTP(writer, request)

				request, _ = http.NewRequest(tc.method, tc.path, tc.data)
				request.Header.Set("X-CSRF-Token", csrfToken)
				request.Header.Set("Cookie", writer.Header().Get("Set-Cookie"))

				writer = httptest.NewRecorder()

				request.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
			}

			if tc.usecookie == true {
				t.Logf("Using cookie %v", cookie)
				request.Header.Set("Cookie", cookie)
			}
			p.ServeHTTP(writer, request)

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
