package botify

type Message struct{}

func (m *Message) UpdateType() UpdateType {
	return UpdateTypeMessage
}
