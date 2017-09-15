package main

import (
  "log"
  "strings"
  "testing"
  "io/ioutil"
)

var mbapi *mbApi = GetMbApiInstance()

func expectNMessages(p Payload, n int, t *testing.T) {
  result, err := mbapi.SendMessage(p)
  if err != nil {
    t.Errorf("Expected error not to be nil: %#v", err)
  }
  if len(result) != n {
    t.Errorf("Expected %d message(s), got %d", n, len(result))
  }
}

func TestTotalSentMessages(t *testing.T) {
  log.SetOutput(ioutil.Discard)
  p := Payload{
    Originator: "Diogo",
    Recipient: "1234",
    Message: "Test",
  }

  // Test non existent phone number in a valid format
  if result, err := mbapi.SendMessage(p); err == nil {
    t.Errorf(
      "Expected error not to be nil. Result: %#v", result)
  }

  // Test valid plain message
  p.Recipient = "5531988174420"
  expectNMessages(p, 1, t)

  // Test valid unicode concatenated message
  p.Message = strings.Repeat("日", 71)
  expectNMessages(p, 2, t)
}

func TestApiThroughput(t *testing.T) {
  log.SetOutput(ioutil.Discard)
  p := Payload{
    Originator: "Diogo",
    Recipient: "5531988174420",
    Message: strings.Repeat("a", 1400), // 10 messages
  }

  result, _ := mbapi.SendMessage(p);
  for i := 1; i < len(result); i++ {
    if diff := result[i].CreatedDatetime.Sub(
      *result[i-1].CreatedDatetime).Seconds();
    diff < 1.0 {
      t.Errorf("Expected diff >= 1; Got: %d", diff)
    }
  }
}
