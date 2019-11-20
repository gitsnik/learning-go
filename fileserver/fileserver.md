## Serving Static Assets

Technically we already have (the code)[https://gowebexamples.com/http-server/] for this, but we're going to begin as we intend to finish - tests first.

Because we're starting to need it, we will also refactor our server code to be more robust and get most of it out of main(). We will use our tests to drive our redevelopment. First we redfine our mux calls to use a new function (setupHttpHandlers) in our existing test, and then we add a new test for static resources:

```
func TestSetupHttpHandlers(t *testing.T) {
  mux := setupHttpHandlers()
  writer := httptest.NewRecorder()

  request, _ := http.NewRequest("GET", "/static/test.txt", nil)
  mux.ServeHTTP(writer, request)

  if writer.Code != 200 {
    t.Errorf("Response code %v expected 200", writer.Code)
  }

  if response, err := ioutil.ReadAll(writer.Body); err != nil {
    t.Fail()
  } else {
    if !strings.Contains(string(response), "Hello") {
      t.Errorf("Response does not contain `Hello`. Got %v", string(response))
    }
  }
}
```

Create a file and folder:

```
mkdir static
echo dog > static/test.txt
go test
# _/vagrant/fileserver [_/vagrant/fileserver.test]
./main.go:15:10: undefined: setupHttpHandlers
./main_test.go:12:10: undefined: setupHttpHandlers
./main_test.go:40:10: undefined: setupHttpHandlers
FAIL    _/vagrant/fileserver [build failed]
```

Now we will refactor our func main() and add a setupHttpHandlers function. Note the removal of HandleFunc()

```
func setupHttpHandlers() *http.ServeMux {
  mux := http.NewServeMux()
  mux.Handle("/", handleRoot)

  fs := http.FileServer(http.Dir("static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fs))
  return mux
}

func main() {
  mux := setupHttpHandlers()
  http.ListenAndServe(":80", mux)
}
...
$ go test
--- FAIL: TestSetupHttpHandlers (0.00s)
    main_test.go:54: Response does not contain `Hello`. Got dog
FAIL
exit status 1
FAIL    _/vagrant/fileserver    0.006s
```

This shows us that the static resources are working correctly, but the test doesn't pass because the file contains the wrong information. Fix that, and the test passes:

```
$ echo Hello > static/test.txt
$ go test
PASS
ok      _/vagrant/fileserver    0.006s
```

We have noticed, however, that we have lost our Clacks. We will refactor our tests again to reduce duplication:

```
func testSecurityHeaders(t *testing.T, writer *httptest.ResponseRecorder) {
  if clacks := writer.Header().Get("X-Clacks-Overhead"); clacks != "GNU Terry Pratchett" {
    t.Errorf("Require X-Clacks-Overhead: GNU Terry Pratchett. Never forget.")
  }

  if contopts := writer.Header().Get("X-Content-Type-Options"); clacks != "nosniff" {
    t.Errorf("Require X-Content-Type-Options: nosniff")
  }
}

func TestHandleRoot(t *testing.T) {
  ...
  testSecurityHeaders(t, writer)
  ...

func TestSetupHttpHandlers(t *testing.T) {
  ...
  testSecurityHeaders(t, writer)
  ...

$ go test
--- FAIL: TestSetupHttpHandlers (0.00s)
    main_test.go:13: Require X-Clacks-Overhead: GNU Terry Pratchett. Never forget.
    main_test.go:17: Require X-Content-Type-Options: nosniff
FAIL
exit status 1
FAIL    _/vagrant/fileserver    0.007s
```

### Heavy refactoring

Because we want to introduce our security headers to static code, and because we want to reduce duplication of work, it's time to offload the header functionality to a piece of middleware.

To do this, we want to refactor the handleRoot function slightly, and move our Clacks and Content options headers to a middleware function.

```
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
...
mux.Handle("/", securityHeaders(http.HandlerFunc(handleRoot)))
...
mux.Handle("/static/", securityHeaders(http.StripPrefix("/static", fs)))
```

And now our tests pass:

```
$ go test
PASS
ok      _/vagrant/fileserver    0.006s
$
```

#### Finishing Up

This takes me to the end of (the code)[https://gowebexamples.com/http-server/]. Hardening has been built in with new middleware and we've reduced a bit of duplication. Static file delivery will allow us to automatically deploy CSS, JS, and even JSON extensions as needed.
