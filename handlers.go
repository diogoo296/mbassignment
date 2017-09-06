package main

import (
  "fmt"
  "log"
  "net/http"
)

var mbapi = GetMbApiInstance()
var validator *Validator

func Balance(w http.ResponseWriter, r *http.Request) {
  balance := mbapi.getBalance()
  fmt.Fprintf(w, "Payment: %v\n", balance.Payment)
  fmt.Fprintf(w, "Type   : %v\n", balance.Type)
  fmt.Fprintf(w, "Amount : %v"  , balance.Amount)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
  msgPayload, err := validator.DecodeRequestBody(r.Body)
  if err != nil {
    http.Error(w, err.Error(), 400)
    log.Printf("%#v\n", err)
    return
  }
  log.Printf("%#v\n", msgPayload)
  payloadErr := validator.CheckPayload(msgPayload)
  if (payloadErr != nil) {
    http.Error(w, payloadErr.Message, 400)
    log.Printf("%#v\n", payloadErr)
    return
  }
}
