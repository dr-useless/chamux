package chamux

// A serializable message relating to a topic
type Frame struct {
	Topic string
	Body  []byte
}

func NewFrame(body []byte, topicName string) *Frame {
	return &Frame{
		Topic: topicName,
		Body:  body,
	}
}
