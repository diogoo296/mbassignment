package main

import (
  "log"
  "net/http"
)

func main() {
  router := &Router{Routes}
  router.MapRoutes()
  log.Println("Server started!")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
