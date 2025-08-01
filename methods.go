package botify

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/bigelle/formy"
)

// APIMethod describes how the request would be written
// and which API method it belongs to
type APIMethod interface {
	// APIEndpoint returns API method name.
	// E.g. for sendMessage it should return "sendMessage",
	// without the leading or trailing "/", case insensitive
	APIEndpoint() string
	// WritePayload writes the struct into body, which would be used while sending the request.
	// It must return a valid "Content-Type" header value and any write errors
	WritePayload(body io.Writer) (contentType string, err error)
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

func (m methodWithNoParams) WritePayload(_ io.Writer) (string, error) {
	return "", nil
}

const (
	GetWebhookInfo methodWithNoParams = "getWebhookInfo"
	Close          methodWithNoParams = "close"
)

type GetUpdates struct {
	Offset         int      `json:"offset,omitempty"`
	Limit          int      `json:"limit,omitempty"`
	Timeout        int      `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitzero"`
}

func (m *GetUpdates) APIEndpoint() string {
	return "getUpdates"
}

func (m *GetUpdates) WritePayload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}

type SetWebhook struct {
	URL                string    `json:"url"`
	Certificate        InputFile `json:"certificate,omitempty"`
	IPAddress          string    `json:"ip_address,omitempty"`
	MaxConnections     int       `json:"max_connections,omitempty"`
	AllowedUpdates     []string  `json:"allowed_updates,omitzero"`
	DropPendingUpdates bool      `json:"drop_pending_updates,omitempty"`
	SecretToken        string    `json:"secret_token,omitempty"`
}

func (m *SetWebhook) APIEndpoint() string {
	return "setWebhook"
}

func (m *SetWebhook) WritePayload(body io.Writer) (string, error) {
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

type DeleteWebhook struct {
	DropPendingUpdates bool `json:"drop_pending_updates,omitempty"`
}

func (m *DeleteWebhook) APIEndpoint() string {
	return "deleteWebhook"
}

func (m *DeleteWebhook) WritePayload(w io.Writer) (string, error) {
	return jsonPayload(m, w)
}

type SendMessage struct {
	ChatId               string           `json:"chat_id"`
	Text                 string           `json:"text"`
	BusinessConnectionId string           `json:"business_connection_id,omitempty"`
	MessageThreadId      int              `json:"message_thread_id,omitempty"`
	ParseMode            string           `json:"parse_mode,omitempty"`
	Entities             *MessageEntity   `json:"entities,omitempty"`
	DisableNotification  bool             `json:"disable_notification,omitempty"`
	ProtectContent       bool             `json:"protect_content,omitempty"`
	AllowPaidBroadcast   bool             `json:"allow_paid_broadcast,omitempty"`
	MessageEffectId      string           `json:"message_effect_id,omitempty"`
	ReplyParameters      *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup          *ReplyMarkup     `json:"reply_markup,omitempty"`
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

func (m *GetMyCommands) WritePayload(body io.Writer) (string, error) {
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

func (m *SetMyCommands) WritePayload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}
