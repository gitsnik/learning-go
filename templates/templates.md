## HTML Templates!

We can use the html/template package to provide a simple method of passing variables from Go, to the browser. We will basically copy directly from the [HTML Templates examples](https://gowebexamples.com/templates/) here, with some minor customisations to keep our testing up to date.

Because our tests have previously used contains(string()) we can use our existing environment to import html/template, and perform a basic setup.

Create html/layout.html and populate it thus:

```
Hello, you've requested: {{.RequestPath}}
```

Add the package `html/template` to the import pattern, and update the handleRoot function with these changes:

```
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
```

The specific changes here are the `PageData` struct, and replacing the w.Write line with:

```
tmpl := template.Must(template.ParseFiles("html/layout.html"))
data := PageData{
  RequestPath: r.URL.Path,
}
tmpl.Execute(w, data)
```

Because we are searching only for the contains(string()) in our test, we can modify the layout.html file as much as we like without impacting the final product. For example:

```html
<html>
  <head></head>
  <body>
    <p>
Hello, you've requested: {{.RequestPath}}
    </p>
  </body>
</html>
```
