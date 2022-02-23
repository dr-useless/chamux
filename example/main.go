package main

import (
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/dr-useless/chamux"
)

func main() {
	rand.Seed(time.Now().UnixMicro())

	// make listener to accept connections
	ready := makeListener()
	<-ready

	// connect to the listener
	mc, err := chamux.Dial("unix", "/tmp/example", chamux.Gob{}, 2048)
	if err != nil {
		panic(err)
	}

	// create a topic
	topic := chamux.NewTopic("dog")

	// get a channel for messages about dogs
	channel := topic.Subscribe()

	// register our topic
	mc.AddTopic(&topic)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

loop:
	for {
		select {
		case <-sigint:
			mc.Close()
			break loop
		case dog := <-channel:
			// do something with the messages
			log.Println("dog: " + string(dog))
		}
	}
}

func makeListener() chan bool {
	ready := make(chan bool, 1)
	go func(ready chan bool) {
		listener, err := net.Listen("unix", "/tmp/example")
		if err != nil {
			panic(err)
		}
		ready <- true

		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}

			// implement chamux.Serializer to use another encoding
			mc := chamux.NewMConn(conn, chamux.Gob{}, 2048)

			go func(chamux.MConn) {
				for {
					time.Sleep(time.Second)
					msg := []byte(getRandomName())

					// make a new frame with the message & topic name
					frame := chamux.NewFrame(msg, "dog")

					// Publish(s,f) serializes the frame to a byte slice,
					// and then writes the slice to the underlying connection.
					mc.Publish(chamux.Gob{}, frame)
				}
			}(mc)
		}
	}(ready)
	return ready
}
