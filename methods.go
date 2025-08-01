package botify

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/bigelle/formy"
)

type APIMethod interface {
	APIEndpoint() string
	Payload(body io.Writer) (contentType string, err error)
}

func jsonPayload(obj any, body io.Writer) (string, error) {
	enc := json.NewEncoder(body)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(obj); err != nil {
		return "", fmt.Errorf("encoding JSON payload: %w", err)
	}
	return "application/json", nil
}

// methodWithNoParams is used to send a request that requires no parameters,
// meaning there is no request body and it does not require Content-Type header.
type methodWithNoParams string

func (m methodWithNoParams) APIEndpoint() string {
	return string(m)
}

func (m methodWithNoParams) Payload(_ io.Writer) (string, error) {
	return "", nil
}

const (
	GetWebhookInfo methodWithNoParams = "getWebhookInfo"
	Close          methodWithNoParams = "close"
)

type GetUpdates struct {
	Offset         int       `json:"offset,omitempty"`
	Limit          int       `json:"limit,omitempty"`
	Timeout        int       `json:"timeout,omitempty"`
	AllowedUpdates *[]string `json:"allowed_updates,omitempty"`
}

func (m *GetUpdates) APIEndpoint() string {
	return "getUpdates"
}

func (m *GetUpdates) Payload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}

type SetWebhook struct {
	URL                string    `json:"url"`
	Certificate        InputFile `json:"certificate,omitempty"`
	IPAddress          string    `json:"ip_address,omitempty"`
	MaxConnections     int       `json:"max_connections,omitempty"`
	AllowedUpdates     *[]string `json:"allowed_updates,omitempty"`
	DropPendingUpdates bool      `json:"drop_pending_updates,omitempty"`
	SecretToken        string    `json:"secret_token,omitempty"`
}

func (m *SetWebhook) APIEndpoint() string {
	return "setWebhook"
}

func (m *SetWebhook) Payload(body io.Writer) (string, error) {
	if _, ok := m.Certificate.(InputFileRemote); ok {
		return "", fmt.Errorf("can't upload a certificate from a remote source; use a local file")
	}
	if m.Certificate == nil {
		// then there's no need to send multipart
		return jsonPayload(m, body)
	}

	cert := m.Certificate.(InputFileLocal)
	mw := formy.NewWriter(body).
		WriteString("url", m.URL).
		WriteFile("certificate", cert.Name, cert.Data).
		WriteString("ip_address", m.IPAddress).
		WriteInt("max_connections", m.MaxConnections).
		WriteJSON("allowed_updates", m.AllowedUpdates).
		WriteBool("drop_pending_updates", m.DropPendingUpdates).
		WriteString("secret_token", m.SecretToken)

	return mw.FormDataContentType(), mw.Close()
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
	Entities *MessageEntity `json:"entities,omitempty,"`

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

func (m *SendMessage) Method() string {
	return "sendMessage"
}

func (m *SendMessage) Payload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}

type SendPhoto struct {
	ChatID string    `json:"chat_id"`
	Photo  InputFile `json:"photo"`
	// TODO: other fields
}

func (m *SendPhoto) Method() string {
	return "sendPhoto"
}

func (m *SendPhoto) Payload(body io.Writer) (string, error) {
	if _, ok := m.Photo.(InputFileRemote); ok {
		return jsonPayload(m, body)
	}

	photo := m.Photo.(InputFileLocal)
	mw := formy.NewWriter(body).
		WriteFile("photo", photo.Name, photo.Data).
		WriteString("chat_id", m.ChatID)

	return mw.FormDataContentType(), mw.Close()
}

type GetMyCommands struct {
	Scope        BotCommandScope `json:"scope,omitempty"`
	LanguageCode string          `json:"language_code,omitempty"`
}

func (m *GetMyCommands) APIEndpoint() string {
	return "getMyCommands"
}

func (m *GetMyCommands) Payload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}

type SetMyCommands struct {
	Commands     []BotCommand    `json:"commands"`
	Scope        BotCommandScope `json:"scope,omitempty"`
	LanguageCode string          `json:"language_code,omitempty"`
}

func (m *SetMyCommands) APIEndpoint() string {
	return "setMyCommands"
}

func (m *SetMyCommands) Payload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}
