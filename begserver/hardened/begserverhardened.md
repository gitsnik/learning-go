### Running our test

Using the Zap! Baseline scan as described in the tests folder (side note: ''Securing DevOps: Security in the Cloud'' can be considered required reading) we can determine the current status of our go server in a basic format

```
$ go build
$ go test
PASS
ok      _/vagrant       0.041s
$ ./vagrant &
[1] 19411
$ ./zapbaseline.sh
Total of 3 URLs
PASS: Cookie No HttpOnly Flag [10010]
PASS: Cookie Without Secure Flag [10011]
PASS: Incomplete or No Cache-control and Pragma HTTP Header Set [10015]
PASS: Web Browser XSS Protection Not Enabled [10016]
PASS: Cross-Domain JavaScript Source File Inclusion [10017]
PASS: Content-Type Header Missing [10019]
PASS: X-Frame-Options Header Scanner [10020]
PASS: Information Disclosure - Debug Error Messages [10023]
PASS: Information Disclosure - Sensitive Information in URL [10024]
PASS: Information Disclosure - Sensitive Information in HTTP Referrer Header [10025]
PASS: HTTP Parameter Override [10026]
PASS: Information Disclosure - Suspicious Comments [10027]
PASS: Viewstate Scanner [10032]
PASS: Server Leaks Information via "X-Powered-By" HTTP Response Header Field(s) [10037]
PASS: Secure Pages Include Mixed Content [10040]
PASS: Cookie Without SameSite Attribute [10054]
PASS: CSP Scanner [10055]
PASS: X-Debug-Token Information Leak [10056]
PASS: Username Hash Found [10057]
PASS: X-AspNet-Version Response Header Scanner [10061]
PASS: Timestamp Disclosure [10096]
PASS: Cross-Domain Misconfiguration [10098]
PASS: Weak Authentication Method [10105]
PASS: Absence of Anti-CSRF Tokens [10202]
PASS: Private IP Disclosure [2]
PASS: Session ID in URL Rewrite [3]
PASS: Script Passive Scan Rules [50001]
PASS: Insecure JSF ViewState [90001]
PASS: Charset Mismatch [90011]
PASS: Application Error Disclosure [90022]
PASS: Loosely Scoped Cookie [90033]
FAIL-NEW: X-Content-Type-Options Header Missing [10021] x 4
        http://172.17.0.1/ (200 OK)
        http://172.17.0.1/robots.txt (200 OK)
        http://172.17.0.1/sitemap.xml (200 OK)
        http://172.17.0.1 (200 OK)
FAIL-NEW: 1     FAIL-INPROG: 0  WARN-NEW: 0     WARN-INPROG: 0  INFO: 0 IGNORE: 0       PASS: 31
$ killall vagrant
```

We can see here that we have a single failure - X-Content-Type-Options has not been set. The reason we see this for 4 URL's is the way the handler initially functions. Fixing this once will work for everything for now.

A quick bit of research leads us to [Mozilla Developer Pages](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options), and we know what to do.

There are two options here - we can simply continue to run the Zap test every time we test our local code, or we can add a new test to our go tests file. When using a full deployment pipeline, one will do both. For now, we add a test to the main_test.go

```
if contops := writer.Header().Get("X-Content-Type-Options"); contops != "nosniff" {
  t.Errorf("Require X-Content-Type-Options: nosniff")
}
```

And the requisite change to our main.go:

```
w.Header().Set("X-Content-Type-Options", "nosniff")
```

A new build of our service shows us tests passing in go, and re-running the zapbaseline shows the same.

Success. We have hardened a hello world service. The code for this hardened service is in this folder.
