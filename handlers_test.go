package main

import (
  "log"
  "bytes"
  "strings"
  "testing"
  "net/http"
  "io/ioutil"
  "encoding/json"
  "net/http/httptest"
  "github.com/messagebird/go-rest-api"
)

func testRequest(p Payload, expectedCode, expectedNumMsgs int,
t *testing.T) {
  body, _ := json.Marshal(p)

  request, err := http.NewRequest(
    "POST", "/messages", bytes.NewReader(body))

  if err != nil {
    t.Errorf("Expected error to be nil: %#v", err)
  }

  recorder := httptest.NewRecorder()
  handler  := http.HandlerFunc(SendMessage)
  handler.ServeHTTP(recorder, request)

  if recorder.Code != expectedCode {
    t.Errorf("Expected statusCode %d; Got: %d",
      expectedCode, recorder.Code)
  }
  if recorder.Code == http.StatusOK {
    var result []*messagebird.Message
    err = json.NewDecoder(recorder.Body).Decode(&result)
    if len(result) != expectedNumMsgs {
      t.Errorf("Expected %d messages; Got: %d",
        expectedNumMsgs, len(result))
    }
  }
}

// Test only API response
func TestMessageHandler(t *testing.T) {
  log.SetOutput(ioutil.Discard)
  payload := Payload{
    Originator: "Diogo",
    Recipient: "5531988174420",
    Message: "Test message",
  }
  testRequest(payload, 200, 1, t)

  payload = Payload{
    Originator: "",
    Recipient: "5531988174420",
    Message: "Test message",
  }
  testRequest(payload, 400, 0, t)

  payload = Payload{
    Originator: "Diogo",
    Recipient: "5531988174420",
    Message: strings.Repeat("a", 170),
  }
  testRequest(payload, 200, 2, t)
}
