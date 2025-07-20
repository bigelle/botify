package botify

import (
	"context"
	"slices"
	"strings"
	"time"
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

func (u *Update) UpdateType() string {
	checks := []struct {
		condition bool
		result    string
	}{
		{u.Message != nil && u.Message.IsCommand(), func() string {
			cmd, _ := u.Message.GetCommand()
			return cmd
		}()},
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

type HandlerFunc func(ctx *Context)

func ChainHandlers(handlers ...HandlerFunc) HandlerFunc {
	return func(ctx *Context) {
		for _, handler := range handlers {
			handler(ctx)
		}
	}
}

type Context struct {
	bot *Bot

	updType        string
	upd            *Update
	sendedRequests []RequestInfo // for logging

	ctx context.Context
}

type RequestInfo struct {
	Method      string
	ContentType string
	Duration    time.Duration
}

// returns a read-only copy of the currently active bot
func (c *Context) Bot() Bot {
	return *c.bot
}

func (c *Context) Send(obj APIMethod) (*APIResponse, error) {
	start := time.Now()
	resp, err := c.bot.Sender.SendWithContext(c.ctx, obj)
	end := time.Since(start) // used mostly for measuring the network latencies
	// but it also measures the time took by writing the request body
	ct, _, _ := strings.Cut(obj.ContentType(), ";")

	c.sendedRequests = slices.Grow(c.sendedRequests, 1)
	c.sendedRequests = append(c.sendedRequests, RequestInfo{
		Method:      obj.Method(),
		ContentType: ct,
		Duration:    end,
	})

	return resp, err
}

func (c *Context) SendRaw(method string, obj any) (*APIResponse, error) {
	start := time.Now()
	resp, err := c.bot.Sender.SendRawWithContext(c.ctx, method, obj)
	end := time.Since(start)

	c.sendedRequests = slices.Grow(c.sendedRequests, 1)
	c.sendedRequests = append(c.sendedRequests, RequestInfo{
		Method:      method,
		ContentType: "application/json", // i mean, its not possible to send anything besides json
		Duration:    end,
	})

	return resp, err
}

// Use it to make sure that you're working with the expected Update type
func (c *Context) UpdateType() string {
	return c.updType
}

func (c *Context) UpdateID() int {
	return c.upd.UpdateID
}

func (c *Context) SendedRequests() []RequestInfo {
	return c.sendedRequests
}

func (c *Context) Context() context.Context {
	return c.Context()
}

func (c *Context) SetValue(key, val any) {
	ctx := c.Context()

	withVal := context.WithValue(ctx, key, val)

	c.ctx = withVal
}

func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}

// It is safe to call for MustBindMessage() if the handler is subscribed to UpdateTypeMessage,
// the return value is always non-nil.
// Otherwise, it will panic.
func (c *Context) MustBindMessage() *Message {
	if c.updType != UpdateTypeMessage {
		panic("calling MustMessage() when it is known that Update doesn't have any messages")
	}

	if c.upd.Message == nil {
		panic("well, it is a message, but somehow it is not")
	}

	return c.upd.Message
}

// If you're not sure about the update type, it is safe to use TryMessage() and check for a nil value.
// It will never panic.
func (c *Context) BindMessage() *Message {
	return c.upd.Message
}
