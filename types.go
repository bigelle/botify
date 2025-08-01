package botify

import (
	"encoding/json"
	"fmt"
	"io"
)

/*
	BEGIN Getting Updates TYPES
*/

type WebhookInfo struct {
	URL                          string    `json:"url"`
	HasCustomCertificate         *bool     `json:"has_custom_certificate"`
	PendingUpdateCount           *int      `json:"pending_update_count"`
	IPAddress                    *string   `json:"ip_address"`
	LastErrorDate                *int      `json:"last_error_date"`
	LastErrorMessage             *string   `json:"last_error_message"`
	LastSynchronizationErrorDate *int      `json:"last_synchronization_error_date"`
	MaxConnections               *int      `json:"max_connections"`
	AllowedUpdates               *[]string `json:"allowed_updates"`
}

/*
	BEGIN Available types TYPES
*/

type User struct {
	ID                      int     `json:"id"`
	FirstName               string  `json:"first_name"`
	IsBot                   bool    `json:"is_bot"`
	LastName                *string `json:"last_name,omitempty"`
	UserName                *string `json:"user_name,omitempty"`
	LanguageCode            *string `json:"language_code,omitempty"`
	CanJoinGroups           *bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages *bool   `json:"can_read_all_group_messages,omitempty"`
	SupportInlineQueries    *bool   `json:"support_inline_queries,omitempty"`
	IsPremium               *bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu   *bool   `json:"added_to_attachment_menu,omitempty"`
	CanConnectToBusiness    *bool   `json:"can_connect_to_business,omitempty"`
	HasMainWebApp           *bool   `json:"has_main_web_app,omitempty"`
}

type Chat struct {
	ID        int64   `json:"id"`
	Type      string  `json:"type"`
	Title     *string `json:"title,omitempty"`
	UserName  *string `json:"user_name,omitempty"`
	FirstName *string `json:"first_name,omitempty"`
	LastName  *string `json:"last_name,omitempty"`
	IsForum   *bool   `json:"is_forum,omitempty"`
}

type ChatFullInfo struct {
	ID                                 int64                 `json:"id"`
	Type                               string                `json:"type"`
	Title                              *string               `json:"title,omitempty,"`
	UserName                           *string               `json:"user_name,omitempty,"`
	FirstName                          *string               `json:"first_name,omitempty,"`
	LastName                           *string               `json:"last_name,omitempty,"`
	IsForum                            *bool                 `json:"is_forum,omitempty,"`
	AccentColorId                      int                   `json:"accent_color_id"`
	MaxReactionCount                   int                   `json:"max_reaction_count"`
	Photo                              *ChatPhoto            `json:"photo,omitempty,"`
	ActiveUsernames                    *[]string             `json:"active_usernames,omitempty,"`
	BirthDate                          *BirthDate            `json:"birth_date,omitempty,"`
	BusinessIntro                      *BusinessIntro        `json:"business_intro,omitempty,"`
	BusinessLocation                   *BusinessLocation     `json:"business_location,omitempty,"`
	BusinessOpeningHours               *BusinessOpeningHours `json:"business_opening_hours,omitempty,"`
	PersonalChat                       *Chat                 `json:"personal_chat,omitempty,"`
	AvailableReactions                 *[]ReactionType       `json:"available_reactions,omitempty,"`
	BackgroundCustomEmojiId            *string               `json:"background_custom_emoji_id,omitempty,"`
	ProfileAccentColorId               *bool                 `json:"profile_accent_color_id,omitempty,"`
	ProfileBackgroundCustomEmojiId     *string               `json:"profile_background_custom_emoji_id,omitempty,"`
	EmojiStatusCustomEmojiId           *string               `json:"emoji_status_custom_emoji_id,omitempty,"`
	EmojiStatusExpirationDate          *bool                 `json:"emoji_status_expiration_date,omitempty,"`
	Bio                                *string               `json:"bio,omitempty,"`
	HasPrivateForwards                 *bool                 `json:"has_private_forwards,omitempty,"`
	HasRestrictedVoiceAndVideoMessages *bool                 `json:"has_restricted_voice_and_video_messages,omitempty,"`
	JoinToSendMessages                 *bool                 `json:"join_to_send_messages,omitempty,"`
	JoinByRequest                      *bool                 `json:"join_by_request,omitempty,"`
	Description                        *string               `json:"description,omitempty,"`
	InviteLink                         *string               `json:"invite_link,omitempty,"`
	PinnedMessage                      *Message              `json:"pinned_message,omitempty,"`
	Permissions                        *ChatPermissions      `json:"permissions,omitempty,"`
	AcceptedGiftTypes                  AcceptedGiftTypes     `json:"accepted_gift_types"`
	CanSendPaidMedia                   *bool                 `json:"can_send_paid_media,omitempty,"`
	SlowModeDelay                      *int                  `json:"slow_mode_delay,omitempty,"`
	UnrestrictBoostCount               *int                  `json:"unrestrict_boost_count,omitempty,"`
	MessageAutoDeleteTime              *int                  `json:"message_auto_delete_time,omitempty,"`
	HasAggressiveAntiSpamEnabled       *string               `json:"has_aggressive_anti_spam_enabled,omitempty,"`
	HasHiddenMembers                   *bool                 `json:"has_hidden_members,omitempty,"`
	HasProtectedCount                  *bool                 `json:"has_protected_count,omitempty,"`
	HasVisibleHistory                  *bool                 `json:"has_visible_history,omitempty,"`
	StickerSetName                     *string               `json:"sticker_set_name,omitempty,"`
	CanSetStickerSet                   *bool                 `json:"can_set_sticker_set,omitempty,"`
	CustomEmojiStickerSetName          *string               `json:"custom_emoji_sticker_set_name,omitempty,"`
	LinkedChatId                       *int64                `json:"linked_chat_id,omitempty,"`
	Location                           *ChatLocation         `json:"location,omitempty,"`
}

type Message struct {
	MessageId                     int                            `json:"message_id"`
	MessageThreadId               *int                           `json:"message_thread_id,omitempty"`
	From                          *User                          `json:"from,omitempty"`
	SenderChat                    *Chat                          `json:"sender_chat,omitempty"`
	SenderBoostCount              *int                           `json:"sender_boost_count,omitempty"`
	SenderBusinessBot             *User                          `json:"sender_business_bot,omitempty"`
	Date                          int                            `json:"date"`
	BusinessConnectionId          *string                        `json:"business_connection_id,omitempty"`
	Chat                          Chat                           `json:"chat"`
	ForwardOrigin                 *MessageOrigin                 `json:"forward_origin,omitempty"`
	IsTopicMessage                *bool                          `json:"is_topic_message,omitempty"`
	IsAutomaticForward            *bool                          `json:"is_automatic_forward,omitempty"`
	ReplyToMessage                *Message                       `json:"reply_to_message,omitempty"`
	ExternalReply                 *ExternalReplyInfo             `json:"external_reply,omitempty"`
	Quote                         *TextQuote                     `json:"quote,omitempty"`
	ReplyToStory                  *Story                         `json:"reply_to_story,omitempty"`
	ViaBot                        *User                          `json:"via_bot,omitempty"`
	EditDate                      *int                           `json:"edit_date,omitempty"`
	HasProtectedContent           *bool                          `json:"has_protected_content,omitempty"`
	IsFromOffline                 *bool                          `json:"is_from_offline,omitempty"`
	MediaGroupId                  *string                        `json:"media_group_id,omitempty"`
	AuthorSignature               *string                        `json:"author_signature,omitempty"`
	PaidStarCount                 *int                           `json:"paid_star_count,omitempty"`
	Text                          *string                        `json:"text,omitempty"`
	Entities                      *[]MessageEntity               `json:"entities,omitempty"`
	LinkPreviewOptions            *LinkPreviewOptions            `json:"link_preview_options,omitempty"`
	EffectId                      *string                        `json:"effect_id,omitempty"`
	Animation                     *Animation                     `json:"animation,omitempty"`
	Audio                         *Audio                         `json:"audio,omitempty"`
	Document                      *Document                      `json:"document,omitempty"`
	PaidMedia                     *PaidMediaInfo                 `json:"paid_media,omitempty"`
	Photo                         *[]PhotoSize                   `json:"photo,omitempty"`
	Sticker                       *Sticker                       `json:"sticker,omitempty"`
	Story                         *Story                         `json:"story,omitempty"`
	Video                         *Video                         `json:"video,omitempty"`
	VideoNote                     *VideoNote                     `json:"video_note,omitempty"`
	Voice                         *Voice                         `json:"voice,omitempty"`
	Caption                       *string                        `json:"caption,omitempty"`
	CaptionEntities               *[]MessageEntity               `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia         *bool                          `json:"show_caption_above_media,omitempty"`
	HasMediaSpoiler               *bool                          `json:"has_media_spoiler,omitempty"`
	Contact                       *Contact                       `json:"contact,omitempty"`
	Dice                          *Dice                          `json:"dice,omitempty"`
	Game                          *Game                          `json:"game,omitempty"`
	Poll                          *Poll                          `json:"poll,omitempty"`
	Venue                         *Venue                         `json:"venue,omitempty"`
	Location                      *Location                      `json:"location,omitempty"`
	NewChatMembers                *[]User                        `json:"new_chat_members,omitempty"`
	LeftChatMember                *User                          `json:"left_chat_member,omitempty"`
	NewChatTitle                  *string                        `json:"new_chat_title,omitempty"`
	NewChatPhoto                  *[]PhotoSize                   `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto               *bool                          `json:"delete_chat_photo,omitempty"`
	GroupChatCreated              *bool                          `json:"group_chat_created,omitempty"`
	SuperGroupCreated             *bool                          `json:"super_group_created,omitempty"`
	ChannelChatCreated            *bool                          `json:"channel_chat_created"`
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	MigrateToChatId               *int64                         `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatId             *int64                         `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage                 *MaybeInaccessibleMessage      `json:"pinned_message,omitempty"`
	Invoice                       *Invoice                       `json:"invoice,omitempty"`
	SuccessfulPayment             *SuccessfulPayment             `json:"successful_payment,omitempty"`
	RefundedPayment               *RefundedPayment               `json:"refunded_payment,omitempty"`
	UsersShared                   *UsersShared                   `json:"users_shared,omitempty"`
	ChatShared                    *ChatShared                    `json:"chat_shared,omitempty"`
	Gift                          *GiftInfo                      `json:"gift,omitempty"`
	UniqueGift                    *UniqueGiftInfo                `json:"unique_gift,omitempty"`
	ConnectedWebsite              *string                        `json:"connected_website,omitempty"`
	WriteAccessAllowed            *WriteAccessAllowed            `json:"write_access_allowed,omitempty"`
	PassportData                  *PassportData                  `json:"passport_data,omitempty"`
	ProximityAlertTriggered       *ProximityAlertTriggered       `json:"proximity_alert_triggered,omitempty"`
	BoostAdded                    *ChatBoostAdded                `json:"boost_added,omitempty"`
	ChatBackgroundSet             *ChatBackground                `json:"chat_background_set,omitempty"`
	ForumTopicCreated             *ForumTopicCreated             `json:"forum_topic_created,omitempty"`
	ForumTopicEdited              *ForumTopicEdited              `json:"forum_topic_edited,omitempty"`
	ForumTopicClosed              *ForumTopicClosed              `json:"forum_topic_closed,omitempty"`
	ForumTopicReopened            *ForumTopicReopened            `json:"forum_topic_reopened,omitempty"`
	GeneralForumTopicHidden       *GeneralForumTopicHidden       `json:"general_forum_topic_hidden,omitempty"`
	GeneralForumTopicUnhidden     *GeneralForumTopicUnhidden     `json:"general_forum_topic_unhidden,omitempty"`
	GiveawayCreated               *GiveawayCreated               `json:"giveaway_created,omitempty"`
	Giveaway                      *Giveaway                      `json:"giveaway,omitempty"`
	GiveawayWinners               *GiveawayWinners               `json:"giveaway_winners,omitempty"`
	GiveawayCompleted             *GiveawayCompleted             `json:"giveaway_completed,omitempty"`
	PaidMessagePriceChanged       *PaidMessagePriceChanged       `json:"paid_message_price_changed,omitempty"`
	VideoChatScheduled            *VideoChatScheduled            `json:"video_chat_scheduled,omitempty"`
	VideoChatStarted              *VideoChatStarted              `json:"video_chat_started,omitempty"`
	VideoChatEnded                *VideoChatEnded                `json:"video_chat_ended,omitempty"`
	VideoChatParticipantsInvited  *VideoChatParticipantsInvited  `json:"video_chat_participants_invited,omitempty"`
	WebAppData                    *WebAppData                    `json:"web_app_data,omitempty"`
	ReplyMarkup                   *InlineKeyboardMarkup          `json:"reply_markup,omitempty"`
}

func (m Message) IsCommand() bool {
	if m.Entities == nil {
		return false
	}
	if m.Text == nil {
		return false
	}

	for _, ent := range *m.Entities {
		if ent.Type == "bot_command" {
			return true
		}
	}
	return false
}

func (m *Message) GetCommand() (string, error) {
	if m.Entities == nil {
		return "", fmt.Errorf("the message has no entities")
	}
	if m.Text == nil {
		return "", fmt.Errorf("the message text is empty")
	}

	for _, ent := range *m.Entities {
		if ent.Type == "bot_command" {
			runes := []rune(*m.Text)
			if ent.Offset+ent.Length > len(runes) {
				return "", fmt.Errorf("invalid command entity: offset+length out of bounds")
			}
			return string(runes[ent.Offset : ent.Offset+ent.Length]), nil
		}
	}
	return "", fmt.Errorf("the message has no commands")
}

type MessageId struct {
	MessageId int `json:"message_id"`
}

type InaccessibleMessage struct {
	Chat      Chat `json:"chat"`
	MessageId int  `json:"message_id"`
	Date      int  `json:"date"`
}

type MaybeInaccessibleMessage struct {
	MessageId                     int                            `json:"message_id"`
	MessageThreadId               *int                           `json:"message_thread_id,omitempty"`
	From                          *User                          `json:"from,omitempty"`
	SenderChat                    *Chat                          `json:"sender_chat,omitempty"`
	SenderBoostCount              *int                           `json:"sender_boost_count,omitempty"`
	SenderBusinessBot             *User                          `json:"sender_business_bot,omitempty"`
	Date                          int                            `json:"date"`
	BusinessConnectionId          *string                        `json:"business_connection_id,omitempty"`
	Chat                          Chat                           `json:"chat"`
	ForwardOrigin                 *MessageOrigin                 `json:"forward_origin,omitempty"`
	IsTopicMessage                *bool                          `json:"is_topic_message,omitempty"`
	IsAutomaticForward            *bool                          `json:"is_automatic_forward,omitempty"`
	ReplyToMessage                *Message                       `json:"reply_to_message,omitempty"`
	ExternalReply                 *ExternalReplyInfo             `json:"external_reply,omitempty"`
	Quote                         *TextQuote                     `json:"quote,omitempty"`
	ReplyToStory                  *Story                         `json:"reply_to_story,omitempty"`
	ViaBot                        *User                          `json:"via_bot,omitempty"`
	EditDate                      *int                           `json:"edit_date,omitempty"`
	HasProtectedContent           *bool                          `json:"has_protected_content,omitempty"`
	IsFromOffline                 *bool                          `json:"is_from_offline,omitempty"`
	MediaGroupId                  *string                        `json:"media_group_id,omitempty"`
	AuthorSignature               *string                        `json:"author_signature,omitempty"`
	PaidStarCount                 *int                           `json:"paid_star_count,omitempty"`
	Text                          *string                        `json:"text,omitempty"`
	Entities                      *[]MessageEntity               `json:"entities,omitempty"`
	LinkPreviewOptions            *LinkPreviewOptions            `json:"link_preview_options,omitempty"`
	EffectId                      *string                        `json:"effect_id,omitempty"`
	Animation                     *Animation                     `json:"animation,omitempty"`
	Audio                         *Audio                         `json:"audio,omitempty"`
	Document                      *Document                      `json:"document,omitempty"`
	PaidMedia                     *PaidMediaInfo                 `json:"paid_media,omitempty"`
	Photo                         *[]PhotoSize                   `json:"photo,omitempty"`
	Sticker                       *Sticker                       `json:"sticker,omitempty"`
	Story                         *Story                         `json:"story,omitempty"`
	Video                         *Video                         `json:"video,omitempty"`
	VideoNote                     *VideoNote                     `json:"video_note,omitempty"`
	Voice                         *Voice                         `json:"voice,omitempty"`
	Caption                       *string                        `json:"caption,omitempty"`
	CaptionEntities               *[]MessageEntity               `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia         *bool                          `json:"show_caption_above_media,omitempty"`
	HasMediaSpoiler               *bool                          `json:"has_media_spoiler,omitempty"`
	Contact                       *Contact                       `json:"contact,omitempty"`
	Dice                          *Dice                          `json:"dice,omitempty"`
	Game                          *Game                          `json:"game,omitempty"`
	Poll                          *Poll                          `json:"poll,omitempty"`
	Venue                         *Venue                         `json:"venue,omitempty"`
	Location                      *Location                      `json:"location,omitempty"`
	NewChatMembers                *[]User                        `json:"new_chat_members,omitempty"`
	LeftChatMember                *User                          `json:"left_chat_member,omitempty"`
	NewChatTitle                  *string                        `json:"new_chat_title,omitempty"`
	NewChatPhoto                  *[]PhotoSize                   `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto               *bool                          `json:"delete_chat_photo,omitempty"`
	GroupChatCreated              *bool                          `json:"group_chat_created,omitempty"`
	SuperGroupCreated             *bool                          `json:"super_group_created,omitempty"`
	ChannelChatCreated            *bool                          `json:"channel_chat_created"`
	MessageAutoDeleteTimerChanged *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
	MigrateToChatId               *int64                         `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatId             *int64                         `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage                 *MaybeInaccessibleMessage      `json:"pinned_message,omitempty"`
	Invoice                       *Invoice                       `json:"invoice,omitempty"`
	SuccessfulPayment             *SuccessfulPayment             `json:"successful_payment,omitempty"`
	RefundedPayment               *RefundedPayment               `json:"refunded_payment,omitempty"`
	UsersShared                   *UsersShared                   `json:"users_shared,omitempty"`
	ChatShared                    *ChatShared                    `json:"chat_shared,omitempty"`
	Gift                          *GiftInfo                      `json:"gift,omitempty"`
	UniqueGift                    *UniqueGiftInfo                `json:"unique_gift,omitempty"`
	ConnectedWebsite              *string                        `json:"connected_website,omitempty"`
	WriteAccessAllowed            *WriteAccessAllowed            `json:"write_access_allowed,omitempty"`
	PassportData                  *PassportData                  `json:"passport_data,omitempty"`
	ProximityAlertTriggered       *ProximityAlertTriggered       `json:"proximity_alert_triggered,omitempty"`
	BoostAdded                    *ChatBoostAdded                `json:"boost_added,omitempty"`
	ChatBackgroundSet             *ChatBackground                `json:"chat_background_set,omitempty"`
	ForumTopicCreated             *ForumTopicCreated             `json:"forum_topic_created,omitempty"`
	ForumTopicEdited              *ForumTopicEdited              `json:"forum_topic_edited,omitempty"`
	ForumTopicClosed              *ForumTopicClosed              `json:"forum_topic_closed,omitempty"`
	ForumTopicReopened            *ForumTopicReopened            `json:"forum_topic_reopened,omitempty"`
	GeneralForumTopicHidden       *GeneralForumTopicHidden       `json:"general_forum_topic_hidden,omitempty"`
	GeneralForumTopicUnhidden     *GeneralForumTopicUnhidden     `json:"general_forum_topic_unhidden,omitempty"`
	GiveawayCreated               *GiveawayCreated               `json:"giveaway_created,omitempty"`
	Giveaway                      *Giveaway                      `json:"giveaway,omitempty"`
	GiveawayWinners               *GiveawayWinners               `json:"giveaway_winners,omitempty"`
	GiveawayCompleted             *GiveawayCompleted             `json:"giveaway_completed,omitempty"`
	PaidMessagePriceChanged       *PaidMessagePriceChanged       `json:"paid_message_price_changed,omitempty"`
	VideoChatScheduled            *VideoChatScheduled            `json:"video_chat_scheduled,omitempty"`
	VideoChatStarted              *VideoChatStarted              `json:"video_chat_started,omitempty"`
	VideoChatEnded                *VideoChatEnded                `json:"video_chat_ended,omitempty"`
	VideoChatParticipantsInvited  *VideoChatParticipantsInvited  `json:"video_chat_participants_invited,omitempty"`
	WebAppData                    *WebAppData                    `json:"web_app_data,omitempty"`
	ReplyMarkup                   *InlineKeyboardMarkup          `json:"reply_markup,omitempty"`
}

func (m MaybeInaccessibleMessage) IsAccessible() bool {
	return m.Date != 0
}

type MessageEntity struct {
	Type          string  `json:"type"`
	Offset        int     `json:"offset"`
	Length        int     `json:"length"`
	Url           *string `json:"url,omitempty"`
	User          *User   `json:"user,omitempty"`
	Language      *string `json:"language,omitempty"`
	CustomEmojiId *string `json:"custom_emoji_id,omitempty"`
}

type TextQuote struct {
	Text     string           `json:"text"`
	Entities *[]MessageEntity `json:"entities,omitempty"`
	Position int              `json:"position"`
	IsManual *bool            `json:"is_manual,omitempty"`
}

type ExternalReplyInfo struct {
	Origin             MessageOrigin       `json:"origin"`
	Chat               *Chat               `json:"chat,omitempty"`
	MessageId          *int                `json:"message_id,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	Animation          *Animation          `json:"animation,omitempty"`
	Audio              *Audio              `json:"audio,omitempty"`
	Document           *Document           `json:"document,omitempty"`
	PaidMedia          *PaidMediaInfo      `json:"paid_media,omitempty"`
	Photo              *[]PhotoSize        `json:"photo,omitempty"`
	Sticker            *Sticker            `json:"sticker,omitempty"`
	Story              *Story              `json:"story,omitempty"`
	Video              *Video              `json:"video,omitempty"`
	VideoNote          *VideoNote          `json:"video_note,omitempty"`
	Voice              *Voice              `json:"voice,omitempty"`
	HasMediaSpoiler    *bool               `json:"has_media_spoiler,omitempty"`
	Contact            *Contact            `json:"contact,omitempty"`
	Dice               *Dice               `json:"dice,omitempty"`
	Game               *Game               `json:"game,omitempty"`
	Giveaway           *Giveaway           `json:"giveaway,omitempty"`
	GiveawayWinners    *GiveawayWinners    `json:"giveaway_winners,omitempty"`
	Invoice            *Invoice            `json:"invoice,omitempty"`
	Location           *Location           `json:"location"`
	Poll               *Poll               `json:"poll,omitempty"`
	Venue              *Venue              `json:"venue,omitempty"`
}

type ReplyParameters struct {
	MessageId                int              `json:"message_id"`
	ChatId                   *string          `json:"chat_id,omitempty"`
	AllowSendingWithoutReply *bool            `json:"allow_sending_without_reply,omitempty"`
	Quote                    *string          `json:"quote,omitempty"`
	QuoteParseMode           *string          `json:"quote_parse_mode,omitempty"`
	QuoteEntities            *[]MessageEntity `json:"quote_entities,omitempty"`
	QuotePosition            *int             `json:"quote_position,omitempty"`
}

type MessageOrigin struct {
	Type            string
	Date            int     `json:"date"`
	SenderUser      User    `json:"sender_user"`
	SenderUsername  string  `json:"sender_username"`
	SenderChat      Chat    `json:"sender_chat"`
	MessageId       int     `json:"message_id"`
	AuthorSignature *string `json:"author_signature,omitzero"`
}

type PhotoSize struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	FileSize     *int   `json:"file_size,omitempty"`
}

type Animation struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     *string    `json:"file_name,omitempty"`
	MimeType     *string    `json:"mime_type,omitempty"`
	FileSize     *int64     `json:"file_size,omitempty"`
}

type Audio struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Duration     int        `json:"duration"`
	Performer    *string    `json:"performer,omitempty"`
	Title        *string    `json:"title,omitempty"`
	FileName     *string    `json:"file_name,omitempty"`
	MimeType     *string    `json:"mime_type,omitempty"`
	FileSize     *int64     `json:"file_size,omitempty"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
}

type Document struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileName     *string    `json:"file_name,omitempty"`
	MimeType     *string    `json:"mime_type,omitempty"`
	FileSize     *int64     `json:"file_size,omitempty"`
}

type Story struct {
	Chat Chat `json:"chat"`
	Id   int  `json:"id"`
}

type Video struct {
	FileId         string       `json:"file_id"`
	FileUniqueId   string       `json:"file_unique_id"`
	Width          int          `json:"width"`
	Height         int          `json:"height"`
	Duration       int          `json:"duration"`
	Thumbnail      *PhotoSize   `json:"thumbnail,omitempty,"`
	Cover          *[]PhotoSize `json:"cover,omitempty"`
	StartTimestamp *int         `json:"start_timestamp,omitempty"`
	FileName       *string      `json:"file_name,omitempty,"`
	MimeType       *string      `json:"mime_type,omitempty,"`
	FileSize       *int64       `json:"file_size,omitempty,"`
}

type VideoNote struct {
	FileId       string     `json:"file_id"`
	FileUniqueId string     `json:"file_unique_id"`
	Length       int        `json:"length"`
	Duration     int        `json:"duration"`
	Thumbnail    *PhotoSize `json:"thumbnail,omitempty"`
	FileSize     *int       `json:"file_size,omitempty"`
}

type Voice struct {
	FileId       string  `json:"file_id"`
	FileUniqueId string  `json:"file_unique_id"`
	Duration     int     `json:"duration"`
	MimeType     *string `json:"mime_type,omitempty"`
	FileSize     *int    `json:"file_size,omitempty"`
}

type PaidMediaInfo struct {
	StarCount string      `json:"star_count"`
	PaidMedia []PaidMedia `json:"paid_media"`
}

type PaidMedia struct {
	Type    string `json:"type"`
	Preview *PaidMediaPreview
	Photo   *PaidMediaPhoto
	Video   *PaidMediaVideo
}

type PaidMediaPreview struct {
	Type     string `json:"type"`
	Width    *int   `json:"width,omitempty"`
	Height   *int   `json:"height,omitempty"`
	Duration *int   `json:"duration,omitempty"`
}

type PaidMediaPhoto struct {
	Type  string      `json:"type"`
	Photo []PhotoSize `json:"photo"`
}

type PaidMediaVideo struct {
	Type  string `json:"type"`
	Video Video  `json:"video"`
}

type Contact struct {
	PhoneNumber string  `json:"phone_number"`
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name,omitempty"`
	UserId      *int64  `json:"user_id,omitempty"`
	VCard       *string `json:"v_card,omitempty"`
}

type Dice struct {
	Value int    `json:"value"`
	Emoji string `json:"emoji"`
}

type PollOption struct {
	Text         string           `json:"text"`
	VoterCount   int              `json:"voter_count"`
	TextEntities *[]MessageEntity `json:"text_entities,omitempty"`
}

type InputPollOption struct {
	Text          string           `json:"text"`
	TextParseMode *string          `json:"text_parse_mode,omitempty"`
	TextEntities  *[]MessageEntity `json:"text_entities,omitempty"`
}

type PollAnswer struct {
	PollId    string `json:"poll_id"`
	VoterChat *Chat  `json:"voter_chat"`
	User      *User  `json:"user"`
	OptionIds []int  `json:"option_ids"`
}

type Poll struct {
	Id                    string           `json:"id"`
	Question              string           `json:"question"`
	QuestionEntities      []MessageEntity  `json:"question_entities"`
	Options               []PollOption     `json:"options"`
	TotalVoterCount       int              `json:"total_voter_count"`
	IsClosed              bool             `json:"is_closed"`
	IsAnonymous           bool             `json:"is_anonymous"`
	Type                  string           `json:"type"`
	AllowsMultipleAnswers bool             `json:"allows_multiple_answers"`
	CorrectOptionId       *int             `json:"correct_option_id,omitempty"`
	Explanation           *string          `json:"explanation,omitempty"`
	ExplanationEntities   *[]MessageEntity `json:"explanation_entities,omitempty"`
	OpenPeriod            *int             `json:"open_period,omitempty"`
	CloseDate             *int             `json:"close_date,omitempty"`
}

type Location struct {
	Latitude             float64  `json:"latitude"`
	Longitude            float64  `json:"longitude"`
	HorizontalAccuracy   *float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           *int     `json:"live_period,omitempty"`
	Heading              *int     `json:"heading,omitempty"`
	ProximityAlertRadius *int     `json:"proximity_alert_radius,omitempty"`
}

type Venue struct {
	Location        Location `json:"location"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareId    *string  `json:"foursquare_id,omitempty"`
	FourSquareType  *string  `json:"four_square_type,omitempty"`
	GooglePlaceId   *string  `json:"google_place_id,omitempty"`
	GooglePlaceType *string  `json:"google_place_type,omitempty"`
}

type WebAppData struct {
	Data       string `json:"data"`
	ButtonText string `json:"button_text"`
}

type ProximityAlertTriggered struct {
	Traveler User `json:"traveler"`
	Watcher  User `json:"watcher"`
	Distance int  `json:"distance"`
}

type MessageAutoDeleteTimerChanged struct {
	MessageAutoDeleteTime int `json:"message_auto_delete_time"`
}

type ChatBoostAdded struct {
	BoostCount int `json:"boost_count"`
}

type BackgroundFill struct {
	Type          string `json:"type"`
	Color         int    `json:"color"`
	TopColor      int    `json:"top_color"`
	BottomColor   int    `json:"bottom_color"`
	RotationAngle int    `json:"rotation_angle"`
	Colors        []int  `json:"colors"`
}

type BackgroundType struct {
	Type             string         `json:"type"`
	Fill             BackgroundFill `json:"fill"`
	DarkThemeDimming int            `json:"dark_theme_dimming"`
	Document         Document       `json:"document"`
	IsBlurred        *bool          `json:"is_blurred,omitempty"`
	IsMoving         *bool          `json:"is_moving,omitempty"`
	Intensity        int            `json:"intensity"`
	IsInverted       *bool          `json:"is_inverted,omitempty"`
	ThemeName        string         `json:"theme_name"`
}

type ChatBackground struct {
	Type BackgroundType `json:"type"`
}

type ForumTopicCreated struct {
	Name              string  `json:"name"`
	IconColor         int     `json:"icon_color"`
	IconCustomEmojiId *string `json:"icon_custom_emoji_id,omitempty"`
}

type ForumTopicClosed struct{}

type ForumTopicEdited struct {
	Name              *string `json:"name,omitempty"`
	IconCustomEmojiId *string `json:"icon_custom_emoji_id,omitempty"`
}

type ForumTopicReopened struct{}

type GeneralForumTopicHidden struct{}

type GeneralForumTopicUnhidden struct{}

type SharedUser struct {
	UserId    int64        `json:"user_id"`
	FirstName *string      `json:"first_name,omitempty"`
	LastName  *string      `json:"last_name,omitempty"`
	Username  *string      `json:"username,omitempty"`
	Photo     *[]PhotoSize `json:"photo,omitempty"`
}

type UsersShared struct {
	RequestId string       `json:"request_id"`
	Users     []SharedUser `json:"users"`
}

type ChatShared struct {
	RequestId string       `json:"request_id"`
	ChatId    int64        `json:"chat_id"`
	Title     *string      `json:"title,omitempty"`
	Username  *string      `json:"username,omitempty"`
	Photo     *[]PhotoSize `json:"photo,omitempty"`
}

type WriteAccessAllowed struct {
	FromRequest        *bool   `json:"from_request,omitempty"`
	WebAppName         *string `json:"web_app_name,omitempty"`
	FromAttachmentMenu *bool   `json:"from_attachment_menu,omitempty"`
}

type VideoChatScheduled struct {
	StartDate int `json:"start_date"`
}

type VideoChatStarted struct{}

type VideoChatEnded struct {
	Duration int `json:"duration"`
}

type VideoChatParticipantsInvited struct {
	Users []User `json:"users"`
}

type PaidMessagePriceChanged struct {
	PaidMessageStarCount int `json:"paid_message_star_count"`
}

type GiveawayCreated struct {
	PrizeStarCount *int `json:"prize_star_count,omitempty"`
}

type Giveaway struct {
	Chats                         []Chat    `json:"chats"`
	WinnerSelectionDate           int       `json:"winner_selection_date"`
	WinnerCount                   int       `json:"winner_count"`
	OnlyNewMembers                *bool     `json:"only_new_members,omitempty"`
	HasPublicWinners              *bool     `json:"has_public_winners,omitempty"`
	PrizeDescription              *string   `json:"prize_description,omitempty"`
	CountryCodes                  *[]string `json:"country_codes,omitempty"`
	PrizeStarCount                *int      `json:"prize_star_count,omitempty"`
	PremiumSubscriptionMonthCount *int      `json:"premium_subscription_month_count,omitempty"`
}

type GiveawayWinners struct {
	Chat                          Chat    `json:"chat"`
	GiveawayMessageId             int     `json:"giveaway_message_id"`
	WinnersSelectionDate          int     `json:"winners_selection_date"`
	WinnerCount                   int     `json:"winner_count"`
	Winners                       []User  `json:"winners"`
	AdditionalChatCount           *int    `json:"additional_chat_count,omitempty"`
	PrizeStarCount                *int    `json:"prize_star_count,omitempty"`
	PremiumSubscriptionMonthCount *int    `json:"premium_subscription_month_count,omitempty"`
	UnclaimedPrizeCount           *int    `json:"unclaimed_prize_count,omitempty"`
	OnlyNewMembers                *bool   `json:"only_new_members,omitempty"`
	WasRefunded                   *bool   `json:"was_refunded,omitempty"`
	PrizeDescription              *string `json:"prize_description,omitempty"`
}

type GiveawayCompleted struct {
	WinnerCount         int      `json:"winner_count"`
	UnclaimedPrizeCount *int     `json:"unclaimed_prize_count,omitempty"`
	GiveawayMessage     *Message `json:"giveaway_message,omitempty"`
	IsStarGiveaway      *bool    `json:"is_star_giveaway,omitempty"`
}

type LinkPreviewOptions struct {
	IsDisabled       *bool   `json:"is_disabled,omitempty"`
	UrlFileId        *string `json:"url_file_id,omitempty"`
	PreferSmallMedia *bool   `json:"prefer_small_media,omitempty"`
	PreferLargeMedia *bool   `json:"prefer_large_media,omitempty"`
	ShowAboveText    *bool   `json:"show_above_text,omitempty"`
}

type UserProfilePhotos struct {
	TotalCount int           `json:"total_count"`
	Photos     [][]PhotoSize `json:"photos"`
}

type File struct {
	FileId       string  `json:"file_id"`
	FileUniqueId string  `json:"file_unique_id"`
	FileSize     *int64  `json:"file_size,omitempty"`
	FilePath     *string `json:"file_path,omitempty"`
}

type WebAppInfo struct {
	Url string `json:"url"`
}

// ReplyMarkup is basically a union type for
// [InlineKeyboardMarkup], [ReplyKeyboardMarkup], [ReplyKeyboardRemove] and [ForceReply]
type ReplyMarkup interface {
	replyKeyboardContract()
}

type ReplyKeyboardMarkup struct {
	Keyboard              [][]KeyboardButton `json:"keyboard"`
	IsPersistent          *bool              `json:"is_persistent,omitempty"`
	ResizeKeyboard        *bool              `json:"resize_keyboard,omitempty"`
	OneTimeKeyboard       *bool              `json:"one_time_keyboard,omitempty"`
	InputFieldPlaceholder *string            `json:"input_field_placeholder,omitempty"`
	Selective             *bool              `json:"selective,omitempty"`
}

func (f ReplyKeyboardMarkup) replyKeyboardContract() {}

type KeyboardButton struct {
	Text            string                      `json:"text"`
	RequestUsers    *KeyboardButtonRequestUsers `json:"request_users,omitempty"`
	RequestChat     *KeyboardButtonRequestChat  `json:"request_chat,omitempty"`
	RequestContact  *bool                       `json:"request_contact,omitempty"`
	RequestLocation *bool                       `json:"request_location,omitempty"`
	RequestPoll     *KeyboardButtonPollType     `json:"request_poll,omitempty"`
	WebApp          *WebAppInfo                 `json:"web_app,omitempty"`
}

type KeyboardButtonRequestUsers struct {
	RequestId       int32 `json:"request_id"`
	UserIsBot       *bool `json:"user_is_bot,omitempty"`
	UserIsPremium   *bool `json:"user_is_premium,omitempty"`
	MaxQuantity     *int  `json:"max_quantity,omitempty"`
	RequestName     *bool `json:"request_name,omitempty"`
	RequestUsername *bool `json:"request_username,omitempty"`
	RequestPhoto    *bool `json:"request_photo,omitempty"`
}

type KeyboardButtonRequestChat struct {
	RequestId               int32                    `json:"request_id"`
	ChatIsChannel           bool                     `json:"chat_is_channel"`
	ChatIsForum             *bool                    `json:"chat_is_forum,omitempty"`
	ChatHasUsername         *bool                    `json:"chat_has_username,omitempty"`
	ChatIsCreated           *bool                    `json:"chat_is_created,omitempty"`
	UserAdministratorRights *ChatAdministratorRights `json:"user_administrator_rights,omitempty"`
	BotAdministratorRights  *ChatAdministratorRights `json:"bot_administrator_rights,omitempty"`
	BotIsMember             *bool                    `json:"bot_is_member,omitempty"`
	RequestTitle            *bool                    `json:"request_title,omitempty"`
	RequestUsername         *bool                    `json:"request_username,omitempty"`
	RequestPhoto            *bool                    `json:"request_photo,omitempty"`
}

type KeyboardButtonPollType struct {
	Type *string `json:"type,omitempty"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool  `json:"remove_keyboard"`
	Selective      *bool `json:"selective,omitempty"`
}

func (f ReplyKeyboardRemove) replyKeyboardContract() {}

type InlineKeyboardMarkup struct {
	Keyboard [][]InlineKeyboardButton `json:"keyboard"`
}

func (f InlineKeyboardMarkup) replyKeyboardContract() {}

type InlineKeyboardButton struct {
	Text                         string                       `json:"text"`
	Url                          *string                      `json:"url,omitempty"`
	CallbackData                 *string                      `json:"callback_data,omitempty"`
	WebApp                       *WebAppInfo                  `json:"web_app,omitempty"`
	LoginUrl                     *LoginUrl                    `json:"login_url,omitempty"`
	SwitchInlineQuery            *string                      `json:"switch_inline_query,omitempty"`
	SwitchInlineQueryCurrentChat *string                      `json:"switch_inline_query_current_chat,omitempty"`
	SwitchInlineQueryChosenChat  *SwitchInlineQueryChosenChat `json:"switch_inline_query_chosen_chat,omitempty"`
	CopyText                     *CopyTextButton              `json:"copy_text,omitempty"`
	CallbackGame                 *CallbackGame                `json:"callback_game,omitempty"`
	Pay                          *bool                        `json:"pay,omitempty"`
}

type LoginUrl struct {
	Url                string  `json:"url"`
	ForwardText        *string `json:"forward_text,omitempty"`
	BotUsername        *string `json:"bot_username,omitempty"`
	RequestWriteAccess *bool   `json:"request_write_access,omitempty"`
}

type SwitchInlineQueryChosenChat struct {
	Query             *string `json:"query,omitempty"`
	AllowUserChats    *bool   `json:"allow_user_chats,omitempty"`
	AllowBotChats     *bool   `json:"allow_bot_chats,omitempty"`
	AllowGroupChats   *bool   `json:"allow_group_chats,omitempty"`
	AllowChannelChats *bool   `json:"allow_channel_chats,omitempty"`
}

type CopyTextButton struct {
	Text string
}

type CallbackQuery struct {
	Id              string                    `json:"id"`
	From            User                      `json:"from"`
	Message         *MaybeInaccessibleMessage `json:"message,omitempty"`
	InlineMessageId *string                   `json:"inline_message_id,omitempty"`
	ChatInstance    string                    `json:"chat_instance"`
	Data            *string                   `json:"data,omitempty"`
	GameShortName   *string                   `json:"game_short_name,omitempty"`
}

type ForceReply struct {
	ForceReply            bool    `json:"force_reply"`
	InputFieldPlaceholder *string `json:"input_field_placeholder,omitempty"`
	Selective             *bool   `json:"selective,omitempty"`
}

func (f ForceReply) replyKeyboardContract() {}

type ChatPhoto struct {
	SmallFileId       string `json:"small_file_id"`
	SmallFileUniqueId string `json:"small_file_unique_id"`
	BigFileId         string `json:"big_file_id"`
	BigFileUniqueId   string `json:"big_file_unique_id"`
}

type ChatInviteLink struct {
	InviteLink              string  `json:"invite_link"`
	Creator                 User    `json:"creator"`
	CreatesJoinRequest      bool    `json:"creates_join_request"`
	IsPrimary               bool    `json:"is_primary"`
	IsRevoked               bool    `json:"is_revoked"`
	Name                    *string `json:"name,omitempty"`
	ExpireDate              *int    `json:"expire_date,omitempty"`
	MemberLimit             *bool   `json:"member_limit,omitempty"`
	PendingJoinRequestCount *int    `json:"pending_join_request_count,omitempty"`
	SubscriptionPeriod      *int    `json:"subscription_period,omitempty"`
	SubscriptionPrice       *int    `json:"subscription_price,omitempty"`
}

type ChatAdministratorRights struct {
	IsAnonymous         bool  `json:"is_anonymous"`
	CanManageChat       bool  `json:"can_manage_chat"`
	CanDeleteMessages   bool  `json:"can_delete_messages"`
	CanManageVideoChats bool  `json:"can_manage_video_chats"`
	CanRestrictMembers  bool  `json:"can_restrict_members"`
	CanPromoteMembers   bool  `json:"can_promote_members"`
	CanChangeInfo       bool  `json:"can_change_info"`
	CanInviteUsers      bool  `json:"can_invite_users"`
	CanPostStories      bool  `json:"can_post_stories"`
	CanEditStories      bool  `json:"can_edit_stories"`
	CanDeleteStories    bool  `json:"can_delete_stories"`
	CanPostMessages     *bool `json:"can_post_messages,omitempty"`
	CanEditMessages     *bool `json:"can_edit_messages,omitempty"`
	CanPinMessages      *bool `json:"can_pin_messages,omitempty"`
	CanManageTopics     *bool `json:"can_manage_topics,omitempty"`
}

type ChatMemberUpdated struct {
	Chat                    Chat            `json:"chat"`
	From                    User            `json:"from"`
	Date                    int             `json:"date"`
	OldChatMember           ChatMember      `json:"old_chat_member"`
	NewChatMember           ChatMember      `json:"new_chat_member"`
	InviteLink              *ChatInviteLink `json:"invite_link,omitempty"`
	ViaJoinRequest          *bool           `json:"via_join_request,omitempty"`
	ViaChatFolderInviteLink *bool           `json:"via_chat_folder_invite_link,omitempty"`
}

type ChatMember struct {
	Status                string  `json:"status"`
	User                  User    `json:"user"`
	IsAnonymous           bool    `json:"is_anonymous"`
	CanBeEdited           bool    `json:"can_be_edited"`
	CanManageChat         bool    `json:"can_manage_chat"`
	CanDeleteMessages     bool    `json:"can_delete_messages"`
	CanManageVideoChats   bool    `json:"can_manage_video_chats"`
	CanRestrictMembers    bool    `json:"can_restrict_members"`
	CanPromoteMembers     bool    `json:"can_promote_members"`
	CanChangeInfo         bool    `json:"can_change_info"`
	CanInviteUsers        bool    `json:"can_invite_users"`
	CanPostStories        bool    `json:"can_post_stories"`
	CanEditStories        bool    `json:"can_edit_stories"`
	CanDeleteStories      bool    `json:"can_delete_stories"`
	IsMember              bool    `json:"is_member"`
	CanSendMessages       bool    `json:"can_send_messages"`
	CanSendAudios         bool    `json:"can_send_audios"`
	CanSendDocuments      bool    `json:"can_send_documents"`
	CanSendPhotos         bool    `json:"can_send_photos"`
	CanSendVideos         bool    `json:"can_send_videos"`
	CanSendVideoNotes     bool    `json:"can_send_video_notes"`
	CanSendVoiceNotes     bool    `json:"can_send_voice_notes"`
	CanSendPolls          bool    `json:"can_send_polls"`
	CanSendOtherMessages  bool    `json:"can_send_other_messages"`
	CanAddWebpagePreviews bool    `json:"can_add_webpage_previews"`
	CanPostMessages       *bool   `json:"can_post_messages,omitempty"`
	CanEditMessages       *bool   `json:"can_edit_messages,omitempty"`
	CanPinMessages        *bool   `json:"can_pin_messages,omitempty"`
	CanManageTopics       *bool   `json:"can_manage_topics,omitempty"`
	UntilDate             *int    `json:"until_date,omitempty"`
	CustomTitle           *string `json:"custom_title,omitempty"`
}

type ChatJoinRequest struct {
	Chat       Chat            `json:"chat"`
	User       User            `json:"user"`
	UserChatId int64           `json:"user_chat_id"`
	Date       int             `json:"date"`
	Bio        *string         `json:"bio,omitempty"`
	InviteLink *ChatInviteLink `json:"invite_link,omitempty"`
}

type ChatPermissions struct {
	CanSendMessages       *bool `json:"can_send_messages,omitempty"`
	CanSendAudios         *bool `json:"can_send_audios,omitempty"`
	CanSendDocuments      *bool `json:"can_send_documents,omitempty"`
	CanSendPhotos         *bool `json:"can_send_photos,omitempty"`
	CanSendVideos         *bool `json:"can_send_videos,omitempty"`
	CanSendVideoNotes     *bool `json:"can_send_video_notes,omitempty"`
	CanSendVoiceNotes     *bool `json:"can_send_voice_notes,omitempty"`
	CanSendPolls          *bool `json:"can_send_polls,omitempty"`
	CanSendOtherMessages  *bool `json:"can_send_other_messages,omitempty"`
	CanAddWebpagePreviews *bool `json:"can_add_webpage_previews,omitempty"`
	CanChangeInfo         *bool `json:"can_change_info,omitempty"`
	CanInviteUsers        *bool `json:"can_invite_users,omitempty"`
	CanPinMessages        *bool `json:"can_pin_messages,omitempty"`
	CanManageTopics       *bool `json:"can_manage_topics,omitempty"`
}

type BirthDate struct {
	Day   int  `json:"day"`
	Month int  `json:"month"`
	Year  *int `json:"year,omitempty"`
}

type BusinessIntro struct {
	Title   string   `json:"title,omitempty"`
	Message *string  `json:"message,omitempty"`
	Sticker *Sticker `json:"sticker,omitempty"`
}

type BusinessLocation struct {
	Address  string    `json:"address"`
	Location *Location `json:"location,omitempty"`
}

type BusinessOpeningHoursInterval struct {
	OpeningMinute int `json:"opening_minute"`
	ClosingMinute int `json:"closing_minute"`
}

type BusinessOpeningHours struct {
	TimeZone     string                         `json:"time_zone"`
	OpeningHours []BusinessOpeningHoursInterval `json:"opening_hours"`
}

type StoryAreaPosition struct {
	XPercentage            float64 `json:"x_percentage"`
	YPercentage            float64 `json:"y_percentage"`
	WidthPercentage        float64 `json:"width_percentage"`
	HeightPercentage       float64 `json:"height_percentage"`
	RotationAngle          float64 `json:"rotation_angle"`
	CornerRadiusPercentage float64 `json:"corner_radius_percentage"`
}

type LocationAddress struct {
	CountryCode string  `json:"country_code"`
	State       *string `json:"state,omitempty"`
	City        *string `json:"city,omitempty"`
	Street      *string `json:"street,omitempty"`
}

type StoryAreaType struct {
	Type            string
	Latitude        float64          `json:"latitude"`
	Longtitude      float64          `json:"longtitude"`
	ReactionType    ReactionType     `json:"reaction_type"`
	Url             string           `json:"url"`
	Temperatue      float64          `json:"temperatue"`
	Emoji           string           `json:"emoji"`
	BackgroundColor int              `json:"background_color"`
	Name            string           `json:"name"`
	Address         *LocationAddress `json:"address,omitempty"`
	IsDark          *bool            `json:"is_dark,omitempty"`
	IsFlipped       *bool            `json:"is_flipped,omitempty"`
}

type StoryArea struct {
	Position StoryAreaPosition `json:"position"`
	Type     StoryAreaType     `json:"type"`
}

type ChatLocation struct {
	Location Location `json:"location"`
	Address  string   `json:"address"`
}

type ReactionType struct {
	Type          string `json:"type"`
	Emoji         string `json:"emoji"`
	CustomEmojiId string `json:"custom_emoji_id"`
}

type ReactionCount struct {
	Type       ReactionType `json:"type"`
	TotalCount int          `json:"total_count"`
}

type MessageReactionUpdated struct {
	Chat        Chat           `json:"chat"`
	MessageId   int            `json:"message_id"`
	User        *User          `json:"user,omitempty"`
	ActorChat   *Chat          `json:"actor_chat,omitempty"`
	Date        int            `json:"date"`
	OldReaction []ReactionType `json:"old_reaction"`
	NewReaction []ReactionType `json:"new_reaction"`
}

type MessageReactionCountUpdated struct {
	Chat      Chat            `json:"chat"`
	MessageId int             `json:"message_id"`
	Date      int             `json:"date"`
	Reactions []ReactionCount `json:"reactions"`
}

type ForumTopic struct {
	MessageThreadId   int    `json:"message_thread_id"`
	Name              string `json:"name"`
	IconColor         int    `json:"icon_color"`
	IconCustomEmojiId string `json:"icon_custom_emoji_id"`
}

type Gift struct {
	Id               string  `json:"id"`
	Sticker          Sticker `json:"sticker"`
	StarCount        int     `json:"star_count"`
	UpgradeStarCount *int    `json:"upgrade_star_count,omitempty"`
	TotalCount       *int    `json:"total_count,omitempty,"`
	RemainingCount   *int    `json:"remaining_count,omitempty,"`
}

type Gifts struct {
	Gifts []Gift `json:"gifts"`
}

type UniqueGiftModel struct {
	Name          string  `json:"name"`
	Sticker       Sticker `json:"sticker"`
	RarityPerMile int     `json:"rarity_per_mile"`
}

type UniqueGiftSymbol struct {
	Name          string  `json:"name"`
	Sticker       Sticker `json:"sticker"`
	RarityPerMile int     `json:"rarity_per_mile"`
}

type UniqueGiftBackdropColors struct {
	CenterColor int `json:"center_color"`
	EdgeColor   int `json:"edge_color"`
	SymbolColor int `json:"symbol_color"`
	TextColor   int `json:"text_color"`
}

type UniqueGiftBackdrop struct {
	Name          string                   `json:"name"`
	Colors        UniqueGiftBackdropColors `json:"colors"`
	RarityPerMile int                      `json:"rarity_per_mile"`
}

type UniqueGift struct {
	BaseName string             `json:"base_name"`
	Name     string             `json:"name"`
	Number   int                `json:"number"`
	Model    UniqueGiftModel    `json:"model"`
	Symbol   UniqueGiftSymbol   `json:"symbol"`
	Backdrop UniqueGiftBackdrop `json:"backdrop"`
}

type GiftInfo struct {
	Gift                    Gift             `json:"gift"`
	OwnedGiftId             *string          `json:"owned_gift_id,omitempty"`
	ConvertStarCount        *int             `json:"convert_star_count,omitempty"`
	PrepaidUpgradeStarCount *int             `json:"prepaid_upgrade_star_count,omitempty"`
	CanBeUpgraded           *bool            `json:"can_be_upgraded,omitempty"`
	Text                    *string          `json:"text,omitempty"`
	Entities                *[]MessageEntity `json:"entities,omitempty"`
	IsPrivate               *bool            `json:"is_private,omitempty"`
}

type UniqueGiftInfo struct {
	Gift              UniqueGift `json:"gift"`
	Origin            string     `json:"origin"`
	OwnedGiftId       *string    `json:"owned_gift_id,omitempty"`
	TransferStarCount *int       `json:"transfer_star_count,omitempty"`
}

type OwnedGift struct {
	Type                    string           `json:"type"`
	Gift                    Gift             `json:"gift"`
	SendDate                int              `json:"send_date"`
	SenderDate              int              `json:"sender_date"`
	OwnedGiftId             *string          `json:"owned_gift_id,omitempty"`
	SenderUser              *User            `json:"sender_user,omitempty"`
	Text                    *string          `json:"text,omitempty"`
	Entities                *[]MessageEntity `json:"entities,omitempty"`
	IsPrivate               *bool            `json:"is_private,omitempty"`
	IsSaved                 *bool            `json:"is_saved,omitempty"`
	CanBeUpgraded           *bool            `json:"can_be_upgraded,omitempty"`
	WasRefunded             *bool            `json:"was_refunded,omitempty"`
	ConvertStarCount        *int             `json:"convert_star_count,omitempty"`
	PrepaidUpgradeStarCount *int             `json:"prepaid_upgrade_star_count,omitempty"`
	CanBeTransferred        *bool            `json:"can_be_transferred,omitempty"`
	TransferStarCount       *int             `json:"transfer_star_count,omitempty"`
}

type OwnedGifts struct {
	TotalCount int         `json:"total_count"`
	Gifts      []OwnedGift `json:"gifts"`
	NextOffset *string     `json:"next_offset.omitempty"`
}

type StarAmount struct {
	Amount         int  `json:"amount"`
	NanostarAmount *int `json:"nanostar_amount,omitempty"`
}

type AcceptedGiftTypes struct {
	UnlimitedGifts      bool `json:"unlimited_gifts"`
	LimitedGifts        bool `json:"limited_gifts"`
	UniqueGifts         bool `json:"unique_gifts"`
	PremiumSubscription bool `json:"premium_subscription"`
}

type BotCommand struct {
	Command     string `json:"command"`
	Description string `json:"description"`
}

type BotCommandScope interface {
	Scope() string
	json.Marshaler
}

type botCommandScopeNoParams string

func (b botCommandScopeNoParams) Scope() string {
	return string(b)
}

func (b botCommandScopeNoParams) MarshalJSON() ([]byte, error) {
	return []byte(`{"type": "` + b.Scope() + `"}`), nil
}

const (
	BotCommandScopeDefault               botCommandScopeNoParams = "default"
	BotCommandScopeAllPrivateChats       botCommandScopeNoParams = "all_private_chats"
	BotCommandScopeAllGroupChats         botCommandScopeNoParams = "all_group_chats"
	BotCommandScopeAllChatAdministrators botCommandScopeNoParams = "all_chat_administrators"
)

type BotCommandScopeChat string

func (b BotCommandScopeChat) Scope() string {
	return "chat"
}

func (b BotCommandScopeChat) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `{"type": "%s", "chat_id": "%s"}`, b.Scope(), string(b)), nil
}

type BotCommandScopeChatAdministrators string

func (b BotCommandScopeChatAdministrators) Scope() string {
	return "chat_administrators"
}

func (b BotCommandScopeChatAdministrators) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `{"type": "%s", "chat_id": "%s"}`, b.Scope(), string(b)), nil
}

type BotCommandScopeChatMember struct {
	ChatID string `json:"chat_id"`
	UserID int    `json:"user_id"`
}

func (b BotCommandScopeChatMember) Scope() string {
	return "chat_member"
}

func (b BotCommandScopeChatMember) MarshalJSON() ([]byte, error) {
	return fmt.Appendf(nil, `{"type": "%s", "chat_id": "%s", "user_id": %d}`, b.Scope(), b.ChatID, b.UserID), nil
}

type BotName struct {
	Name string
}

type BotDescription struct {
	Description string `json:"description"`
}

type BotShortDescription struct {
	ShortDescription string `json:"short_description"`
}

type MenuButton struct {
	Type   string      `json:"type"`
	Text   *string     `json:"text,omitempty"`
	WebApp *WebAppInfo `json:"web_app,omitempty"`
}

type ChatBoostSource struct {
	Source            string  `json:"source"`
	User              *User   `json:"user,omitempty"`
	GiveawayMessageId *string `json:"giveaway_message_id,omitempty"`
	PrizeStarCount    *int    `json:"prize_star_count,omitempty"`
	IsUnclaimed       *bool   `json:"is_unclaimed,omitempty"`
}

type ChatBoostSourcePremium struct {
	Source string `json:"source"`
	User   User   `json:"user"`
}

type ChatBoostSourceGiftCode struct {
	Source string `json:"source"`
	User   User   `json:"user"`
}

type ChatBoostSourceGiveaway struct {
	Source            string `json:"source"`
	GiveawayMessageId string `json:"giveaway_message_id"`
	User              *User  `json:"user,omitempty"`
	PrizeStarCount    *int   `json:"prize_star_count,omitempty"`
	IsUnclaimed       *bool  `json:"is_unclaimed,omitempty"`
}

type ChatBoost struct {
	BoostId        string          `json:"boost_id"`
	AddDate        int             `json:"add_date"`
	ExpirationDate int             `json:"expiration_date"`
	Source         ChatBoostSource `json:"source"`
}

type ChatBoostUpdated struct {
	Chat  Chat      `json:"chat"`
	Boost ChatBoost `json:"boost"`
}

type ChatBoostRemoved struct {
	Chat       Chat            `json:"chat"`
	BoostId    string          `json:"boost_id"`
	RemoveDate int             `json:"remove_date"`
	Source     ChatBoostSource `json:"source"`
}

type UserChatBoosts struct {
	Boosts []ChatBoost `json:"boosts"`
}

type BusinessBotRights struct {
	CanReply                   *bool `json:"can_reply,omitempty"`
	CanReadMessages            *bool `json:"can_read_messages,omitempty"`
	CanDeleteOutgoingMessages  *bool `json:"can_delete_outgoing_messages,omitempty"`
	CanDeleteAllMessages       *bool `json:"can_delete_all_messages,omitempty"`
	CanEditName                *bool `json:"can_edit_name,omitempty"`
	CanEditBio                 *bool `json:"can_edit_bio,omitempty"`
	CanEditProfilePhoto        *bool `json:"can_edit_profile_photo,omitempty"`
	CanEditUsername            *bool `json:"can_edit_username,omitempty"`
	CanChangeGiftSettings      *bool `json:"can_change_gift_settings,omitempty"`
	CanViewGiftsAndStars       *bool `json:"can_view_gifts_and_stars,omitempty"`
	CanConvertGiftsGoStars     *bool `json:"can_convert_gifts_go_stars,omitempty"`
	CanTransferAndUpgradeGifts *bool `json:"can_transfer_and_upgrade_gifts,omitempty"`
	CanTransferStars           *bool `json:"can_transfer_stars,omitempty"`
	CanManageStories           *bool `json:"can_manage_stories,omitempty"`
}

type BusinessConnection struct {
	Id         string             `json:"id"`
	User       User               `json:"user"`
	UserChatId int                `json:"user_chat_id"`
	Date       int                `json:"date"`
	Rights     *BusinessBotRights `json:"rights,omitempty"`
	IsEnabled  bool               `json:"is_enabled"`
}

type BusinessMessagesDeleted struct {
	BusinessConnectionId string `json:"business_connection_id"`
	Chat                 Chat   `json:"chat"`
	MessageIds           []int  `json:"message_ids"`
}

type InputMedia interface {
	GetMedia() (media io.Reader)
}

type InputMediaPhoto struct {
	Type                  string           `json:"type"`
	Media                 string           `json:"media"`
	Caption               *string          `json:"caption,omitempty"`
	ParseMode             *string          `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool            `json:"show_caption_above_media,omitempty"`
	HasSpoiler            *bool            `json:"has_spoiler,omitempty"`

	Photo io.Reader `json:"-"`
}

func (m InputMediaPhoto) GetMedia() io.Reader {
	return m.Photo
}

type InputMediaVideo struct {
	Type                  string           `json:"type"`
	Media                 string           `json:"media"`
	Thumbnail             *string          `json:"thumbnail,omitempty"`
	Cover                 *string          `json:"cover,omitempty"`
	StartTimestamp        *int             `json:"start_timestamp,omitempty"`
	Caption               *string          `json:"caption,omitempty"`
	ParseMode             *string          `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool            `json:"show_caption_above_media,omitempty"`
	Width                 *int             `json:"width,omitempty"`
	Height                *int             `json:"height,omitempty"`
	Duration              *int             `json:"duration,omitempty"`
	SupportsStreaming     *bool            `json:"supports_streaming,omitempty"`
	HasSpoiler            *bool            `json:"has_spoiler,omitempty"`

	Video      io.Reader `json:"-"`
	ThumbnailR io.Reader `json:"-"`
}

func (m InputMediaVideo) GetMedia() io.Reader {
	return m.Video
}

func (m InputMediaVideo) GetThumbnail() io.Reader {
	return m.ThumbnailR
}

type InputMediaAnimation struct {
	Media                 string           `json:"media"`
	Thumbnail             *string          `json:"thumbnail,omitempty"`
	Caption               *string          `json:"caption,omitempty"`
	ParseMode             *string          `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool            `json:"show_caption_above_media,omitempty"`
	Width                 *int             `json:"width,omitempty"`
	Height                *int             `json:"height,omitempty"`
	Duration              *int             `json:"duration,omitempty"`
	HasSpoiler            *bool            `json:"has_spoiler,omitempty"`

	Animation  io.Reader `json:"-"`
	ThumbnailR io.Reader `json:"-"`
}

func (m InputMediaAnimation) GetMedia() io.Reader {
	return m.Animation
}

func (m InputMediaAnimation) GetThumbnail() io.Reader {
	return m.ThumbnailR
}

type InputMediaAudio struct {
	Media           string           `json:"media"`
	Thumbnail       *string          `json:"thumbnail,omitempty"`
	Caption         *string          `json:"caption,omitempty"`
	ParseMode       *string          `json:"parse_mode,omitempty"`
	CaptionEntities *[]MessageEntity `json:"caption_entities,omitempty"`
	Duration        *int             `json:"duration,omitempty"`
	Performer       *string          `json:"performer,omitempty"`
	Title           *string          `json:"title,omitempty"`

	Audio      io.Reader `json:"-"`
	ThumbnailR io.Reader `json:"-"`
}

func (m InputMediaAudio) GetMedia() io.Reader {
	return m.Audio
}

func (m InputMediaAudio) GetThumbnail() io.Reader {
	return m.ThumbnailR
}

type InputMediaDocument struct {
	Media                       string           `json:"media"`
	Thumbnail                   *string          `json:"thumbnail,omitempty"`
	Caption                     *string          `json:"caption,omitempty"`
	ParseMode                   *string          `json:"parse_mode,omitempty"`
	CaptionEntities             *[]MessageEntity `json:"caption_entities,omitempty"`
	DisableContentTypeDetection *bool            `json:"disable_content_type_detection,omitempty"`

	Document   io.Reader `json:"-"`
	ThumbnailR io.Reader `json:"-"`
}

func (m InputMediaDocument) GetMedia() io.Reader {
	return m.Document
}

func (m InputMediaDocument) GetThumbnail() io.Reader {
	return m.ThumbnailR
}

type InputFile interface {
	IsRemote() bool
}

type InputFileRemote string

func (i InputFileRemote) IsRemote() bool {
	return true
}

type InputFileLocal struct {
	Data io.Reader
	Name string
}

func (i InputFileLocal) IsRemote() bool {
	return false
}

type InputPaidMedia interface {
	GetPaidMedia() io.Reader
}

type InputPaidMediaPhoto struct {
	Media string `json:"media"`

	Photo io.Reader `json:"-"`
}

func (m InputPaidMediaPhoto) GetPaidMedia() io.Reader {
	return m.Photo
}

type InputPaidMediaVideo struct {
	Media             string  `json:"media"`
	Thumbnail         *string `json:"thumbnail,omitempty,"`
	Cover             *string `json:"cover,omitempty"`
	StartTimestamp    *int    `json:"start_timestamp,omitempty"`
	Width             *int    `json:"width,omitempty,"`
	Height            *int    `json:"height,omitempty,"`
	Duration          *int    `json:"duration,omitempty,"`
	SupportsStreaming *bool   `json:"supports_streaming,omitempty,"`

	Video      io.Reader `json:"-"`
	ThumbnailR io.Reader `json:"-"`
	CoverR     io.Reader `json:"-"`
}

func (m InputPaidMediaVideo) GetPaidMedia() io.Reader {
	return m.Video
}

func (m InputPaidMediaVideo) GetThumbnail() io.Reader {
	return m.ThumbnailR
}

func (m InputPaidMediaVideo) GetCover() io.Reader {
	return m.CoverR
}

type InputProfilePhoto interface {
	GetPhoto() io.Reader
}

type InputProfilePhotoStatic struct {
	Photo string `json:"photo"`

	PhotoR io.Reader `json:"-"`
}

func (pfp InputProfilePhotoStatic) GetPhoto() io.Reader {
	return pfp.PhotoR
}

type InputProfilePhotoAnimated struct {
	Animation          string   `json:"animation"`
	MainFrameTimeStamp *float64 `json:"main_frame_time_stamp,omitempty"`

	AnimationR io.Reader `json:"-"`
}

func (pfp InputProfilePhotoAnimated) GetPhoto() io.Reader {
	return pfp.AnimationR
}

type InputStoryContent interface {
	GetContent() io.Reader
}

type InputStoryContentPhoto struct {
	Photo string `json:"photo"`

	PhotoR io.Reader `json:"-"`
}

func (s InputStoryContentPhoto) GetContent() io.Reader {
	return s.PhotoR
}

type InputStoryContentVideo struct {
	Video               string   `json:"photo"`
	Duration            *float64 `json:"duration,omitempty"`
	CoverFrameTimestamp *float64 `json:"cover_frame_timestamp,omitempty"`
	IsAnimation         bool     `json:"is_animation,omitempty"`

	VideoR io.Reader `json:"-"`
}

func (s InputStoryContentVideo) GetContent() io.Reader {
	return s.VideoR
}

/*
	BEGIN Stickers TYPES
*/

type Sticker struct {
	FileId           string        `json:"file_id"`
	FileUniqueId     string        `json:"file_unique_id"`
	Type             string        `json:"type"`
	Width            int           `json:"width"`
	Height           int           `json:"height"`
	IsAnimated       bool          `json:"is_animated"`
	IsVideo          bool          `json:"is_video"`
	Thumbnail        *PhotoSize    `json:"thumbnail,omitempty"`
	Emoji            *string       `json:"emoji,omitempty"`
	SetName          *string       `json:"set_name,omitempty"`
	PremiumAnimation *File         `json:"premium_animation,omitempty"`
	MaskPosition     *MaskPosition `json:"mask_position,omitempty"`
	CustomEmojiId    *string       `json:"custom_emoji_id,omitempty"`
	NeedsRepainting  *bool         `json:"needs_repainting,omitempty"`
	FileSize         *int          `json:"file_size,omitempty"`
}

type StickerSet struct {
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	StickerType string     `json:"sticker_type"`
	Stickers    []Sticker  `json:"stickers"`
	Thumbnail   *PhotoSize `json:"thumbnail,omitempty"`
}

type MaskPosition struct {
	Point  string   `json:"point"`
	XShift *float32 `json:"x_shift"`
	YShift *float32 `json:"y_shift"`
	Scale  *float32 `json:"scale"`
}

type InputSticker struct {
	Sticker      InputFile     `json:"sticker"`
	Format       string        `json:"format"`
	EmojiList    []string      `json:"emoji_list"`
	MaskPosition *MaskPosition `json:"mask_position,omitempty"`
	Keywords     *[]string     `json:"keywords,omitempty"`
}

/*
	BEGIN Inline Mode TYPES
*/

type InlineQuery struct {
	Id       string    `json:"id"`
	From     User      `json:"from"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
	ChatType *string   `json:"chat_type,omitempty"`
	Location *Location `json:"location,omitempty"`
}

type InlineQueryResultsButton struct {
	Text           string      `json:"text"`
	WebApp         *WebAppInfo `json:"web_app"`
	StartParameter *string     `json:"start_parameter"`
}

type InlineQueryResult interface {
	GetInlineQueryResultType() string
}

type InlineQueryResultArticle struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	Title               string                `json:"title"`
	InputMessageContent InputMessageContent   `json:"input_message_content"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	Url                 *string               `json:"url,omitempty"`
	HideUrl             *bool                 `json:"hide_url,omitempty"`
	Description         *string               `json:"description,omitempty"`
	ThumbnailUrl        *string               `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      *int                  `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     *int                  `json:"thumbnail_height,omitempty"`
}

func (i InlineQueryResultArticle) GetInlineQueryResultType() string {
	return "article"
}

type InlineQueryResultPhoto struct {
	Type                  string                `json:"type"`
	Id                    string                `json:"id"`
	PhotoUrl              string                `json:"photo_url"`
	ThumbnailUrl          string                `json:"thumbnail_url"`
	PhotoWidth            *int                  `json:"photo_width,omitempty"`
	PhotoHeight           *int                  `json:"photo_height,omitempty"`
	Title                 *string               `json:"title,omitempty"`
	Description           *string               `json:"description,omitempty"`
	Caption               *string               `json:"caption,omitempty"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity      `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool                 `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultPhoto) GetInlineQueryResultType() string {
	return "photo"
}

type InlineQueryResultGif struct {
	Type                  string                `json:"type"`
	Id                    string                `json:"id"`
	GifUrl                string                `json:"gif_url"`
	GifWidth              *int                  `json:"gif_width,omitempty"`
	GifHeight             *int                  `json:"gif_height,omitempty"`
	GifDuration           *int                  `json:"gif_duration,omitempty"`
	ThumbnailUrl          string                `json:"thumbnail_url"`
	ThumbnailMimeType     *string               `json:"thumbnail_mime_type,omitempty"`
	Title                 *string               `json:"title,omitempty"`
	Caption               *string               `json:"caption,omitempty"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity      `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool                 `json:"show_caption_above_media,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (i InlineQueryResultGif) GetInlineQueryResultType() string {
	return "gif"
}

type InlineQueryResultMpeg4Gif struct {
	Type                  string                `json:"type"`
	Id                    string                `json:"id"`
	Mpeg4Url              string                `json:"mpeg4_url"`
	Mpeg4Width            *int                  `json:"mpeg4_width,omitempty"`
	Mpeg4Height           *int                  `json:"mpeg4_height,omitempty"`
	Mpeg4Duration         *int                  `json:"mpeg4_duration,omitempty"`
	ThumbnailUrl          string                `json:"thumbnail_url"`
	ThumbnailMimeType     *string               `json:"thumbnail_mime_type,omitempty"`
	Title                 *string               `json:"title,omitempty"`
	Caption               *string               `json:"caption,omitempty"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity      `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool                 `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultMpeg4Gif) GetInlineQueryResultType() string {
	return "mpeg4_gif"
}

type InlineQueryResultVideo struct {
	Type                  string                `json:"type"`
	Id                    string                `json:"id"`
	VideoUrl              string                `json:"video_url"`
	MimeType              string                `json:"mime_type"`
	ThumbnailUrl          string                `json:"thumbnail_url"`
	Title                 string                `json:"title"`
	Caption               *string               `json:"caption,omitempty"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity      `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool                 `json:"show_caption_above_media,omitempty"`
	VideoWidth            *int                  `json:"video_width,omitempty"`
	VideoHeight           *int                  `json:"video_height,omitempty"`
	VideoDuration         *int                  `json:"video_duration,omitempty"`
	Description           *string               `json:"description,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultVideo) GetInlineQueryResultType() string {
	return "video"
}

type InlineQueryResultAudio struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	AudioUrl            string                `json:"audio_url"`
	Title               string                `json:"title"`
	Caption             *string               `json:"caption,omitempty"`
	ParseMode           *string               `json:"parse_mode,omitempty"`
	CaptionEntities     *[]MessageEntity      `json:"caption_entities,omitempty"`
	Performer           *string               `json:"performer,omitempty"`
	AudioDuration       *int                  `json:"audio_duration,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultAudio) GetInlineQueryResultType() string {
	return "audio"
}

type InlineQueryResultVoice struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	VoiceUrl            string                `json:"voice_url"`
	Title               string                `json:"title"`
	Caption             *string               `json:"caption,omitempty"`
	ParseMode           *string               `json:"parse_mode,omitempty"`
	CaptionEntities     *[]MessageEntity      `json:"caption_entities,omitempty"`
	VoiceDuration       *int                  `json:"voice_duration,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultVoice) GetInlineQueryResultType() string {
	return "voice"
}

type InlineQueryResultDocument struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	Title               string                `json:"title"`
	Caption             *string               `json:"caption,omitempty"`
	ParseMode           *string               `json:"parse_mode,omitempty"`
	CaptionEntities     *[]MessageEntity      `json:"caption_entities,omitempty"`
	DocumentUrl         string                `json:"document_url"`
	MimeType            string                `json:"mime_type"`
	Description         *string               `json:"description,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
	ThumbnailUrl        *string               `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      *int                  `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     *int                  `json:"thumbnail_height,omitempty"`
}

func (i InlineQueryResultDocument) GetInlineQueryResultType() string {
	return "document"
}

type InlineQueryResultLocation struct {
	Type                 string                `json:"type"`
	Id                   string                `json:"id"`
	Latitude             *float32              `json:"latitude"`
	Longitude            *float32              `json:"longitude"`
	Title                string                `json:"title"`
	HorizontalAccuracy   *float32              `json:"horizontal_accuracy,omitempty"`
	LivePeriod           *int                  `json:"live_period,omitempty"`
	Heading              *int                  `json:"heading,omitempty"`
	ProximityAlertRadius *int                  `json:"proximity_alert_radius,omitempty"`
	ReplyMarkup          *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent  InputMessageContent   `json:"input_message_content,omitempty"`
	ThumbnailUrl         *string               `json:"thumbnail_url,omitempty"`
	ThumbnailWidth       *int                  `json:"thumbnail_width,omitempty"`
	ThumbnailHeight      *int                  `json:"thumbnail_height,omitempty"`
}

func (i InlineQueryResultLocation) GetInlineQueryResultType() string {
	return "location"
}

type InlineQueryResultVenue struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	Latitude            *float32              `json:"latitude"`
	Longitude           *float32              `json:"longitude"`
	Title               string                `json:"title"`
	Address             string                `json:"address"`
	FoursquareId        *string               `json:"foursquare_id,omitempty"`
	FourSquareType      *string               `json:"four_square_type,omitempty"`
	GooglePlaceId       *string               `json:"google_place_id,omitempty"`
	GooglePlaceType     *string               `json:"google_place_type,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
	ThumbnailUrl        *string               `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      *int                  `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     *int                  `json:"thumbnail_height,omitempty"`
}

func (i InlineQueryResultVenue) GetInlineQueryResultType() string {
	return "venue"
}

type InlineQueryResultContact struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	PhoneNumber         string                `json:"phone_number"`
	FirstName           string                `json:"first_name"`
	LastName            *string               `json:"last_name,omitempty"`
	VCard               *string               `json:"v_card,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
	ThumbnailUrl        string                `json:"thumbnail_url,omitempty"`
	ThumbnailWidth      *int                  `json:"thumbnail_width,omitempty"`
	ThumbnailHeight     *int                  `json:"thumbnail_height,omitempty"`
}

func (i InlineQueryResultContact) GetInlineQueryResultType() string {
	return "contact"
}

type InlineQueryResultGame struct {
	Type          string                `json:"type"`
	Id            string                `json:"id"`
	GameShortName string                `json:"game_short_name"`
	ReplyMarkup   *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

func (i InlineQueryResultGame) GetInlineQueryResultType() string {
	return "game"
}

type InlineQueryResultCachedPhoto struct {
	Type                  string                `json:"type"`
	Id                    string                `json:"id"`
	PhotoFileId           string                `json:"photo_file_id"`
	Title                 *string               `json:"title,omitempty"`
	Description           *string               `json:"description,omitempty"`
	Caption               *string               `json:"caption,omitempty"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity      `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool                 `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultCachedPhoto) GetInlineQueryResultType() string {
	return "photo"
}

type InlineQueryResultCachedGif struct {
	Type                  string                `json:"type"`
	Id                    string                `json:"id"`
	GifFileId             string                `json:"gif_file_id"`
	Title                 *string               `json:"title,omitempty"`
	Caption               *string               `json:"caption,omitempty"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity      `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool                 `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultCachedGif) GetInlineQueryResultType() string {
	return "gif"
}

type InlineQueryResultCachedMpeg4Gif struct {
	Type                  string                `json:"type"`
	Id                    string                `json:"id"`
	Mpeg4FileId           string                `json:"mpeg_4_file_id"`
	Title                 *string               `json:"title,omitempty"`
	Caption               *string               `json:"caption,omitempty"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity      `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool                 `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultCachedMpeg4Gif) GetInlineQueryResultType() string {
	return "mpeg4_gif"
}

type InlineQueryResultCachedSticker struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	StickerFileId       string                `json:"sticker_file_id"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultCachedSticker) GetInlineQueryResultType() string {
	return "sticker"
}

type InlineQueryResultCachedDocument struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	Title               string                `json:"title"`
	DocumentFileId      string                `json:"document_file_id"`
	Description         *string               `json:"description,omitempty"`
	Caption             *string               `json:"caption,omitempty"`
	ParseMode           *string               `json:"parse_mode,omitempty"`
	CaptionEntities     *[]MessageEntity      `json:"caption_entities,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultCachedDocument) GetInlineQueryResultType() string {
	return "document"
}

type InlineQueryResultCachedVideo struct {
	Type                  string                `json:"type"`
	Id                    string                `json:"id"`
	VideoFileId           string                `json:"video_file_id"`
	Title                 string                `json:"title"`
	Description           *string               `json:"description,omitempty"`
	Caption               *string               `json:"caption,omitempty"`
	ParseMode             *string               `json:"parse_mode,omitempty"`
	CaptionEntities       *[]MessageEntity      `json:"caption_entities,omitempty"`
	ShowCaptionAboveMedia *bool                 `json:"show_caption_above_media,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent   InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultCachedVideo) GetInlineQueryResultType() string {
	return "video"
}

type InlineQueryResultCachedVoice struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	VoiceFileId         string                `json:"voice_file_id"`
	Title               string                `json:"title"`
	Caption             *string               `json:"caption,omitempty"`
	ParseMode           *string               `json:"parse_mode,omitempty"`
	CaptionEntities     *[]MessageEntity      `json:"caption_entities,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultCachedVoice) GetInlineQueryResultType() string {
	return "voice"
}

type InlineQueryResultCachedAudio struct {
	Type                string                `json:"type"`
	Id                  string                `json:"id"`
	AudioFileId         string                `json:"audio_file_id"`
	Caption             *string               `json:"caption,omitempty"`
	ParseMode           *string               `json:"parse_mode,omitempty"`
	CaptionEntities     *[]MessageEntity      `json:"caption_entities,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
	InputMessageContent InputMessageContent   `json:"input_message_content,omitempty"`
}

func (i InlineQueryResultCachedAudio) GetInlineQueryResultType() string {
	return "audio"
}

type InputMessageContent interface {
	GetInputMessageContentType() string
}

type InputTextMessageContent struct {
	MessageText        string              `json:"message_text"`
	ParseMode          *string             `json:"parse_mode,omitempty"`
	Entities           *[]MessageEntity    `json:"entities,omitempty"`
	LinkPreviewOptions *LinkPreviewOptions `json:"link_preview_options,omitempty"`
}

func (i InputTextMessageContent) GetInputMessageContentType() string {
	return "text"
}

type InputLocationMessageContent struct {
	Latitude             *float64 `json:"latitude"`
	Longitude            *float64 `json:"longitude"`
	HorizontalAccuracy   *float64 `json:"horizontal_accuracy,omitempty"`
	LivePeriod           *int     `json:"live_period,omitempty"`
	Heading              *int     `json:"heading,omitempty"`
	ProximityAlertRadius *int     `json:"proximity_alert_radius,omitempty"`
}

func (i InputLocationMessageContent) GetInputMessageContentType() string {
	return "location"
}

type InputVenueMessageContent struct {
	Latitude        *float64 `json:"latitude"`
	Longitude       *float64 `json:"longitude"`
	Title           string   `json:"title"`
	Address         string   `json:"address"`
	FoursquareId    *string  `json:"foursquare_id,omitempty"`
	FoursquareType  *string  `json:"foursquare_type,omitempty"`
	GooglePlaceId   *string  `json:"google_place_id,omitempty"`
	GooglePlaceType *string  `json:"google_place_type,omitempty"`
}

func (i InputVenueMessageContent) GetInputMessageContentType() string {
	return "venue"
}

type InputContactMessageContent struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	VCard       string `json:"v_card"`
}

func (i InputContactMessageContent) GetInputMessageContentType() string {
	return "contact"
}

type InputInvoiceMessageContent struct {
	Title                     string         `json:"title"`
	Description               string         `json:"description"`
	Payload                   string         `json:"payload"`
	ProviderToken             *string        `json:"provider_token,omitempty"`
	Currency                  string         `json:"currency"`
	Prices                    []LabeledPrice `json:"prices"`
	MaxTipAmount              *int           `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts       *[]int         `json:"suggested_tip_amounts,omitempty"`
	ProviderData              *string        `json:"provider_data,omitempty"`
	PhotoUrl                  *string        `json:"photo_url,omitempty"`
	PhotoSize                 *int           `json:"photo_size,omitempty"`
	PhotoWidth                *int           `json:"photo_width,omitempty"`
	PhotoHeight               *int           `json:"photo_height,omitempty"`
	NeedName                  *bool          `json:"need_name,omitempty"`
	NeedPhoneNumber           *bool          `json:"need_phone_number,omitempty"`
	NeedEmail                 *bool          `json:"need_email,omitempty"`
	NeedShippingAddress       *bool          `json:"need_shipping_address,omitempty"`
	SendPhoneNumberToProvider *bool          `json:"send_phone_number_to_provider,omitempty"`
	SendEmailToProvider       *bool          `json:"send_email_to_provider,omitempty"`
	IsFlexible                *bool          `json:"is_flexible,omitempty"`
}

func (i InputInvoiceMessageContent) GetInputMessageContentType() string {
	return "invoice"
}

type ChosenInlineResult struct {
	ResultId        string    `json:"result_id"`
	From            User      `json:"from"`
	Location        *Location `json:"location,omitempty"`
	InlineMessageId *string   `json:"inline_message_id,omitempty"`
	Query           string    `json:"query"`
}

type SentWebAppMessage struct {
	InlineMessageId *string `json:"inline_message_id,omitempty"`
}

type PreparedInlineMessage struct {
	Id             string
	ExpirationDate int
}

/*
	BEGIN Payments TYPES
*/

type LabeledPrice struct {
	Label  string `json:"label"`
	Amount int    `json:"amount"`
}

type Invoice struct {
	Title          string `json:"title"`
	Description    string `json:"description"`
	StartParameter string `json:"start_parameter"`
	Currency       string `json:"currency"`
	TotalAmount    string `json:"total_amount"`
}

type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

type OrderInfo struct {
	Name            *string          `json:"name,omitempty"`
	PhoneNumber     *string          `json:"phone_number,omitempty"`
	Email           *string          `json:"email,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

type ShippingOption struct {
	Id     string         `json:"id"`
	Title  string         `json:"title"`
	Prices []LabeledPrice `json:"prices"`
}

type SuccessfulPayment struct {
	Currency                   string     `json:"currency"`
	TotalAmount                string     `json:"total_amount"`
	InvoicePayload             string     `json:"invoice_payload"`
	SubscriptionExpirationDate *int       `json:"subscription_expiration_date,omitempty"`
	IsRecurring                *bool      `json:"is_recurring,omitempty"`
	IsFirstRecurring           *bool      `json:"is_first_recurring,omitempty"`
	ShippingOptionId           *string    `json:"shipping_option_id,omitempty"`
	OrderInfo                  *OrderInfo `json:"order_info,omitempty"`
	TelegramPaymentChargeId    string     `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeId    string     `json:"provider_payment_charge_id"`
}

type RefundedPayment struct {
	Currency                string  `json:"currency"`
	TotalAmount             int     `json:"total_amount"`
	InvoicePayload          string  `json:"invoice_payload"`
	TelegramPaymentChargeId string  `json:"telegram_payment_charge_id"`
	ProviderPaymentChargeId *string `json:"provider_payment_charge_id,omitempty"`
}

type ShippingQuery struct {
	Id              string          `json:"id"`
	From            User            `json:"from"`
	InvoicePayload  string          `json:"invoice_payload"`
	ShippingAddress ShippingAddress `json:"shipping_address"`
}

type PreCheckoutQuery struct {
	Id               string     `json:"id"`
	From             *User      `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int        `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionId *string    `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

type PaidMediaPurchased struct {
	User             User   `json:"user"`
	PaidMediaPayload string `json:"paid_media_payload"`
}

type RevenueWithdrawalState struct {
	Type string `json:"type"`
	Date int    `json:"date"`
	Url  string `json:"url"`
}

type AffiliateInfo struct {
	AffiliateUser     *User `json:"affiliate_user,omitempty"`
	AffiliateChat     *Chat `json:"affiliate_chat,omitempty"`
	CommissionPerMile int   `json:"commission_per_mile"`
	Amount            int   `json:"amount"`
	NanostarAmount    *int  `json:"nanostar_amount,omitempty"`
}

type TransactionPartner struct {
	Type                      string
	TransactionType           string                  `json:"transaction_type"`
	User                      User                    `json:"user"`
	Affiliate                 *AffiliateInfo          `json:"affiliate,omitempty"`
	InvoicePayload            *string                 `json:"invoice_payload,omitempty,"`
	SubscriptionPeriod        *int                    `json:"subscription_period"`
	PaidMedia                 *[]PaidMedia            `json:"paid_media,omitempty,"`
	PaidMediaPayload          *string                 `json:"paid_media_payload,omitempty,"`
	Gift                      *Gift                   `json:"gift,omitempty"`
	PremiumSubscriptionPeriod *int                    `json:"premium_subscription_period"`
	Chat                      Chat                    `json:"chat"`
	SponsorUser               *User                   `json:"sponsor_user,omitempty"`
	CommissionPerMile         int                     `json:"commission_per_mile"`
	WithdrawalState           *RevenueWithdrawalState `json:"withdrawal_state,omitempty"`
	RequestCount              int                     `json:"request_count"`
}

type StarTransaction struct {
	Id             string              `json:"id"`
	Amount         int                 `json:"amount"`
	NanostarAmount *int                `json:"nanostar_amount,omitempty"`
	Date           int                 `json:"date"`
	Source         *TransactionPartner `json:"source,omitempty,"`
	Receiver       *TransactionPartner `json:"receiver,omitempty,"`
}

type StarTransactions struct {
	Transactions []StarTransaction `json:"transactions"`
}

/*
	BEGIN Telegram Passport TYPES
*/

type PassportData struct {
	Data        []EncryptedPassportElement `json:"data"`
	Credentials EncryptedCredentials       `json:"credentials"`
}

type PassportFile struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	FileDate     string `json:"file_date"`
}

type EncryptedPassportElement struct {
	Type         string          `json:"type"`
	Data         *string         `json:"data,omitempty"`
	PhoneNumber  *string         `json:"phone_number,omitempty"`
	Email        *string         `json:"email,omitempty"`
	Files        *[]PassportFile `json:"files,omitempty"`
	FrontSide    *PassportFile   `json:"front_side,omitempty"`
	ReverseSide  *PassportFile   `json:"reverse_side,omitempty"`
	Selfie       *PassportFile   `json:"selfie,omitempty"`
	Translations *[]PassportFile `json:"translations,omitempty"`
	Hash         string          `json:"hash"`
}

type EncryptedCredentials struct {
	Data   string `json:"data"`
	Hash   string `json:"hash"`
	Secret string `json:"secret"`
}

type PassportElementError interface {
	GetPassportElementErrorSource() string
}

type PassportElementErrorDataField struct {
	Source    string `json:"source"`
	Type      string `json:"type"`
	FieldName string `json:"field_name"`
	DataHash  string `json:"data_hash"`
	Message   string `json:"message"`
}

func (p PassportElementErrorDataField) GetPassportElementErrorSource() string {
	return "data"
}

type PassportElementErrorFrontSide struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

func (p PassportElementErrorFrontSide) GetPassportElementErrorSource() string {
	return "front_side"
}

type PassportElementErrorReverseSide struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

func (p PassportElementErrorReverseSide) GetPassportElementErrorSource() string {
	return "reverse_side"
}

type PassportElementErrorSelfie struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

func (p PassportElementErrorSelfie) GetPassportElementErrorSource() string {
	return "selfie"
}

type PassportElementErrorFile struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

func (p PassportElementErrorFile) GetPassportElementErrorSource() string {
	return "file"
}

type PassportElementErrorFiles struct {
	Source     string   `json:"source"`
	Type       string   `json:"type"`
	FileHashes []string `json:"file_hashes"`
	Message    string   `json:"message"`
}

func (p PassportElementErrorFiles) GetPassportElementErrorSource() string {
	return "files"
}

type PassportElementErrorTranslationFile struct {
	Source   string `json:"source"`
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

func (p PassportElementErrorTranslationFile) GetPassportElementErrorSource() string {
	return "translation_file"
}

type PassportElementErrorTranslationFiles struct {
	Source     string   `json:"source"`
	Type       string   `json:"type"`
	FileHashes []string `json:"file_hashes"`
	Message    string   `json:"message"`
}

func (p PassportElementErrorTranslationFiles) GetPassportElementErrorSource() string {
	return "translation_files"
}

type PassportElementErrorUnspecified struct {
	Source      string `json:"source"`
	Type        string `json:"type"`
	ElementHash string `json:"element_hash"`
	Message     string `json:"message"`
}

func (p PassportElementErrorUnspecified) GetPassportElementErrorSource() string {
	return "unspecified"
}

/*
	BEGIN Games TYPES
*/

type Game struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Photo        []PhotoSize      `json:"photo"`
	Text         *string          `json:"text,omitempty"`
	TextEntities *[]MessageEntity `json:"text_entities,omitempty"`
	Animation    *Animation       `json:"animation,omitempty"`
}

type CallbackGame struct{}

type GameHighScore struct {
	Position int  `json:"position"`
	User     User `json:"user"`
	Score    int  `json:"score"`
}
