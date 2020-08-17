package Gateway

type ContentMessageHandler interface {
	HandleMessage(message *Message) bool
}