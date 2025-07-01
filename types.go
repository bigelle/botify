package botify

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// This object represents a message.
type Message struct {
	// Unique message identifier inside this chat.
	// In specific instances (e.g., message containing a video sent to a big chat),
	// the server might automatically schedule a message instead of sending it immediately.
	// In such cases, this field will be 0 and the relevant message will be unusable until it is actually sent
	MessageId int `json:"message_id"`
	// Date the message was sent in Unix time. It is always a positive number, representing a valid date.
	Date int `json:"date"`
	// Chat the message belongs to
	Chat Chat `json:"chat"`

	// Optional. Unique identifier of a message thread to which the message belongs; for supergroups only
	// MessageThreadId *int `json:"message_thread_id,omitempty"`
	// Optional. Sender of the message; may be empty for messages sent to channels.
	// For backward compatibility, if the message was sent on behalf of a chat,
	// the field contains a fake sender user in non-channel chats
	From *User `json:"from,omitempty"`
	// Optional. Sender of the message when sent on behalf of a chat.
	// For example, the supergroup itself for messages sent by its anonymous administrators or
	// a linked channel for messages automatically forwarded to the channel's discussion group.
	// For backward compatibility, if the message was sent on behalf of a chat,
	// the field from contains a fake sender user in non-channel chats.
	SenderChat *Chat `json:"sender_chat,omitempty"`
	// Optional. If the sender of the message boosted the chat, the number of boosts added by the user
	SenderBoostCount *int `json:"sender_boost_count,omitempty"`
	// Optional. The bot that actually sent the message on behalf of the business account.
	// Available only for outgoing messages sent on behalf of the connected business account.
	SenderBusinessBot *User `json:"sender_business_bot,omitempty"`
	// Optional. Unique identifier of the business connection from which the message was received.
	// If non-empty, the message belongs to a chat of the corresponding business account that is
	// independent from any potential bot chat which might share the same identifier.
	BusinessConnectionId *string `json:"business_connection_id,omitempty"`
	// Optional. Information about the original message for forwarded messages
	// ForwardOrigin *MessageOrigin `json:"forward_origin,omitempty"`
	// Optional. True, if the message is sent to a forum topic
	IsTopicMessage *bool `json:"is_topic_message,omitempty"`
	// Optional. True, if the message is a channel post that was automatically forwarded to the connected discussion group
	IsAutomaticForward *bool `json:"is_automatic_forward,omitempty"`
	// Optional. For replies in the same chat and message thread, the original message.
	// Note that the Message object in this field will not contain further reply_to_message fields even if it itself is a reply.
	ReplyToMessage *Message `json:"reply_to_message,omitempty"`
	// Optional. Information about the message that is being replied to, which may come from another chat or forum topic
	// ExternalReply *ExternalReplyInfo `json:"external_reply,omitempty"`
	// Optional. For replies that quote part of the original message, the quoted part of the message
	// Quote *TextQuote `json:"quote,omitempty"`
	// Optional. For replies to a story, the original story
	// ReplyToStory *Story `json:"reply_to_story,omitempty"`
	// Optional. Bot through which the message was sent
	ViaBot *User `json:"via_bot,omitempty"`
	// Optional. Date the message was last edited in Unix time
	EditDate *int `json:"edit_date,omitempty"`
	// Optional. True, if the message can't be forwarded
	HasProtectedContent *bool `json:"has_protected_content,omitempty"`
	// Optional. True, if the message was sent by an implicit action, for example,
	// as an away or a greeting business message, or as a scheduled message
	IsFromOffline *bool `json:"is_from_offline,omitempty"`
	// Optional. The unique identifier of a media message group this message belongs to
	MediaGroupId *string `json:"media_group_id,omitempty"`
	// Optional. Signature of the post author for messages in channels, or the custom title of an anonymous group administrator
	AuthorSignature *string `json:"author_signature,omitempty"`
	// Optional. The number of Telegram Stars that were paid by the sender of the message to send it
	PaidStarCount *int `json:"paid_star_count,omitempty"`
	// Optional. For text messages, the actual UTF-8 text of the message
	Text *string `json:"text,omitempty"`
	// Optional. For text messages, special entities like usernames, URLs, bot commands, etc. that appear in the text
	// Entities *[]MessageEntity `json:"entities,omitempty"`
	// Optional. Options used for link preview generation for the message,
	// if it is a text message and link preview options were changed
	// LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	// Optional. Unique identifier of the message effect added to the message
	EffectId *string `json:"effect_id,omitempty"`
	// Optional. Message is an animation, information about the animation.
	// For backward compatibility, when this field is set, the document field will also be set
	// Animation *Animation `json:"animation,omitempty"`
	// Optional. Message is an audio file, information about the file
	// Audio *Audio `json:"audio,omitempty"`
	// Optional. Message is a general file, information about the file
	// Document *Document `json:"document,omitempty"`
	// Optional. Message contains paid media; information about the paid media
	// PaidMedia *PaidMediaInfo `json:"paid_media,omitempty"`
	// Optional. Message is a photo, available sizes of the photo
	// Photo *[]PhotoSize `json:"photo,omitempty"`
	// Optional. Message is a sticker, information about the sticker
	// Sticker *Sticker `json:"sticker,omitempty"`
	// Optional. Message is a forwarded story
	// Story *Story `json:"story,omitempty"`
	// Optional. Message is a video, information about the video
	// Video *Video `json:"video,omitempty"`
	// Optional. Message is a video note, information about the video message
	// VideoNote *VideoNote `json:"video_note,omitempty"`
	// Optional. Message is a voice message, information about the file
	// Voice *Voice `json:"voice,omitempty"`
	// Optional. Caption for the animation, audio, document, paid media, photo, video or voice
	// Caption *string `json:"caption,omitempty"`
	// Optional. For messages with a caption, special entities like usernames, URLs, bot commands, etc. that appear in the caption
	// CaptionEntities *[]MessageEntity `json:"caption_entities,omitempty"`
	// Optional. True, if the caption must be shown above the message media
	ShowCaptionAboveMedia *bool `json:"show_caption_above_media,omitempty"`
	// Optional. True, if the message media is covered by a spoiler animation
	HasMediaSpoiler *bool `json:"has_media_spoiler,omitempty"`
	// Optional. Message is a shared contact, information about the contact
	// Contact *Contact `json:"contact,omitempty"`
	// Optional. Message is a dice with random value
	// Dice *Dice `json:"dice,omitempty"`
	// Optional. Message is a game, information about the game. More about games » https://core.telegram.org/bots/api#games
	// Game *Game `json:"game,omitempty"`
	// Optional. Message is a native poll, information about the poll
	// Poll *Poll `json:"poll,omitempty"`
	// Optional. Message is a venue, information about the venue. For backward compatibility, when this field is set, the location field will also be set
	// Venue *Venue `json:"venue,omitempty"`
	// Optional. Message is a shared location, information about the location
	// Location *Location `json:"location,omitempty"`
	// Optional. New members that were added to the group or supergroup and information about them (the bot itself may be one of these members)
	NewChatMembers *[]User `json:"new_chat_members,omitempty"`
	// Optional. A member was removed from the group, information about them (this member may be the bot itself)
	LeftChatMember *User `json:"left_chat_member,omitempty"`
	// Optional. A chat title was changed to this value
	NewChatTitle *string `json:"new_chat_title,omitempty"`
	// Optional. A chat photo was change to this value
	// NewChatPhoto *[]PhotoSize `json:"new_chat_photo,omitempty"`
	// Optional. Service message: the chat photo was deleted
	DeleteChatPhoto *bool `json:"delete_chat_photo,omitempty"`
	// Optional. Service message: the group has been created
	GroupChatCreated *bool `json:"group_chat_created,omitempty"`
	// Optional. Service message: the supergroup has been created.
	// This field can't be received in a message coming through updates, because bot can't be a member of a supergroup when it is created.
	// It can only be found in reply_to_message if someone replies to a very first message in a directly created supergroup.
	SuperGroupCreated *bool `json:"super_group_created,omitempty"`
	// Optional. Service message: the channel has been created.
	// This field can't be received in a message coming through updates, because bot can't be a member of a channel when it is created.
	// It can only be found in reply_to_message if someone replies to a very first message in a channel.
	ChannelChatCreated *bool `json:"channel_chat_created"`
	// Optional. Service message: auto-delete timer settings changed in the chat
	// MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	// Optional. The group has been migrated to a supergroup with the specified identifier.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	MigrateToChatId *int64 `json:"migrate_to_chat_id,omitempty"`
	// Optional. The supergroup has been migrated from a group with the specified identifier.
	// This number may have more than 32 significant bits and some programming languages may have difficulty/silent defects in interpreting it.
	// But it has at most 52 significant bits, so a signed 64-bit integer or double-precision float type are safe for storing this identifier.
	MigrateFromChatId *int64 `json:"migrate_from_chat_id,omitempty"`
	// Optional. Specified message was pinned. Note that the Message object in this field
	// will not contain further reply_to_message fields even if it itself is a reply.
	// PinnedMessage *MaybeInaccessibleMessage `json:"pinned_message,omitempty"`
	// Optional. Message is an invoice for a payment, information about the invoice.
	// More about payments » https://core.telegram.org/bots/api#payments
	// Invoice *Invoice `json:"invoice,omitempty"`
	// Optional. Message is a service message about a successful payment, information about the payment.
	// More about payments » https://core.telegram.org/bots/api#payments
	// SuccessfulPayment *SuccessfulPayment `json:"successful_payment,omitempty"`
	// Optional. Message is a service message about a refunded payment, information about the payment.
	// More about payments » https://core.telegram.org/bots/api#payments
	// RefundedPayment *RefundedPayment `json:"refunded_payment,omitempty"`
	// Optional. Service message: users were shared with the bot
	// UsersShared *UsersShared `json:"users_shared,omitempty"`
	// Optional. Service message: a chat was shared with the bot
	// ChatShared *ChatShared `json:"chat_shared,omitempty"`
	// Optional. Service message: a regular gift was sent or received
	// Gift *GiftInfo `json:"gift,omitempty"`
	// Optional. Service message: a unique gift was sent or received
	// UniqueGift *UniqueGiftInfo `json:"unique_gift,omitempty"`
	// Optional. The domain name of the website on which the user has logged in.
	// More about Telegram Login » https://core.telegram.org/widgets/login
	// ConnectedWebsite *string `json:"connected_website,omitempty"`
	// Optional. Service message: the user allowed the bot to write messages after adding it to the attachment or side menu,
	// launching a Web App from a link, or accepting an explicit request from a Web App sent by the method requestWriteAccess
	// WriteAccessAllowed *WriteAccessAllowed `json:"write_access_allowed,omitempty"`
	// Optional. Telegram Passport data
	// PassportData *PassportData `json:"passport_data,omitempty"`
	// Optional. Service message. A user in the chat triggered another user's proximity alert while sharing Live Location.
	// ProximityAlertTriggered *ProximityAlertTriggered `json:"proximity_alert_triggered,omitempty"`
	// Optional. Service message: user boosted the chat
	// BoostAdded *ChatBoostAdded `json:"boost_added,omitempty"`
	// Optional. Service message: chat background set
	// ChatBackgroundSet *ChatBackground `json:"chat_background_set,omitempty"`
	// Optional. Service message: forum topic created
	// ForumTopicCreated *ForumTopicCreated `json:"forum_topic_created,omitempty"`
	// Optional. Service message: forum topic edited
	// ForumTopicEdited *ForumTopicEdited `json:"forum_topic_edited,omitempty"`
	// Optional. Service message: forum topic closed
	// ForumTopicClosed *ForumTopicClosed `json:"forum_topic_closed,omitempty"`
	// Optional. Service message: forum topic reopened
	// ForumTopicReopened *ForumTopicReopened `json:"forum_topic_reopened,omitempty"`
	// Optional. Service message: the 'General' forum topic hidden
	// GeneralForumTopicHidden *GeneralForumTopicHidden `json:"general_forum_topic_hidden,omitempty"`
	// Optional. Service message: the 'General' forum topic unhidden
	// GeneralForumTopicUnhidden *GeneralForumTopicUnhidden `json:"general_forum_topic_unhidden,omitempty"`
	// Optional. Service message: a scheduled giveaway was created
	// GiveawayCreated *GiveawayCreated `json:"giveaway_created,omitempty"`
	// Optional. The message is a scheduled giveaway message
	// Giveaway *Giveaway `json:"giveaway,omitempty"`
	// Optional. A giveaway with public winners was completed
	// GiveawayWinners *GiveawayWinners `json:"giveaway_winners,omitempty"`
	// Optional. Service message: a giveaway without public winners was completed
	// GiveawayCompleted *GiveawayCompleted `json:"giveaway_completed,omitempty"`
	// Optional. Service message: the price for paid messages has changed in the chat
	// PaidMessagePriceChanged *PaidMessagePriceChanged `json:"paid_message_price_changed,omitempty"`
	// Optional. Service message: video chat scheduled
	// VideoChatScheduled *VideoChatScheduled `json:"video_chat_scheduled,omitempty"`
	// Optional. Service message: video chat started
	// VideoChatStarted *VideoChatStarted `json:"video_chat_started,omitempty"`
	// Optional. Service message: video chat ended
	// VideoChatEnded *VideoChatEnded `json:"video_chat_ended,omitempty"`
	// Optional. Service message: new participants invited to a video chat
	// VideoChatParticipantsInvited *VideoChatParticipantsInvited `json:"video_chat_participants_invited,omitempty"`
	// Optional. Service message: data sent by a Web App
	// WebAppData *WebAppData `json:"web_app_data,omitempty"`
	// Optional. Inline keyboard attached to the message. login_url buttons are represented as ordinary url buttons.
	// ReplyMarkup *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

type User struct{}

type Chat struct {
	ID        int     `json:"id"`
	Type      string  `json:"type"`
	Title     *string `json:"title"`
	Username  *string `json:"username"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	IsForum   *bool   `json:"is_forum"`
}

type ChatBoostSource struct {
	Source string `json:"source"`
	raw    json.RawMessage
}

func (c *ChatBoostSource) UnmarshalJSON(data []byte) error {
	type Alias ChatBoostSource
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(c),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	c.raw = data
	return nil
}

func (s *ChatBoostSource) AsPremium() (*ChatBoostSourcePremium, error) {
	if s.Source != "premium" {
		return nil, fmt.Errorf("failed conversion; expected source premium, got %s", s.Source)
	}

	var result ChatBoostSourcePremium
	if err := json.NewDecoder(bytes.NewReader(s.raw)).Decode(&result); err != nil {
		return nil, fmt.Errorf("converting to ChatBoostSourcePremium: %w", err)
	}

	return &result, nil
}

func (s *ChatBoostSource) AsGiftCode() (*ChatBoostSourceGiftCode, error) {
	if s.Source != "gift_code" {
		return nil, fmt.Errorf("failed conversion; expected source gift_code, got %s", s.Source)
	}

	var result ChatBoostSourceGiftCode
	if err := json.NewDecoder(bytes.NewReader(s.raw)).Decode(&result); err != nil {
		return nil, fmt.Errorf("converting to ChatBoostSourcePremium: %w", err)
	}

	return &result, nil
}

func (s *ChatBoostSource) AsGiveaway() (*ChatBoostSourceGiveaway, error) {
	if s.Source != "giveaway" {
		return nil, fmt.Errorf("failed conversion; expected source giveaway, got %s", s.Source)	}

	var result ChatBoostSourceGiveaway
	if err := json.NewDecoder(bytes.NewReader(s.raw)).Decode(&result); err != nil {
		return nil, fmt.Errorf("converting to ChatBoostSourcePremium: %w", err)	}

	return &result, nil
}

type ChatBoostSourcePremium struct {
	User User `json:"user"`
}

type ChatBoostSourceGiftCode struct {
	User User `json:"user"`
}

type ChatBoostSourceGiveaway struct {
	GiveawayMessageID int   `json:"giveaway_message_id"`
	User              *User `json:"user"`
	PrizeStarCount    *int  `json:"prize_star_count"`
	IsUnclaimed       *bool `json:"is_unclaimed"`
}
