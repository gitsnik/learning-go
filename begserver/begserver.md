## Beginning Server

We start with the [Hello World](https://gowebexamples.com/hello-world/) example from [gowebexamples.com](https://gowebexamples.com), to which I have added a working (but not pleasant to read) test. Compile and test with:

```
go build
go test
```

This verifies simple functionality, but does not include any security. Before we continue down any path of complexity, let us add two very important tests.

### Terry Pratchett

The most important header to be added to any webserver is:

```
X-Clacks-Overhead: GNU Terry Pratchett
```

We add this test to our main_test.go, and verify that the test fails:

```
if clacks := writer.Header().Get("X-Clacks-Overhead"); clacks != "GNU Terry Pratchett" {
  t.Errorf("Require X-Clacks-Overhead: GNU Terry Pratchett. Never forget.")
}
...
go test
--- FAIL: TestHandleRoot (0.00s)
    main_test.go:31: Require X-Clacks-Overhead: GNU Terry Pratchett. Never forget.
FAIL
exit status 1
FAIL    _/vagrant      0.045s
```

Now we add this line to our main.go:

```
w.Header().Set("X-Clacks-Overhead", "GNU Terry Pratchett")
...
go test
PASS    _/vagrant      0.045s
```

Moving on to hardening our build

### Hardening the beginning server

[Hardening Beginning Server](begserverhardened.md)
