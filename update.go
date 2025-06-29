package botify

type UpdateType int

const (
	UpdateTypeAll UpdateType = iota
	UpdateTypeMessage
	UpdateTypeEditedMessage
	UpdateTypeChannelPost
	UpdateTypeEditedChannelPost
	UpdateTypeBusinessConnection
	UpdateTypeBusinessMessage
	UpdateTypeEditedBusinessMessage
	UpdateTypeDeletedBusinessMessages
	UpdateTypeMessageReaction
	UpdateTypeMessageReactionCount
	UpdateTypeInlineQuery
	UpdateTypeChosenInlineResult
	UpdateTypeCallbackQuery
	UpdateTypeShippingQuery
	UpdateTypePreCheckoutQuery
	UpdateTypePurchasedPaidMedia
	UpdateTypePoll
	UpdateTypePollAnswer
	UpdateTypeMyChatMember
	UpdateTypeChatMember
	UpdateTypeChatJoinRequest
	UpdateTypeChatBoost
	UpdateTypeRemovedChatBoost
)

type HandlerFunc func(ctx Context)

type Context struct {
	bot *Bot

	updType UpdateType
	updObj  UpdateObject
}

func (c *Context) Send(obj APIMethod) (*APIResponse, error) {
	return c.bot.sender.Send(obj)
}

func (c *Context) SendRaw(method string, obj any) (*APIResponse, error) {
	return c.bot.sender.SendRaw(method, obj)
}

// Use it to make sure that you're working with the expected Update type
func (c *Context) UpdateType() UpdateType {
	return c.updType
}

// It is safe to call for Message() if the handler is subscribed to UpdateTypeMessage,
// the return value is always non-nil.
// Otherwise, it will panic.
func (c *Context) Message() *Message {
	if c.updType != UpdateTypeMessage {
		panic("calling Message() when it is known that Update doesn't have any messages")
	}

	m, ok := c.updObj.(*Message)
	if !ok {
		panic("well, it is a message, but somehow it is not")
	}

	return m
}

// If you're not sure about the update type, it is safe to use TryMessage() and check for a nil value.
// It will never panic.
func (c *Context) TryMessage() *Message {
	if c.updType != UpdateTypeMessage {
		return nil
	}

	m, ok := c.updObj.(*Message)
	if !ok {
		return nil
	}

	return m
}

type UpdateObject interface {
	UpdateType() UpdateType
}
