# URL Path Parameters

The `requests` library also supports working with [path parameters](https://swagger.io/docs/specification/describing-parameters/#path-parameters).

## What are Path parameters?

REST APIs and web servers often use variable path parameters. For example, you might have an API with the following schema:

```
/api/users/{userId}/reports/{reportId}
```

So in this example, you could access report `r456`, which is owned by user `u123`, at the following path:

```
/api/users/u123/reports/r456
```

## How can `requests` help?

`requests` can help you handle path parameters using [Go templates](https://pkg.go.dev/text/template).

When creating a Request you can include Go template syntax in the URL parameter and call the `ParsePathParams` method to fill in the data.

To continue the above example, say we're accessing multiple user's reports. We might have data in the following form:

```go
type UserReport struct {
    UserID   string
    ReportID string
}
```

We can create our Request object like so:

```go
baseReq := requests.Request{
    URL: `/api/users/{{ .UserID }}/reports/{{ .ReportID }}`
}
fmt.Println(baseReq.URL) 
// Output: /api/users/{{ .UserID }}/reports/{{ .ReportID }}
```

And we can create new request objects from our `baseReq` using instances of `UserReport`.

```go
data := []UserReport{
    {"u123", "r456"},
    {"u321", "r654"},
    {"u111", "r222"},
}

for _, ur := range data {
    req, err := baseReq.ParsePathParams(&ur)
    if err != nil {
        panic(err)
    }
    fmt.Println(req.URL)
}
```

You'll get the following output:

```
/api/users/u123/reports/r456
/api/users/u321/reports/r654
```

Go templates are a really powerful feature of the language and you can read more about them [here](https://pkg.go.dev/text/template).

## Safely Encoding Path Parameters

There is a limited set of characters that are considered _safe_ for URLs (See section 2.3 of [RFC 3986](https://www.ietf.org/rfc/rfc3986.txt)), and they are:

> ALPHA  DIGIT  "-" / "." / "_" / "~"

`requests` includes a helper function, `URLEncode`, to allow you to escape _unsafe_ characters when parsing Path parameters.

```go
func URLEncode(data interface{}) string
```

`URLEncode` is a convenience function that applies `fmt.Sprint` and then `url.PathEscape`.

In addition to being a package function, `URLEscape` is also available in the URL Go template.

For example, say we have the following URL template:

```
/api/{{ .Message }}
```

and we want to pass in the message "Hello, World!". We'll run into a bit of an issue, since the comma, the space, and the exclamation point are all considered unsafe characters.

We can update our template as follows:

```
/api/{{ .Message | URLEscape }}
```

and by _piping_ `.Message` into the `URLEscape` function, we'll end up with the following:

```
/api/Hello%2C%20World%21
```

Yay! Back to safety!

If you know your parameters are safe, you can skip `URLEscape`, but it's there if you need it.

## Handling Errors

Keep in mind that the `ParsePathParams` method returns both a new request and an error. If you're confident this call won't return an error, you can use the `MustParsePathParams` method, which only returns the new request and panics if an error is encountered.



