package main

import (
  "log"
  "net/http"
)

func main() {
  Router()
  log.Println("Server started!")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
