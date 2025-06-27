package botify

type Update struct{}

type Message struct{}

func (m *Message) UpdateType() UpdateType {
	return UpdateTypeMessage
}
