package main

import (
  "log"
  "net/http"
  "unicode/utf8"
  "encoding/json"
)

var mbapi = GetMbApiInstance()
var validator *Validator

func SendMessage(w http.ResponseWriter, r *http.Request) {
  // Decode payload
  var payload Payload
  if err := json.NewDecoder(r.Body).Decode(&payload);
  err != nil {
    http.Error(w, err.Error(), 400)
    log.Printf("%#v\n", err)
    return
  }

  // Validate payload
  if err := validator.CheckPayload(payload); err != nil {
    http.Error(w, err.Message, err.Code)
    log.Printf("%#v\n", err)
    return
  }

  log.Printf("#Message: %#v", len(payload.Message))
  log.Printf("#Runes:   %#v",
    utf8.RuneCountInString(payload.Message))

  // Send message
  msg, err := mbapi.SendMessage(payload)
  if err != nil {
    http.Error(w, "500 internal server error", 500)
    log.Printf("%#v\n", err)
    return
  }

  // Write reply
  if err := json.NewEncoder(w).Encode(msg); err != nil {
    http.Error(w, "500 internal server error", 500)
    log.Printf("%#v\n", err)
  }
}
