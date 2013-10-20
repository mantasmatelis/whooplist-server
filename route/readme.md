# Route #

A very efficient and minimal HTTP Go router, using a Trie data structure for efficiency.

Based on https://github.com/ant0ine/go-json-rest with the kitchen-sink philosophy stuff removed (so much that I wouldn't particularly consider it a fork).

It aims to be the innermost handler in a composition of handlers. It doesn't require methods to have http.HandlerFunc signature, route paths are interface{}. (See example.go for usage details.)

## Install ##

```go
go get github.com/mantasmatelis/go-route
```

## Basic Usage ##


```go
import(
    "net/http"
    "github.com/mantasmatelis/go-route")

var router route.Router

func main() {
    router.SetRoutes(
        route.Route{"GET",  "/count"            GetCount},
        route.Route{"POST", "/count",           IncrementCount},
        route.Route{"POST", "/count/:Count",    SetCount},
        route.Route{"POST", "/reset",           ResetCount},
    )
    
    http.ListenAndServe(":3000", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
    route, params, pathMatched := router.FindRouteFromURL(r.Method, r.URL)
    if route != nil {
        route.Func.(func(http.ResponseWriter, *http.Request))(w, r)
    }
}

```

## Ideal Usage ##

## Differences ##

How is this different from ant0ine/go-json-rest? This doesn't:
  * log
  * handle errors
  * route automatically
  * have any settings
  * override http.Request, http.ResponseWriter (making interfacing with othe rlibraries a little difficult)
  * force specific handler signatures on you (in the example you can see the ```route.Func``` is typecast before being called)

It exposes the basic router so that your own stack of handlers can be easily written, customized to your needs.
