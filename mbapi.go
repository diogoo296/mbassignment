package main

import (
  "fmt"
  "log"
  "math"
  "sync"
  "time"
  "github.com/messagebird/go-rest-api"
)

const (
  SMS_MAX_LEN  = 160
  CSMS_MAX_LEN = 153
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

func (mbap mbApi) splitMessage(body string) []message {
  bodyLen := len(body)

  if bodyLen > SMS_MAX_LEN {
    var messages []message
    total := int(
      math.Ceil(float64(bodyLen) / float64(CSMS_MAX_LEN)))

    for i := 0; i < total; i++ {
      start := i * CSMS_MAX_LEN
      end   := (i + 1) * CSMS_MAX_LEN
      if end > bodyLen {
        end = bodyLen
      }

      typeDetails := make(messagebird.TypeDetails)
      typeDetails["udh"] = buildUDH(1, total, i+1)
      messages = append(
        messages, message{Body: strHex(body[start:end]),
        //messages, message{Body: body[start:end],
        Params: &messagebird.MessageParams{
          Type: "binary", TypeDetails: typeDetails }},
      )
    }
    return messages
  }
  return []message{message{Body: body, Params: nil}}
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
