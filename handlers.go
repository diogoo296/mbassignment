package main

import (
  "log"
  "time"
  "net/http"
  "encoding/json"
)

func SendMessage(w http.ResponseWriter, r *http.Request) {
  // Decode payload
  var mbapi = GetMbApiInstance()
  var validator *Validator
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

  // Send message
  msgs, err := mbapi.SendMessage(payload)
  if err != nil {
    http.Error(w, "500 internal server error", 500)
    log.Printf("%#v\n", err)
    for _, msg := range msgs {
      log.Printf("%#v\n", msg.Errors)
    }
    return
  }

  // Write reply
  var times []*time.Time
  for _, msg := range msgs {
    times = append(times, msg.CreatedDatetime)
  }
  if err := json.NewEncoder(w).Encode(times); err != nil {
    http.Error(w, "500 internal server error", 500)
    log.Printf("%#v\n", err)
  }
}
