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

func strHex(str string) string {
  return fmt.Sprintf("%02x", str)
}

func (mbapi mbApi) splitMessage(body string) []message {
  tHelper := TextHelperInit(body)
  log.Printf("%#v", tHelper)

  if tHelper.NumParts > 1 {
    var messages []message
    refNo := rand.Intn(256)

    for i := 0; i < tHelper.NumParts; i++ {
      params := &messagebird.MessageParams{
        DataCoding: "auto", Type: "binary" }
      params.TypeDetails = make(messagebird.TypeDetails)
      params.TypeDetails["udh"] = buildUDH(
        refNo, tHelper.NumParts, i+1)

      messages = append(messages, message{
        Body: strHex(tHelper.Parts[i]), Params: params,
        //Body: body[start:end]
      })
    }

    return messages
  }

  params := &messagebird.MessageParams{ DataCoding: "auto" }
  return []message{ message{ Body: body, Params: params } }
}

func (mbapi *mbApi) SendMessage(p Payload) (
[]*messagebird.Message, error) {
  messages := mbapi.splitMessage(p.Message)
  var result []*messagebird.Message

  for _, msg := range messages {
    mbapi.checkThroughput()
    row, err := mbapi.Client.NewMessage(
      p.Originator, []string{p.Recipient},
      msg.Body, msg.Params)

    mbapi.LastRequest = time.Now()
    if err != nil {
      return result, err
    }
    result = append(result, row)
  }

  return result, nil
}
