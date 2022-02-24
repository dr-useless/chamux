package chamux

import (
	"errors"
	"log"
	"net"
	"sync"
)

const pf = "mconn: "

// Multiplexed connection.
// Communication is divided across a range of named topics.
// It's safe to call all methods across goroutines.
type MConn struct {
	conn       net.Conn
	ser        Serializer
	bufferSize int
	topics     map[string]*Topic
	mu         *sync.Mutex
	close      chan bool
}

// Returns a new MConn wrapping the given net.Conn.
// Set bufferSize to the maximum message length you want to receive.
func NewMConn(conn net.Conn, ser Serializer, bufferSize int) MConn {
	mc := MConn{
		conn:       conn,
		ser:        ser,
		bufferSize: bufferSize,
		topics:     make(map[string]*Topic),
		mu:         new(sync.Mutex),
		close:      make(chan bool, 1),
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
func Dial(network, address string, s Serializer, bufferSize int) (MConn, error) {
	conn, err := net.Dial(network, address)
	return NewMConn(conn, s, bufferSize), err
}

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
	_, err = mc.conn.Write(data)
	return err
}

// Reads & decodes incomming frames,
// then sends the data on the topic sub channels
func (mc *MConn) read() {
	buf := make([]byte, mc.bufferSize)
loop:
	for {
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
			_, err := mc.conn.Read(buf)
			if err != nil {
				mc.Close()
				continue
			}

			frame, err := mc.ser.Deserialize(buf)
			if err != nil {
				log.Println(pf+"error deserializing frame:", err)
				continue
			}

			topic := mc.topics[frame.Topic]
			if topic == nil {
				log.Println(pf+"unknown topic:", frame.Topic)
				continue
			}

			for _, sub := range topic.subs {
				sub <- frame.Body
			}
		}
	}
}
