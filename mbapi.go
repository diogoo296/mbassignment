package main

import (
  "fmt"
  "log"
  "time"
  "math/rand"
  "github.com/messagebird/go-rest-api"
)

const THROUGHPUT = 1  // in seconds

type message struct {
  Body   string
  Params *messagebird.MessageParams
}

type mbApi struct {
  Throttle <-chan time.Time
  Client   *messagebird.Client
  RefNo    int
}

func initMbApi() *mbApi {
  var mbapi *mbApi
  config := LoadConfig()
  if config != nil {
    rand.Seed(time.Now().UTC().UnixNano())
    mbapi = &mbApi{
      Throttle: time.Tick(time.Second * THROUGHPUT),
      Client: messagebird.New(config.MbApiKey[GetEnv()]),
      RefNo: rand.Intn(256),
    }
  }
  return mbapi
}

var MbApiInstance *mbApi = initMbApi()

func buildUDH(refNo, total, idx int) string {
  return fmt.Sprintf("050003%02x%02x%02x", refNo, total, idx)
}

func (mbapi *mbApi) getRefNo() int {
  refNo := mbapi.RefNo
  mbapi.RefNo += 1
  if mbapi.RefNo > 255 {
    mbapi.RefNo = 0
  }
  return refNo
}

func (mbapi *mbApi) splitMessage(body string) (
[]message, error) {
  tHelper, err := InitTextHelper(body)
  if err != nil {
    return nil, err
  }
  log.Printf("%#v", tHelper)

  // Returns a single message
  if tHelper.NumParts == 1 {
    params := &messagebird.MessageParams{ DataCoding: "auto" }
    return []message{
      message{ Body: body, Params: params } }, nil
  }

  var messages []message
  refNo := mbapi.getRefNo()

  // Returns N binary messages with the proper UDH
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
    <-mbapi.Throttle  // throughput
    log.Printf("Request: %s", time.Now().String())
    row, err := mbapi.Client.NewMessage(
      p.Originator, []string{p.Recipient},
      msg.Body, msg.Params,
    )
    result = append(result, row)

    if err != nil {
      return result, err
    }
  }

  return result, nil
}
