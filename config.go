package main

import (
  "os"
  "log"
  "path"
  "runtime"
  "encoding/json"
)

type Config struct {
  MbApiKey map[string]string 
}

func LoadConfig() *Config {
  var config *Config
  _, filename, _, ok := runtime.Caller(1)
  if !ok {
    log.Println("Error loading config file")
    return nil
  }
  configFile, err := os.
    Open(path.Join(filename, "../config.json"))
  if err != nil {
    log.Println("Error loading config file")
    return nil
  }
  if err := json.NewDecoder(configFile).Decode(&config);
  err != nil {
    log.Println("Error loading config file")
    return nil
  }
  return config
}

func GetEnv() string {
  env := os.Getenv("ENV")
  if env != "production" {
    env = "development"
  }
  log.Println("ENV:", env)
  return env
}
