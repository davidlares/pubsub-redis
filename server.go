package main

import (
  "fmt"
  "gopkg.in/redis.v3"
  "encoding/json"
  "github.com/gorilla/mux"
  "github.com/gorilla/websocket"
  "net/http"
)

type Request struct {
  Id int
  Name string
}

type Client struct {
  Id int
  websocket *websocket.Conn
}

// store Client structs for bidi communications
var Clients = make(map[int] Client)

func ConnectNewClient(channel_request chan Request){
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

    // getting data
    channel_request <- request

    // fmt.Println(message.Channel)
    // fmt.Println(message.Payload)
  }
  fmt.Println("Connection stablished")
}

func Subscribe(w http.ResponseWriter, r *http.Request){
  ws, err := websocket.Upgrade(w,r,nil,1024,1024)
  if err != nil {
    return
  }
  fmt.Println("New Websocket")
  // counting clients
  count := len(Clients)
  // adding clients
  new_client := Client{count, ws}
  Clients[count] = new_client
  fmt.Println("New client")
  for {
    _,_,err := ws.ReadMessage()
    if err != nil {
      delete(Clients, new_client.Id)
      fmt.Println("Client gone")
      return
    }
  }
}

func ValidateChannel(request chan Request){
  for {
    select{
    case r := <- request:
      // send messages
      SendMessage(r)
    }
  }
}

func SendMessage(request Request){
  for _, client := range Clients{
    if err := client.websocket.WriteJSON(request); err != nil {
      return
    }
  }
}

func main(){
  fmt.Println("Hello, World")
  // channel request
  channel_request := make(chan Request)
  go ValidateChannel(channel_request)
  // go Routing
  go ConnectNewClient(channel_request)
  // Websocket server
  mux := mux.NewRouter()
  mux.HandleFunc("/subscribe", Subscribe).Methods("GET")
  http.Handle("/", mux)
  fmt.Println("Server running on 8000")
  http.ListenAndServe(":8000", nil)
}
