package client

type MessageType int

const (
	NewMessage = iota
	Connect
	Disconnect
)

func (s MessageType) String() string {
	switch s {
	case NewMessage:
		return "NewMessage"
	case Connect:
		return "Connect"
	case Disconnect:
		return "Disconnect"
	}
	return "Unknown"
}

type Message struct {
	Type    MessageType
	Name    string
	Content string
}
