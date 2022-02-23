package chamux

// A serializable message relating to a topic
type Frame struct {
	Topic string
	Body  []byte
}

func MakeFrame(body []byte, topicName string) *Frame {
	return &Frame{
		Topic: topicName,
		Body:  body,
	}
}
