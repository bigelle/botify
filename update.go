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

type Update struct {
	UpdateID          int      `json:"update_id"`
	Message           *Message `json:"message,omitempty"`
	EditedMessage     *Message `json:"edited_message,omitempty"`
	ChannelPost       *Message `json:"channel_post,omitempty"`
	EditedChannelPost *Message `json:"edited_channel_post,omitempty"`
	// BusinessConnection      *BusinessConnection          `json:"business_connection,omitempty"`
	BusinessMessage       *Message `json:"business_message,omitempty"`
	EditedBusinessMessage *Message `json:"edited_business_message,omitempty"`
	// DeletedBusinessMessages *BusinessMessagesDeleted     `json:"deleted_business_messages,omitempty"`
	// MessageReaction         *MessageReactionUpdated      `json:"message_reaction,omitempty"`
	// MessageReactionCount    *MessageReactionCountUpdated `json:"message_reaction_count,omitempty"`
	// InlineQuery             *InlineQuery                 `json:"inline_query,omitempty"`
	// ChosenInlineResult      *ChosenInlineResult          `json:"chosen_inline_result,omitempty"`
	// CallbackQuery           *CallbackQuery               `json:"callback_query,omitempty"`
	// ShippingQuery           *ShippingQuery               `json:"shipping_query,omitempty"`
	// PreCheckoutQuery        *PreCheckoutQuery            `json:"pre_checkout_query,omitempty"`
	// PurchasedPaidMedia      *PaidMediaPurchased          `json:"purchased_paid_media,omitempty"`
	// Poll                    *Poll                        `json:"poll,omitempty"`
	// PollAnswer              *PollAnswer                  `json:"poll_answer,omitempty"`
	// MyChatMember            *ChatMemberUpdated           `json:"my_chat_member,omitempty"`
	// ChatMember              *ChatMemberUpdated           `json:"chat_member,omitempty"`
	// ChatJoinRequest         *ChatJoinRequest             `json:"chat_join_request,omitempty"`
	// ChatBoost               *ChatBoostUpdated            `json:"chat_boost,omitempty"`
	// RemovedChatBoost        *ChatBoostRemoved            `json:"removed_chat_boost,omitempty"`
}

func (u *Update) UpdateType() UpdateType {
	if u.Message != nil {
		return UpdateTypeMessage
	}
	if u.EditedMessage != nil {
		return UpdateTypeEditedMessage
	}
	if u.ChannelPost != nil {
		return UpdateTypeChannelPost
	}
	if u.EditedChannelPost != nil {
		return UpdateTypeEditedChannelPost
	}
	if u.BusinessMessage != nil {
		return UpdateTypeBusinessMessage
	}
	if u.EditedBusinessMessage != nil {
		return UpdateTypeEditedBusinessMessage
	}
	return UpdateType(-1)
}

type HandlerFunc func(ctx Context)

type Context struct {
	bot *Bot

	updType UpdateType
	upd     *Update
}

func (c *Context) Send(obj APIMethod) (*APIResponse, error) {
	// NOTE: maybe i should store some info about the obj for logging
	return c.bot.sender.Send(obj)
}

func (c *Context) SendRaw(method string, obj any) (*APIResponse, error) {
	// NOTE: same here also
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

	if c.upd.Message == nil {
		panic("well, it is a message, but somehow it is not")
	}

	return c.upd.Message
}

// If you're not sure about the update type, it is safe to use TryMessage() and check for a nil value.
// It will never panic.
func (c *Context) TryMessage() *Message {
	return c.upd.Message
}
