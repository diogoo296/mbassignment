package main

import (
  "os"
  "testing"
  "log"
  "io/ioutil"
)

func TestLoadConfig(t *testing.T) {
  log.SetOutput(ioutil.Discard)
  config := LoadConfig()
  if config == nil {
    t.Error("Expected config to not be nil")
  }
  if config.MbApiKey == nil {
    t.Error("Expected config.MbApiKey to not be nil")
  }
  
  if key, ok := config.MbApiKey[GetEnv()]; !ok || key == "" {
    t.Error(
      "Expected config.MbApiKey to be set for current ENV")
  }
}

func TestGetEnv(t *testing.T) {
  log.SetOutput(ioutil.Discard)
  env := os.Getenv("ENV")
  funcEnv := GetEnv()
  if env == "production" && funcEnv != "production" {
    t.Error("Expected: production, got: ", funcEnv)
  } else if env != "production" && funcEnv != "development" {
    t.Error("Expected: development, got: ", funcEnv)
  }
}
