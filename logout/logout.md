## Logging Out

This one is easy - we just need to destroy the sessions when visiting the /logout URL. We will add the /logout link to our main page only when logged in, and destroy the session.

Tests first of course, but we will be required to run these in order:
```
 		{
 			name:       "Secret permits secure cookie",
 			path:       "/secret",
 			code:       200,
 			method:     "GET",
 			expected:   "You found it!",
 			savecookie: true,
 			usecookie:  true,
 		},
 		{
 			name:       "Logging out redirects to /",
 			path:       "/logout",
 			code:       302,
 			method:     "GET",
 			expected:   "/",
 			usecookie:  true,
 		},
 		{
 			name:       "Secret is again protected after a log out",
 			path:       "/secret",
 			code:       302,
 			method:     "GET",
 			expected:   "/login",
 			usecookie:  true,
 		},
```

And handle the /logout URI in those tests:

```
	mux.Handle("/logout", securityHeaders(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csrfToken = csrf.Token(r)
		handleLogout(w, r)
	})))
```

For our main function we will need to expire the session and trigger the logout:

```
func handleLogout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "my-go-session")
	session.Values["Authenticated"] = "No"
	session.Options = &sessions.Options{
		MaxAge: -1,
	}
	session.Save(r, w)
	http.Redirect(w, r, "/", 302)
}
```

And of course we will want to handle the /logout path as well:

```
	r.Path("/logout").Handler(securityHeaders(csrfMiddleware(http.HandlerFunc(handleLogout))))
```
