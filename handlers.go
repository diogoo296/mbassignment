package main

import (
  "log"
  "net/http"
  "encoding/json"
)

var mbapi = GetMbApiInstance()
var validator *Validator

func SendMessage(w http.ResponseWriter, r *http.Request) {
  payload, err := validator.DecodeRequestBody(r.Body)
  if err != nil {
    http.Error(w, err.Error(), 400)
    log.Printf("%#v\n", err)
    return
  }

  payloadErr := validator.CheckPayload(payload)
  if (payloadErr != nil) {
    http.Error(w, payloadErr.Message, 400)
    log.Printf("%#v\n", payloadErr)
    return
  }

  msg, mbErr := mbapi.SendMessage(payload)
  if err != nil {
    http.Error(w, "500 internal server error", 500)
    log.Printf("%#v\n", mbErr)
    return
  }

  err = json.NewEncoder(w).Encode(msg)
  if err != nil {
    http.Error(w, "500 internal server error", 500)
    log.Printf("%#v\n", err)
  }
}
