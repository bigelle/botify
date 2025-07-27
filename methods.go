package botify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
)

type APIMethod interface {
	ContentType() string
	APIEndpoint() string
	Payload() (io.Reader, error)
}

type contentTyperJSON struct{}

func (c contentTyperJSON) ContentType() string {
	return "application/json"
}

// methodWithNoParams is used to send a request that requires no parameters,
// meaning there is no request body and it does not require Content-Type header.
type methodWithNoParams string

func (m methodWithNoParams) ContentType() string {
	return ""
}

func (m methodWithNoParams) APIEndpoint() string {
	return string(m)
}

func (m methodWithNoParams) Payload() (io.Reader, error) {
	return nil, nil
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
	contentTyperJSON
}

func (m *GetUpdates) APIEndpoint() string {
	return "getUpdates"
}

func (m *GetUpdates) Payload() (io.Reader, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 512))
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(m); err != nil {
		return nil, fmt.Errorf("encoding getUpdates payload: %w", err)
	}
	return buf, nil
}

type SetWebhook struct {
	URL                string    `json:"url"`
	Certificate        InputFile `json:"certificate,omitempty"`
	IPAddress          string    `json:"ip_address,omitempty"`
	MaxConnections     int       `json:"max_connections,omitempty"`
	AllowedUpdates     *[]string `json:"allowed_updates,omitempty"`
	DropPendingUpdates bool      `json:"drop_pending_updates,omitempty"`
	SecretToken        string    `json:"secret_token,omitempty"`

	ct string
}

func (m *SetWebhook) ContentType() string {
	// it's always multipart so if ct is empty we fallback
	if m.ct == "" && m.Certificate != nil {
		return "multipart/form-data"
	}
	// it means there's no Certificate so it is safe to send a JSON
	if m.ct == "" {
		return "application/json"
	}
	return m.ct
}

func (m *SetWebhook) APIEndpoint() string {
	return "setWebhook"
}

func (m *SetWebhook) Payload() (io.Reader, error) {
	if _, ok := m.Certificate.(InputFileRemote); ok {
		return nil, fmt.Errorf("can't upload a certificate from a remote source; use a local file")
	}

	if m.Certificate == nil {
		// then there's no need to send multipart
		buf := bytes.NewBuffer(make([]byte, 0, 512))
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)

		if err := enc.Encode(m); err != nil {
			return nil, fmt.Errorf("encoding setWebhook payload as JSON: %w", err)
		}
		return buf, nil
	}
	return m.multipart()
}

func (m *SetWebhook) multipart() (io.Reader, error) {
	var err error

	buf := bytes.NewBuffer(make([]byte, 0, 4*1024))
	mw := multipart.NewWriter(buf)
	m.ct = mw.FormDataContentType()
	defer mw.Close()

	if err = mw.WriteField("url", m.URL); err != nil {
		return nil, fmt.Errorf("writing form field: %w", err)
	}

	cert := m.Certificate.(InputFileLocal)
	part, err := mw.CreateFormFile("certificate", cert.Name)
	if err != nil {
		return nil, fmt.Errorf("creating form file: %w", err)
	}
	_, err = io.Copy(part, cert.Data)
	if err != nil {
		return nil, fmt.Errorf("writing form file: %w", err)
	}

	if m.IPAddress != "" {
		err = mw.WriteField("ip_address", m.IPAddress)
		if err != nil {
			return nil, fmt.Errorf("writing form field: %w", err)
		}
	}

	if m.MaxConnections != 0 {
		err = mw.WriteField("max_connections", fmt.Sprint(m.MaxConnections))
		if err != nil {
			return nil, fmt.Errorf("writing form field: %w", err)
		}
	}

	if m.AllowedUpdates != nil {
		var b []byte
		b, err = json.Marshal(m.AllowedUpdates)
		if err != nil {
			return nil, fmt.Errorf("encoding allowed updates as JSON: %w", err)
		}
		err = mw.WriteField("allowed_updates", string(b))
		if err != nil {
			return nil, fmt.Errorf("writing form field: %w", err)
		}
	}

	if m.DropPendingUpdates {
		err = mw.WriteField("drop_pending_updates", fmt.Sprint(m.DropPendingUpdates))
		if err != nil {
			return nil, fmt.Errorf("writing form field: %w", err)
		}
	}

	if m.SecretToken != "" {
		err = mw.WriteField("secret_token", m.SecretToken)
		if err != nil {
			return nil, fmt.Errorf("writing form field: %w", err)
		}
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
	contentTyperJSON
}

func (m *SendMessage) Method() string {
	return "sendMessage"
}

func (m *SendMessage) Payload() (io.Reader, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(m); err != nil {
		return nil, fmt.Errorf("encoding sendMessage payload: %w", err)
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
	if local, ok := m.Photo.(InputFileLocal); ok {
		buf := bytes.NewBuffer(make([]byte, 4*1024))
		mw := multipart.NewWriter(buf)
		defer mw.Close()
		m.ct = mw.FormDataContentType()

		part, err := mw.CreateFormFile("photo", local.Name)
		if err != nil {
			return nil, fmt.Errorf("creating form file: %w", err)
		}
		_, err = io.Copy(part, local.Data)
		if err != nil {
			return nil, fmt.Errorf("writing form file: %w", err)
		}

		err = mw.WriteField("chat_id", m.ChatID)
		if err != nil {
			return nil, fmt.Errorf("writing form field: %w", err)
		}

		return buf, nil
	}

	m.ct = "application/json"
	buf := bytes.NewBuffer(make([]byte, 1024))
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(m); err != nil {
		return nil, fmt.Errorf("encoding sendPhoto JSON payload: %w", err)
	}
	return buf, nil
}

type GetMyCommands struct {
	Scope        BotCommandScope `json:"scope,omitempty"`
	LanguageCode string          `json:"language_code,omitempty"`
	contentTyperJSON
}

func (m *GetMyCommands) APIEndpoint() string {
	return "getMyCommands"
}

func (m *GetMyCommands) Payload() (io.Reader, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 128))
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(m); err != nil {
		return nil, fmt.Errorf("encoding getMyCommands JSON payload: %w", err)
	}
	return buf, nil
}

type SetMyCommands struct {
	Commands     []BotCommand    `json:"commands"`
	Scope        BotCommandScope `json:"scope,omitempty"`
	LanguageCode string          `json:"language_code,omitempty"`
	contentTyperJSON
}

func (m *SetMyCommands) APIEndpoint() string {
	return "setMyCommands"
}

func (m *SetMyCommands) Payload() (io.Reader, error) {
	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)

	if err := enc.Encode(m); err != nil {
		return nil, fmt.Errorf("encoding setMyCommands JSON payload: %w", err)
	}
	return buf, nil
}
