package chamux

import "sync"

// A named group,
// to which multiple things (goroutines) can subscribe.
type Topic struct {
	name string
	mu   *sync.Mutex
	subs []chan []byte
}

func NewTopic(name string) Topic {
	return Topic{
		name: name,
		mu:   new(sync.Mutex),
	}
}

func (t *Topic) Name() string {
	return t.name
}

// Returns a recvr channel for updates on the topic
func (t *Topic) Subscribe() <-chan []byte {
	sub := make(chan []byte)
	t.mu.Lock()
	t.subs = append(t.subs, sub)
	t.mu.Unlock()
	return sub
}
