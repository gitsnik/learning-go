## Moving to gorilla/mux

Because we have the vast majority of our testing platform in place, it should be very easy and quite simple to safely move to a new routing framework. We will do so without writing any additional tests.

We install the gorilla/mux package

```
go get -u github.com/gorilla/mux
```

And make some modifications to our main.go file (without changing our tests):

```
import (
  "net/http"
  "github.com/gorilla/mux"
)
...
func setupHttpHandlers() *mux.Router {
  r := mux.NewRouter()

  fs := http.FileServer(http.Dir("static/"))
  r.PathPrefix("/static/").Handler(securityHeaders(http.StripPrefix("/static", fs)))

  r.PathPrefix("/").Handler(securityHeaders(http.HandlerFunc(handleRoot)))
  return r
}
```

Note particularly that we've only had to make changes to the setupHttpHandlers function to use a different name (because mux is now a package - we've replaced with r) and update the static file delivery method. By changing the order of the /static/ and / paths we can be certain that the appropriate handler is used by the router.

```
$ go test
PASS
ok      _/vagrant/gorilla       0.006s
```

### Technically...

Technically, we should be testing that /static/test.txt and / return different values to be certain that our ordering is correct. We'll update the test to search for IamGroot! instead:

```
echo 'IamGroot!' > static/test.txt
# Modify main_test.go
if !strings.Contains(string(response), "IamGroot!") {
  t.Errorf("Response does not contain `IamGroot!`. Got %v", string(response))
}
```
