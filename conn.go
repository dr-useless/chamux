package chamux

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

const pf = "mconn: "

// Multiplexed connection.
// Communication is divided across a range of named topics.
// It's safe to call all methods across goroutines.
type MConn struct {
	conn   net.Conn
	ser    Serializer
	topics map[string]*Topic
	mu     *sync.Mutex
	close  chan bool
}

type Options struct {
	ReadDeadline time.Duration
}

// Returns a new MConn wrapping the given net.Conn
func NewMConn(conn net.Conn, ser Serializer, opt Options) MConn {
	mc := MConn{
		conn:   conn,
		ser:    ser,
		topics: make(map[string]*Topic),
		mu:     new(sync.Mutex),
		close:  make(chan bool, 1),
	}

	if opt.ReadDeadline != 0 {
		mc.conn.SetReadDeadline(time.Now().Add(opt.ReadDeadline))
	}

	go mc.read()

	return mc
}

// Shortcut that dials for a net.Conn & returns an MConn.
//
// Same as calling:
//
// conn, _ := net.Dial(network, address)
//
// mc, _ := NewMConn(conn, s, bufferSize)
func Dial(network, address string, s Serializer, opt Options) (MConn, error) {
	conn, err := net.Dial(network, address)
	return NewMConn(conn, s, opt), err
}

// Closes all subscription channels,
// then closes the underlying connection
func (mc *MConn) Close() error {
	mc.close <- true
	return mc.conn.Close()
}

// Adds a named topic to MConn.
// Returns an error if there is already a topic with the same name.
func (mc *MConn) AddTopic(t *Topic) error {
	mc.mu.Lock()
	if mc.topics[t.name] != nil {
		return errors.New(pf + "already has a topic with that name")
	}
	mc.topics[t.name] = t
	mc.mu.Unlock()
	return nil
}

// Serializes & writes a Frame to the underlying connection
func (mc *MConn) Publish(f *Frame) error {
	data, err := mc.ser.Serialize(f)
	if err != nil {
		return err
	}

	// append marker for split function
	data = append(data, []byte("+END")...)

	_, err = mc.conn.Write(data)
	return err
}

// Reads & decodes incomming frames,
// then sends the data on the topic sub channels
func (mc *MConn) read() {
	scan := bufio.NewScanner(mc.conn)
	scan.Split(splitFunc)
loop:
	for scan.Scan() {
		select {
		case <-mc.close:
			for _, topic := range mc.topics {
				for _, sub := range topic.subs {
					close(sub)
				}
				delete(mc.topics, topic.name)
			}
			break loop

		default:
			frameBytes := scan.Bytes()

			if frameBytes == nil {
				mc.Close()
				continue
			}

			frame, err := mc.ser.Deserialize(frameBytes)
			if err != nil {
				fmt.Println(pf + err.Error())
				continue
			}

			topic := mc.topics[frame.Topic]
			if topic == nil {
				fmt.Println(pf+"unknown topic:", frame.Topic)
				continue
			}

			for _, sub := range topic.subs {
				sub <- frame.Body
			}
		}
	}
}
