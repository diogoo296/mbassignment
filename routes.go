package main

import (
  "log"
  "html"
  "net/http"
)

type Route struct {
  Method  string
  Pattern string
  Handler http.HandlerFunc
}

type Routes []Route

var routes = Routes {
  Route {
    "GET",
    "/",
    Index,
  },
  Route {
    "POST",
    "/test",
    Test,
  },
}

func Router() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    found := false
    url := html.EscapeString(r.URL.Path)
    log.Printf("%q %q", r.Method, url)

    for _, route := range routes {
      if url == route.Pattern && r.Method == route.Method {
        found = true
        route.Handler(w, r)
      }
    }

    if !found {
      http.NotFound(w, r)
    }
  }) 
}
