# chamux
A simple net.Conn multiplexer based on named topics & channels. Subscribing to a topic returns a channel on which updates will be sent.

## Usage
```go
func dialer() {
  // ignore the errors for brevity (it's an example)
  conn, _ := net.Dial("unix", "/tmp/example")

  // implement chamux.Serializer to use your own encoding
  mc := chamux.NewMConn(conn, chamux.GobSerializer{}, 2048)

  // get a channel for a topic named: coffee
  topic := chamux.NewTopic("coffee")
  subscription := topic.Subscribe()
  mc.AddTopic(&topic)

  sigint := make(chan os.Signal, 1)
  signal.Notify(sigint, os.Interrupt)

loop:
  for {
    select {
    case <-sigint:
      mc.Close()
      break loop
    case msg := <-subscription:
      log.Println("something about coffee:", string(msg))
  }
}

func listener() {
  listener, _ := net.Listen("unix", "/tmp/example")
  for {
    conn, _ := listener.Accept()
    mc := chamux.NewMConn(conn, 2048)

    // send lots of messages about coffee
    go func(chamux.MConn) {
      for {
        msg := []byte("we need more coffee")
        frame := chamux.NewFrame(msg, "coffee")
        mc.Publish([]byte(chamux.GobSerializer{}, frame)
      }
    }(mc)
  }
}
```

## Serialization
This package exports `chamux.GobSerializer{}`. If you want to use another encoding, simply implement `Serializer`.
```go
type Serializer interface {
  Serialize(f *Frame) ([]byte, error)
  Deserialize(f []byte) (*Frame, error)
}
```

