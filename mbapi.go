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
    instance = &mbApi{
      LastRequest: time.Now().AddDate(0, 0, -1),
      Client: messagebird.New(MB_API_KEY),
    }
  })
  return instance
}

func (mbapi mbApi) checkThroughput() {
    diff := time.Since(mbapi.LastRequest)
    if diff.Seconds() < 1.0 {
      log.Printf("Sleep:", diff.Seconds())
      time.Sleep(diff)
    }
}

func (mbapi *mbApi) getBalance() *messagebird.Balance {
  mbapi.checkThroughput()
  // Request the balance information, returned as a Balance object.
  balance, err := mbapi.Client.Balance()
  mbapi.LastRequest = time.Now()

  if err != nil {
    // messagebird.ErrResponse means custom JSON errors.
    if err == messagebird.ErrResponse {
      for _, mbError := range balance.Errors {
        log.Printf("Error: %#v\n", mbError)
      }
    }

    return balance
  }

  return balance
}
