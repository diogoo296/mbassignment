package main

import (
  "fmt"
  "log"
  "net/http"
  "encoding/json"
)

var mbapi = GetMbApiInstance()

func Balance(w http.ResponseWriter, r *http.Request) {
  balance := mbapi.getBalance()
  fmt.Fprintf(w, "Payment: %v\n", balance.Payment)
  fmt.Fprintf(w, "Type   : %v\n", balance.Type)
  fmt.Fprintf(w, "Amount : %v"  , balance.Amount)
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
