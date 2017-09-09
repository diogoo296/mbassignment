package main

import (
  "log"
  "sync"
  "time"
  "github.com/messagebird/go-rest-api"
)

type mbApi struct {
  LastRequest time.Time
  Client *messagebird.Client
}

var instance *mbApi
var once sync.Once

func GetMbApiInstance() *mbApi {
  once.Do(func() {
    config := LoadConfig()
    if config != nil {
      instance = &mbApi{
        LastRequest: time.Now().AddDate(0, 0, -1),
        Client: messagebird.New(config.MbApiKey[GetEnv()]),
      }
    }
  })
  return instance
}

func (mbapi mbApi) checkThroughput() {
    diff := time.Since(mbapi.LastRequest)
    if diff.Seconds() < 1.0 {
      log.Printf("Sleep (ms):", diff.Seconds() * 1000)
      time.Sleep(diff)
    }
}

func (mbapi *mbApi) SendMessage(p Payload) (
*messagebird.Message, error) {
  mbapi.checkThroughput()
  msg, err := mbapi.Client.NewMessage(
    p.Originator, []string{p.Recipient}, p.Message, nil)
  mbapi.LastRequest = time.Now()
  return msg, err
}
