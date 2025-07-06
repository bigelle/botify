package botify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/bigelle/botify/internal/reused"
)

type APIMethod interface {
	ContentType() string
	Method() string
	Payload() (io.Reader, error)
}

// MethodWithNoParams is used to send a request that requires no parameters,
// meaning there is no request body and it does not require Content-Type header.
type MethodWithNoParams string

func (m MethodWithNoParams) ContentType() string {
	return ""
}

func (m MethodWithNoParams) Method() string {
	return string(m)
}

func (m MethodWithNoParams) Payload() (io.Reader, error) {
	return nil, nil
}

const (
	GetWebhookInfo MethodWithNoParams = "getWebhookInfo"
)

type GetUpdates struct {
	Offset         int       `json:"offset"`
	Limit          int       `json:"limit"`
	Timeout        int       `json:"timeout"`
	AllowedUpdates *[]string `json:"allowed_updates"`
}

func (m *GetUpdates) ContentType() string {
	return "application/json"
}

func (m *GetUpdates) Method() string {
	return "getUpdates"
}

func (m *GetUpdates) Payload() (io.Reader, error) {
	b := reused.Buf()
	defer reused.PutBuf(b)

	buf := bytes.NewBuffer(*b)

	err := json.NewEncoder(buf).Encode(m)
	if err != nil {
		return nil, fmt.Errorf("encoding getUpdates payload: %w", err)
	}

	return buf, nil
}

// FIXME: commented fields

// Use this method to send text messages.
// On success, the sent [Message] is returned.
type SendMessage struct {
	// REQUIRED:
	// Unique identifier for the target chat or username of the target channel (in the format @channelusername)
	ChatId string `json:"chat_id"`
	// REQUIRED:
	// Text of the message to be sent, 1-4096 characters after entities parsing
	Text string `json:"text"`

	// Unique identifier of the business connection on behalf of which the message will be sent
	BusinessConnectionId *string `json:"business_connection_id,omitempty,"`
	// Unique identifier for the target message thread (topic) of the forum; for forum supergroups only
	MessageThreadId *int `json:"message_thread_id,omitempty,"`
	// Mode for parsing entities in the message text.
	// See https://core.telegram.org/bots/api#formatting-options for more details.
	ParseMode *string `json:"parse_mode,omitempty,"`
	// A JSON-serialized list of special entities that appear in message text, which can be specified instead of parse_mode
	// Entities *MessageEntity `json:"entities,omitempty,"`
	// Link preview generation options for the message
	// LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty,"`
	// Sends the message silently. Users will receive a notification with no sound.
	DisableNotification *bool `json:"disable_notification,omitempty,"`
	// Protects the contents of the sent message from forwarding and saving
	ProtectContent *bool `json:"protect_content,omitempty,"`
	// Pass True to allow up to 1000 messages per second, ignoring broadcasting limits for a fee of 0.1 Telegram Stars per message.
	// The relevant Stars will be withdrawn from the bot's balance
	AllowPaidBroadcast *bool `json:"allow_paid_broadcast,omitempty"`
	// Unique identifier of the message effect to be added to the message; for private chats only
	MessageEffectId *string `json:"message_effect_id,omitempty,"`
	// Description of the message to reply to
	// ReplyParameters *ReplyParameters `json:"reply_parameters,omitempty,"`
	// Additional interface options. A JSON-serialized object for an inline keyboard,
	// custom reply keyboard, instructions to remove a reply keyboard or to force a reply from the user
	// ReplyMarkup *ReplyMarkup `json:"reply_markup,omitempty,"`
}

func (m *SendMessage) ContentType() string {
	return "application/json"
}

func (m *SendMessage) Method() string {
	return "sendMessage"
}

func (m *SendMessage) Payload() (io.Reader, error) {
	b := reused.Buf()
	defer reused.PutBuf(b)

	buf := bytes.NewBuffer(*b)

	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	err := enc.Encode(m)
	if err != nil {
		return nil, fmt.Errorf("encoding getUpdates payload: %w", err)
	}

	return buf, nil
}

type SendPhoto struct {
	ChatID string    `json:"chat_id"`
	Photo  InputFile `json:"photo"`
	// TODO: other fields

	ct string
}

func (m *SendPhoto) ContentType() string {
	// if it returned an empty string, it means something went wrong in Payload() method
	return m.ct
}

func (m *SendPhoto) Method() string {
	return "sendPhoto"
}

func (m *SendPhoto) Payload() (io.Reader, error) {
	// TODO:
	return nil, nil
}
