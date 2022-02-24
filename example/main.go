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
	topic := chamux.NewTopic("coffee")

	// get a channel for messages about coffee
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
		case msg := <-channel:
			// do something with the messages
			log.Println("coffee: " + string(msg))
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
					msg := []byte(getRandomCoffeeStrain())

					// make a new frame with the message & topic name
					frame := chamux.NewFrame(msg, "coffee")

					// Publish(f) serializes the frame to a byte slice,
					// and then writes the slice to the underlying connection.
					mc.Publish(frame)
				}
			}(mc)
		}
	}(ready)
	return ready
}

func getRandomCoffeeStrain() string {
	strains := []string{"Arabica", "Robusta", "Liberica", "Bourbon",
		"Catimor", "Catuai", "Caturra", "Heirloom", "Geisha", "Mundo Novo",
		"Timor", "Typica", "Variedad Colombia", "Yellow Bourbon"}
	return strains[rand.Intn(len(strains)-1)]
}
