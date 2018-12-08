package gotgbot

import (
	"encoding/json"

	"github.com/PaulSonOfLars/gotgbot/ext"
)

type Handler interface {
	HandleUpdate(u *Update, d Dispatcher) error
	CheckUpdate(u *Update) (bool, error)
	GetName() string
}

type Update struct {
	UpdateId           int                     `json:"update_id"`
	Message            *ext.Message            `json:"message"`
	EditedMessage      *ext.Message            `json:"edited_message"`
	ChannelPost        *ext.Message            `json:"channel_post"`
	EditedChannelPost  *ext.Message            `json:"edited_channel_post"`
	InlineQuery        *ext.Message            `json:"inline_query"`
	ChosenInlineResult *ext.ChosenInlineResult `json:"chosen_inline_result"`
	CallbackQuery      *ext.CallbackQuery      `json:"callback_query"`
	ShippingQuery      *ext.ShippingQuery      `json:"shipping_query"`
	PreCheckoutQuery   *ext.PreCheckoutQuery   `json:"pre_checkout_query"`

	// Self added type
	EffectiveMessage *ext.Message `json:"effective_message"`
	EffectiveChat    *ext.Chat    `json:"effective_chat"`
	EffectiveUser    *ext.User    `json:"effective_user"`
	Data             map[string]string
}

// todo: move this into dispatcher update processor to updater CPU cycles
func initUpdate(data json.RawMessage, bot ext.Bot) *Update {
	var upd Update
	json.Unmarshal(data, &upd)
	if upd.Message != nil {
		upd.EffectiveMessage = upd.Message
		upd.EffectiveChat = upd.Message.Chat
		upd.EffectiveUser = upd.Message.From

	} else if upd.EditedMessage != nil {
		upd.EffectiveMessage = upd.EditedMessage
		upd.EffectiveChat = upd.EditedMessage.Chat
		upd.EffectiveUser = upd.EditedMessage.From

	} else if upd.ChannelPost != nil {
		upd.EffectiveMessage = upd.ChannelPost
		upd.EffectiveChat = upd.ChannelPost.Chat

	} else if upd.EditedChannelPost != nil {
		upd.EffectiveMessage = upd.EditedChannelPost
		upd.EffectiveChat = upd.EditedChannelPost.Chat

	} else if upd.InlineQuery != nil {
		upd.EffectiveMessage = upd.InlineQuery
		upd.EffectiveUser = upd.InlineQuery.From

	} else if upd.CallbackQuery != nil && upd.CallbackQuery.Message != nil {
		upd.EffectiveMessage = upd.CallbackQuery.Message
		upd.EffectiveChat = upd.CallbackQuery.Message.Chat
		upd.EffectiveUser = upd.CallbackQuery.From

	} else if upd.ChosenInlineResult != nil {
		upd.EffectiveUser = upd.ChosenInlineResult.From

	} else if upd.ShippingQuery != nil {
		upd.EffectiveUser = upd.ShippingQuery.From

	} else if upd.PreCheckoutQuery != nil {
		upd.EffectiveUser = upd.PreCheckoutQuery.From
	}

	if upd.EffectiveMessage != nil {
		upd.EffectiveMessage.Bot = bot
		if upd.EffectiveMessage.ReplyToMessage != nil {
			upd.EffectiveMessage.ReplyToMessage.Bot = bot
			if upd.EffectiveMessage.ReplyToMessage.From != nil {
				upd.EffectiveMessage.ReplyToMessage.From.Bot = bot
			}
		}
	}
	if upd.EffectiveChat != nil {
		upd.EffectiveChat.Bot = bot
	}
	if upd.EffectiveUser != nil {
		upd.EffectiveUser.Bot = bot
	}
	upd.Data = make(map[string]string)
	return &upd
}
