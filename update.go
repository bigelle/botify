package botify

import (
	"encoding/json"
	"fmt"
)

type UpdateType int

const (
	UpdateTypeMessage UpdateType = iota
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

type HandlerFunc func(ctx Context) error

type Context struct {
	updType UpdateType
	// raw bytes of non-empty Update field
	obj []byte
}

func (c *Context) Bind(dest UpdateObject) (err error) {
	if err = json.Unmarshal(c.obj, dest); err != nil {
		return fmt.Errorf("binding update: %w", err)
	}
	return nil
}

type UpdateObject interface {
	// TODO: idk what belongs here
	UpdateType() UpdateType
}
