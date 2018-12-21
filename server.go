package main

import (
  "fmt"
  "gopkg.in/redis.v3"
  "encoding/json"
)

type Request struct {
  Id int
  Name string
}

func ConnectNewClient(){
  // redis client
  client := redis.NewClient(&redis.Options{
    Addr: "localhost:6379",
    Password: "",
    DB: 0,
  })
  // subscribing to test1
  pubsub, err := client.Subscribe("test1")
  if err != nil {
    fmt.Println("Unable to connect")
  }

  // infinite loop for reach messages
  for {
    message,err := pubsub.ReceiveMessage()
    if err != nil {
      fmt.Println("Unable to read message")
    }
    
    request := Request{}
    if err := json.Unmarshal([]byte(message.Payload), &request); err != nil {
      fmt.Println("Unable to read JSON")
    }
    fmt.Println(request.Id)
    fmt.Println(request.Name)
    // fmt.Println(message.Channel)
    // fmt.Println(message.Payload)
  }

  fmt.Println("Connection stablished")
}

func main(){
  fmt.Println("Hello, World")
  ConnectNewClient()
}
