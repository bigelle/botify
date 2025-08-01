package botify

import (
	"context"
)

const (
	UpdateTypeMessage                 = "message"
	UpdateTypeEditedMessage           = "edited_message"
	UpdateTypeChannelPost             = "channel_post"
	UpdateTypeEditedChannelPost       = "edited_channel_post"
	UpdateTypeBusinessConnection      = "business_connection"
	UpdateTypeBusinessMessage         = "business_message"
	UpdateTypeEditedBusinessMessage   = "edited_business_message"
	UpdateTypeDeletedBusinessMessages = "deleted_business_messages"
	UpdateTypeMessageReaction         = "message_reaction"
	UpdateTypeMessageReactionCount    = "message_reaction_count"
	UpdateTypeInlineQuery             = "inline_query"
	UpdateTypeChosenInlineResult      = "chosen_inline_result"
	UpdateTypeCallbackQuery           = "callback_query"
	UpdateTypeShippingQuery           = "shipping_query"
	UpdateTypePreCheckoutQuery        = "pre_checkout_query"
	UpdateTypePurchasedPaidMedia      = "purchased_paid_media"
	UpdateTypePoll                    = "poll"
	UpdateTypePollAnswer              = "poll_answer"
	UpdateTypeMyChatMember            = "my_chat_member"
	UpdateTypeChatMember              = "chat_member"
	UpdateTypeChatJoinRequest         = "chat_join_request"
	UpdateTypeChatBoost               = "chat_boost"
	UpdateTypeRemovedChatBoost        = "removed_chat_boost"
)

var allUpdTypes = map[string]struct{}{
	UpdateTypeMessage:                 {},
	UpdateTypeEditedMessage:           {},
	UpdateTypeChannelPost:             {},
	UpdateTypeEditedChannelPost:       {},
	UpdateTypeBusinessConnection:      {},
	UpdateTypeBusinessMessage:         {},
	UpdateTypeEditedBusinessMessage:   {},
	UpdateTypeDeletedBusinessMessages: {},
	UpdateTypeMessageReaction:         {},
	UpdateTypeMessageReactionCount:    {},
	UpdateTypeInlineQuery:             {},
	UpdateTypeChosenInlineResult:      {},
	UpdateTypeCallbackQuery:           {},
	UpdateTypeShippingQuery:           {},
	UpdateTypePreCheckoutQuery:        {},
	UpdateTypePurchasedPaidMedia:      {},
	UpdateTypePoll:                    {},
	UpdateTypePollAnswer:              {},
	UpdateTypeMyChatMember:            {},
	UpdateTypeChatMember:              {},
	UpdateTypeChatJoinRequest:         {},
	UpdateTypeChatBoost:               {},
	UpdateTypeRemovedChatBoost:        {},
}

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

// UpdateType returns update type as defined in [Update]
// 
// [Update]: https://core.telegram.org/bots/api#update
func (u *Update) UpdateType() string {
	checks := []struct {
		condition bool
		result    string
	}{
		{u.Message != nil, UpdateTypeMessage},
		{u.EditedMessage != nil, UpdateTypeEditedMessage},
		{u.ChannelPost != nil, UpdateTypeChannelPost},
		{u.EditedChannelPost != nil, UpdateTypeEditedChannelPost},
		{u.BusinessMessage != nil, UpdateTypeBusinessMessage},
		{u.EditedBusinessMessage != nil, UpdateTypeEditedBusinessMessage},
	}

	for _, check := range checks {
		if check.condition {
			return check.result
		}
	}

	return ""
}

// HandlerFunc is a function used to handle incoming updates.
type HandlerFunc func(ctx *Context) error

func ChainHandlers(handlers ...HandlerFunc) HandlerFunc {
	return func(ctx *Context) error {
		for _, handler := range handlers {
			if err := handler(ctx); err != nil {
				return err
			}
		}
		return nil
	}
}

// Context is a "bridge",
// which allows to use bot settings such as logger or [RequestSender]
// right from the handler, without the need to store a pointer to the bot's instance.
type Context struct {
	bot *Bot

	updType        string
	upd            *Update

	ctx context.Context
}

// Bot returns a read-only copy of the [Bot] that created this context
func (c *Context) Bot() Bot {
	return *c.bot
}

// SendRequest is a wrapper around [SendRequestContext],
// which is using the inner [context.Context] as ctx
func (c *Context) SendRequest(obj APIMethod) (*APIResponse, error) {
	return c.SendRequestContext(c.Context(), obj)
}

// SendJSON is a wrapper around [SendJSONContext],
// which is using the inner [context.Context] as ctx
func (c *Context) SendJSON(endpoint string, obj any) (*APIResponse, error) {
	return c.SendJSONContext(c.Context(), endpoint, obj)
}

// SendRequestContext is using bot's [RequestSender]
// to send a request with payload, content-type and to the API endpoint, defined in obj,
// and can be cancelled with ctx
func (c *Context) SendRequestContext(ctx context.Context, obj APIMethod) (*APIResponse, error) {
	return c.bot.Sender.SendWithContext(ctx, obj)
}

// SendJSONContext is using bot's [RequestSender]
// to send obj as JSON to the endpoint,
// and can be cancelled with ctx
func (c *Context) SendJSONContext(ctx context.Context, endpoint string, obj any) (*APIResponse, error) {
	return c.bot.Sender.SendJSONWithContext(ctx, endpoint, obj)
}

// UpdateType returns update type.
// Useful to make sure that the type is what was expected.
// See [Update] for a complete list of available update types
//
// [Update]: https://core.telegram.org/bots/api#update
func (c *Context) UpdateType() string {
	return c.updType
}

// UpdateID returns update ID.
func (c *Context) UpdateID() int {
	return c.upd.UpdateID
}

// Context returns [context.Context],
// which in most cases is [context.WithCancel],
// controlled by the currently active bot instance.
// Useful as a parent context for timeouts and deadlines.
func (c *Context) Context() context.Context {
	return c.Context()
}

// SetValue sets the value for the inner [context.Context],
// which is useful for passing a value between middleware functions
func (c *Context) SetValue(key, val any) {
	ctx := c.Context()

	withVal := context.WithValue(ctx, key, val)
	c.ctx = withVal
}

// Value returns the value from the inner [context.Context],
// associated with the key
func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}

// MustGetMessage returns a pointer to the update's [Message].
// It is safe to call for MustGetMessage if the handler is subscribed to [UpdateTypeMessage],
// the return value is always non-nil.
// Otherwise, it will panic.
func (c *Context) MustGetMessage() *Message {
	if c.updType != UpdateTypeMessage {
		panic("calling MustMessage() when it is known that Update doesn't have any messages")
	}
	if c.upd.Message == nil {
		panic("well, it is a message, but somehow it is not") // never happens (i think)
	}
	return c.upd.Message
}

// GetMessage returns a pointer to the update's [Message].
// If you're not sure about the update type, it is safe to use GetMessage and check for a nil value.
// Unlike [MustGetMessage], it will never panic.
func (c *Context) GetMessage() *Message {
	return c.upd.Message
}
