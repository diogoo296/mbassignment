package main

import (
  "io"
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
    Code: 400, Message: "Bad request", Key: key }
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

func checkRequiredStr(attribute string, key string) *HttpError {
  if attribute == "" {
    return badRequest(key)
  }
  return nil
}

func (v *Validator) CheckPayload(p payload) *HttpError {
  // Check required attributes
  if err := checkRequiredStr(p.Recipient, "recipient");
  err != nil {
    return err
  }
  if err := checkRequiredStr(p.Originator, "originator");
  err != nil {
    return err
  }
  if err := checkRequiredStr(p.Message, "message"); err != nil {
    return err
  }

  // Check phone numbers

  return nil
}
