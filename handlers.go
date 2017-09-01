package main

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
)

func Index(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Index")
}

func Test(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "Test")
  decoder := json.NewDecoder(r.Body)
  type Payload struct {
    Username, Password string
  }
  var t Payload
  err := decoder.Decode(&t)
  if err != nil {
    fmt.Fprintln(w, err)
  }
  log.Printf("%#v\n", t)
}
