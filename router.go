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

var Routes = []Route {
  Route {
    "POST",
    "/messages",
    SendMessage,
  },
}

type Router struct {
  Routes []Route
}

// Link endpoints to handlers
func (router *Router) MapRoutes() {
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    found := false
    url := html.EscapeString(r.URL.Path)
    log.Printf("%q %q", r.Method, url)

    for _, route := range router.Routes {
      if url == route.Pattern && r.Method == route.Method {
        found = true
        route.Handler.ServeHTTP(w, r)
      }
    }

    if !found {
      http.NotFound(w, r)
    }
  }) 
}
