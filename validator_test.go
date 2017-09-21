package main

import "testing"

func expectKeyError(err *httpError, key string, t *testing.T) {
  if err == nil {
    t.Error("Expected error to not be nil")
  } else if err.Key != key {
    t.Error("Expected error.Key to not be equal ", key)
  }
}

func expectNilError(err *httpError, t *testing.T) {
  if err != nil {
    t.Error("Expected error to be nil, got: ", err)
  }
}

func TestCheckPayloadMessage(t *testing.T) {
  var validator *Validator
  p := Payload{
    Recipient: "123456789",
    Originator: "123456789",
    Message: "",
  }

  // Test empty message
  expectKeyError(validator.CheckPayload(p), "message", t)

  // Test whitespaces message
  p.Message = " \n     "
  expectKeyError(validator.CheckPayload(p), "message", t)

  // Test not empty message
  p.Message = "Testing 1234 []()"
  expectNilError(validator.CheckPayload(p), t)
}

func TestCheckPayloadRecipient(t *testing.T) {
  var validator *Validator
  p := Payload{
    Recipient: "",
    Originator: "123456789",
    Message: "Test message",
  }

  // Test empty recipient
  expectKeyError(validator.CheckPayload(p), "recipient", t);

  // Test alphanumeric recipient
  p.Recipient = "a123456789"
  expectKeyError(validator.CheckPayload(p), "recipient", t);

  // Test recipient starting with plus sign
  p.Recipient = "+123456789"
  expectNilError(validator.CheckPayload(p), t)

  // Test recipient with only numbers
  p.Recipient = "123456789"
  expectNilError(validator.CheckPayload(p), t)
}

func TestCheckPayloadOriginator(t *testing.T) {
  var validator *Validator
  p := Payload{
    Recipient: "123456789",
    Originator: "",
    Message: "Test message",
  }

  // Test empty originator
  expectKeyError(validator.CheckPayload(p), "originator", t);

  // Test non alphanumeric originator
  p.Originator = "Test)*"
  expectKeyError(validator.CheckPayload(p), "originator", t);

  // Test too long alphanumeric originator (len > 11)
  p.Originator = "TestHasLenghGreaterThan11"
  expectKeyError(validator.CheckPayload(p), "originator", t);

  // Test valid alphanumeric originator
  p.Originator = "TestIsFine"
  expectNilError(validator.CheckPayload(p), t)

  // Test recipient starting with plus sign
  p.Originator = "+123456789"
  expectNilError(validator.CheckPayload(p), t)

  // Test recipient with only numbers
  p.Originator = "123456789"
  expectNilError(validator.CheckPayload(p), t)
}
