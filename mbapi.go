package main

import (
  "fmt"
  "log"
  "sync"
  "time"
  "math/rand"
  "github.com/messagebird/go-rest-api"
)

type message struct {
  Body string
  Params *messagebird.MessageParams
}

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
    log.Printf("Sleep (ms): %f", diff.Seconds() * 1000)
    time.Sleep(diff)
  }
}

func buildUDH(refNo, total, idx int) string {
  return fmt.Sprintf("050003%02x%02x%02x", refNo, total, idx)
}

func (mbapi mbApi) splitMessage(body string) (
[]message, error) {
  tHelper, err := TextHelperInit(body)
  if err != nil {
    return nil, err
  }
  log.Printf("%#v", tHelper)

  if tHelper.NumParts == 1 {
    params := &messagebird.MessageParams{ DataCoding: "auto" }
    return []message{
      message{ Body: body, Params: params } }, nil
  }

  var messages []message
  refNo := rand.Intn(256)

  for i := 0; i < tHelper.NumParts; i++ {
    typeDetails := make(messagebird.TypeDetails)
    typeDetails["udh"] = buildUDH(
      refNo, tHelper.NumParts, i+1,
    )

    params := &messagebird.MessageParams{
      Type: "binary",
      DataCoding: "auto",
      TypeDetails: typeDetails,
    }

    messages = append(messages, message{
      Body: tHelper.Parts[i], Params: params,
    })
  }

  return messages, nil
}

func (mbapi *mbApi) SendMessage(p Payload) (
[]*messagebird.Message, error) {
  var result []*messagebird.Message
  messages, err := mbapi.splitMessage(p.Message)

  if err != nil {
    return nil, err
  }

  for _, msg := range messages {
    mbapi.checkThroughput()
    row, err := mbapi.Client.NewMessage(
      p.Originator, []string{p.Recipient},
      msg.Body, msg.Params,
    )
    mbapi.LastRequest = time.Now()

    if err != nil {
      return result, err
    }
    result = append(result, row)
  }

  return result, nil
}
