package main

import (
  "io"
  "regexp"
  "strconv"
  "unicode"
  "encoding/json"
)

type Validator struct {}

type HttpError struct {
  Code    int
  Message string
  Key     string
}

func badRequest(key string) *HttpError {
  return &HttpError{
    Code: 400, Message: "400 bad request", Key: key }
}

type payload struct {
  Recipient   string
  Originator  string
  Message     string
}

func (v *Validator) DecodeRequestBody(body io.ReadCloser) (payload, error) {
  var p payload
  err := json.NewDecoder(body).Decode(&p)
  return p, err
}

func removePlusSign(phoneNo string) string {
  if string(phoneNo[0]) == "+" {
    return phoneNo[1:len(phoneNo)]
  }
  return phoneNo
}

func (v *Validator) validPhoneNo(phoneNo string) bool {
  if _, err := strconv.Atoi(phoneNo); err == nil {
    return true
  } else {
    return false
  }
}

var isAlphanumeric = regexp.
  MustCompile(`^[a-zA-Z0-9]+$`).MatchString

var isPhoneNo = regexp.
  MustCompile(`^\+?[0-9]+$`).MatchString

func (v *Validator) CheckPayload(p payload) *HttpError {
  // Check required attributes
  if p.Recipient == "" {
    return badRequest("recipient")
  }
  if p.Originator == "" {
    return badRequest("originator")
  }
  if p.Message == "" {
    return badRequest("message")
  }

  // Check recipient
  if !isPhoneNo(p.Recipient) {
    return badRequest("recipient")
  }
  // Check originator string or phone number
  if (!isAlphanumeric(p.Originator) && !isPhoneNo(p.Originator)) ||
  (isAlphanumeric(p.Originator) && len(p.Originator) > 11 && 
  !isPhoneNo(p.Originator)) {
    return badRequest("originator")
  }

  return nil
}
