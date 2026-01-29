package event

// Message 接收到的消息
type Message struct {
	Signature string
	Timestamp string
	Nonce     string
	Body      string
}

type Event interface {
	Handle() error
}
