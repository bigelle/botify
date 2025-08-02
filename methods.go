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
	// Getting updates

	GetWebhookInfo methodWithNoParams = "getWebhookInfo"

	// Available methods

	GetMe          methodWithNoParams = "getMe"
	LogOut         methodWithNoParams = "logOut"
	Close          methodWithNoParams = "close"
)

/*
	BEGIN Getting updates TYPES
*/

type GetUpdates struct {
	Offset         int      `json:"offset,omitempty"`
	Limit          int      `json:"limit,omitempty"`
	Timeout        int      `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitzero"`
}

func (m GetUpdates) APIEndpoint() string {
	return "getUpdates"
}

func (m GetUpdates) WritePayload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}

type SetWebhook struct {
	URL                string    `json:"url"`
	Certificate        InputFile `json:"certificate,omitempty"`
	IPAddress          string    `json:"ip_address,omitempty"`
	MaxConnections     int       `json:"max_connections,omitempty"`
	AllowedUpdates     []string  `json:"allowed_updates,omitempty"`
	DropPendingUpdates bool      `json:"drop_pending_updates,omitempty"`
	SecretToken        string    `json:"secret_token,omitempty"`
}

func (m SetWebhook) APIEndpoint() string {
	return "setWebhook"
}

func (m SetWebhook) WritePayload(body io.Writer) (string, error) {
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
		WriteStringCond("ip_address", m.IPAddress, notEmptyString(m.IPAddress)).
		WriteIntCond("max_connections", m.MaxConnections, notEmptyInt(m.MaxConnections)).
		WriteJSONCond("allowed_updates", m.AllowedUpdates, notEmptySlice(m.AllowedUpdates)).
		WriteBoolCond("drop_pending_updates", m.DropPendingUpdates, func() bool { return m.DropPendingUpdates }).
		WriteStringCond("secret_token", m.SecretToken, notEmptyString(m.SecretToken))

	return mw.FormDataContentType(), mw.Close()
}

type DeleteWebhook struct {
	DropPendingUpdates bool `json:"drop_pending_updates,omitempty"`
}

func (m DeleteWebhook) APIEndpoint() string {
	return "deleteWebhook"
}

func (m DeleteWebhook) WritePayload(w io.Writer) (string, error) {
	return jsonPayload(m, w)
}

/*
	BEGING Available methods TYPES
*/

type SendMessage struct {
	ChatID               string           `json:"chat_id"`
	Text                 string           `json:"text"`
	BusinessConnectionID string           `json:"business_connection_id,omitempty"`
	MessageThreadID      int              `json:"message_thread_id,omitempty"`
	ParseMode            string           `json:"parse_mode,omitempty"`
	Entities             *MessageEntity   `json:"entities,omitempty"`
	DisableNotification  bool             `json:"disable_notification,omitempty"`
	ProtectContent       bool             `json:"protect_content,omitempty"`
	AllowPaidBroadcast   bool             `json:"allow_paid_broadcast,omitempty"`
	MessageEffectId      string           `json:"message_effect_id,omitempty"`
	ReplyParameters      *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup          *ReplyMarkup     `json:"reply_markup,omitempty"`
}

func (m SendMessage) APIEndpoint() string {
	return "sendMessage"
}

func (m SendMessage) WritePayload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}

type ForwardMessage struct {
	ChatID              string `json:"chat_id"`
	FromChatID          string `json:"from_chat_id"`
	MessageID           int    `json:"message_id"`
	MessageThreadID     int    `json:"message_thread_id,omitempty"`
	VideoStartTimestamp int    `json:"video_start_timestamp,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
}

func (m ForwardMessage) APIEndpoint() string {
	return "forwardMessage"
}

func (m ForwardMessage) WritePayload(w io.Writer) (string, error) {
	return jsonPayload(m, w)
}

type ForwardMessages struct {
	ChatID              string `json:"chat_id"`
	FromChatID          string `json:"from_chat_id"`
	MessageIDs          []int  `json:"message_ids"`
	MessageThreadID     int    `json:"message_thread_id,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
}

func (m ForwardMessages) APIEndpoint() string {
	return "forwardMessages"
}

func (m ForwardMessages) WritePayload(w io.Writer) (string, error) {
	return jsonPayload(m, w)
}

type CopyMessage struct {
	ChatID                string           `json:"chat_id"`
	FromChatID            string           `json:"from_chat_id"`
	MessageID             int              `json:"message_id"`
	MessageThreadID       int              `json:"message_thread_id,omitempty"`
	VideoStartTimestamp   int              `json:"video_start_timestamp,omitempty"`
	Caption               string           `json:"caption,omitempty"`
	ParseMode             string           `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity  `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool             `json:"show_caption_above_media,omitempty"`
	DisableNotification   bool             `json:"disable_notification,omitempty"`
	ProtectContent        bool             `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool             `json:"allow_paid_broadcast,omitempty"`
	ReplyParameters       *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup           ReplyMarkup      `json:"reply_markup,omitempty"`
}

func (m CopyMessage) APIEndpoint() string {
	return "copyMessage"
}

func (m CopyMessage) WritePayload(w io.Writer) (string, error) {
	return jsonPayload(m, w)
}

type CopyMessages struct {
	ChatID              string `json:"chat_id"`
	FromChatID          string `json:"from_chat_id"`
	MessageID           int    `json:"message_id"`
	MessageThreadID     int    `json:"message_thread_id,omitempty"`
	ParseMode           string `json:"parse_mode"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content"`
	RemoveCaption       bool   `json:"remove_caption,omitempty"`
}

func (m CopyMessages) APIEndpoint() string {
	return "copyMessages"
}

func (m CopyMessages) WritePayload(w io.Writer) (string, error) {
	return jsonPayload(m, w)
}

type SendPhoto struct {
	ChatID                string           `json:"chat_id"`
	Photo                 InputFile        `json:"photo"`
	BusinessConnectionID  string           `json:"business_connection_id,omitempty"`
	MessageThreadID       int              `json:"message_thread_id,omitempty"`
	Caption               string           `json:"caption,omitempty"`
	ParseMode             string           `json:"parse_mode,omitempty"`
	CaptionEntities       []MessageEntity  `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia bool             `json:"show_caption_above_media,omitempty"`
	HasSpoiler            bool             `json:"has_spoiler,omitempty"`
	DisableNotification   bool             `json:"disable_notification,omitempty"`
	ProtectContent        bool             `json:"protect_content,omitempty"`
	AllowPaidBroadcast    bool             `json:"allow_paid_broadcast,omitempty"`
	MessageEffectID       string           `json:"message_effect_id,omitempty"`
	ReplyParameters       *ReplyParameters `json:"reply_parameters,omitempty"`
	ReplyMarkup           ReplyMarkup      `json:"reply_markup,omitempty"`
}

func (m SendPhoto) APIEndpoint() string {
	return "sendPhoto"
}

func (m SendPhoto) WritePayload(body io.Writer) (string, error) {
	if _, ok := m.Photo.(InputFileRemote); ok {
		return jsonPayload(m, body)
	}

	photo := m.Photo.(InputFileLocal)
	mw := formy.NewWriter(body).
		WriteFile("photo", photo.Name, photo.Data).
		WriteString("chat_id", m.ChatID).
		WriteStringCond("business_connection_id", m.BusinessConnectionID, notEmptyString(m.BusinessConnectionID)).
		WriteIntCond("message_thread_id", m.MessageThreadID, notEmptyInt(m.MessageThreadID)).
		WriteStringCond("caption", m.Caption, notEmptyString(m.Caption)).
		WriteJSONCond("message_entities", m.CaptionEntities, notEmptySlice(m.CaptionEntities)).
		WriteBoolCond("show_caption_above_media", m.ShowCaptionAboveMedia, func() bool { return m.ShowCaptionAboveMedia }).
		WriteBoolCond("has_spoiler", m.HasSpoiler, func() bool { return m.HasSpoiler }).
		WriteBoolCond("disable_notification", m.DisableNotification, func() bool { return m.DisableNotification }).
		WriteBoolCond("protect_content", m.ProtectContent, func() bool { return m.ProtectContent }).
		WriteBoolCond("allow_paid_broadcast", m.AllowPaidBroadcast, func() bool { return m.AllowPaidBroadcast }).
		WriteStringCond("message_effect_id", m.MessageEffectID, notEmptyString(m.MessageEffectID)).
		WriteJSONCond("reply_parameters", m.ReplyParameters, func() bool { return m.ReplyParameters != nil }).
		WriteJSONCond("reply_markup", m.ReplyMarkup, func() bool { return m.ReplyMarkup != nil })

	return mw.FormDataContentType(), mw.Close()
}

type GetMyCommands struct {
	Scope        BotCommandScope `json:"scope,omitempty"`
	LanguageCode string          `json:"language_code,omitempty"`
}

func (m GetMyCommands) APIEndpoint() string {
	return "getMyCommands"
}

func (m GetMyCommands) WritePayload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}

type SetMyCommands struct {
	Commands     []BotCommand    `json:"commands"`
	Scope        BotCommandScope `json:"scope,omitempty"`
	LanguageCode string          `json:"language_code,omitempty"`
}

func (m SetMyCommands) APIEndpoint() string {
	return "setMyCommands"
}

func (m SetMyCommands) WritePayload(body io.Writer) (string, error) {
	return jsonPayload(m, body)
}
