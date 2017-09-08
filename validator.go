package main

import "regexp"

type httpError struct {
  Code    int
  Message string
  Key     string
}

func badRequest(key string) *httpError {
  return &httpError{
    Code: 400, Message: "400 bad request", Key: key }
}

type Payload struct {
  Recipient   string
  Originator  string
  Message     string
}

type Validator struct {}

func (v *Validator) isAlphanum(str string) bool {
  return regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(str)
}

func (v *Validator) isPhoneNo(str string) bool {
  return regexp.MustCompile(`^\+?[0-9]+$`).MatchString(str)
}

func (v *Validator) CheckPayload(p Payload) *httpError {
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
  if !v.isPhoneNo(p.Recipient) {
    return badRequest("recipient")
  }
  // Check originator string or phone number
  if (!v.isAlphanum(p.Originator) && !v.isPhoneNo(p.Originator)) ||
  (v.isAlphanum(p.Originator) && len(p.Originator) > 11 && 
  !v.isPhoneNo(p.Originator)) {
    return badRequest("originator")
  }

  return nil
}
