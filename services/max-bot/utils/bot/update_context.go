package bot

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type UpdateContext struct {
	Context context.Context
	Update schemes.UpdateInterface
	Bot *maxbot.Api
}

type HandlerFunc func(*UpdateContext) error

type FilterFunc func(*UpdateContext) bool

type MiddlewareFunc func(HandlerFunc) HandlerFunc

func (ctx *UpdateContext) Answer(text string) error {
	chatID := ctx.Update.GetChatID()
	msg := maxbot.NewMessage().SetChat(chatID).SetText(text)
	return ctx.Bot.Messages.Send(ctx.Context, msg)
}
