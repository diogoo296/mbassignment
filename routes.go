package main

import "net/http"

type Route struct {
  Method  string
  Pattern string
  Handler http.HandlerFunc
}

var Routes = []Route {
  Route {
    "GET",
    "/balance",
    Balance,
  },
  Route {
    "POST",
    "/test",
    Test,
  },
}
