package bot

import (
	"strings"

	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Filter struct{}

var Filters = &Filter{}

func (f *Filter) Command(cmd string) FilterFunc {
	return func(ctx *UpdateContext) bool {
		if upd, ok := ctx.Update.(*schemes.MessageCreatedUpdate); ok {
			msgText := upd.Message.Body.Text
			return strings.HasPrefix(msgText, "/") &&
				strings.Split(strings.TrimPrefix(msgText, "/"), " ")[0] == cmd
		}

		return false
	}
}

func (f *Filter) Text(text string) FilterFunc {
	return func(ctx *UpdateContext) bool {
		if upd, ok := ctx.Update.(*schemes.MessageCreatedUpdate); ok {
			return upd.Message.Body.Text == text
		}

		return false
	}
}

func (f *Filter) StartsWith(text string) FilterFunc {
	return func(ctx *UpdateContext) bool {
		if upd, ok := ctx.Update.(*schemes.MessageCreatedUpdate); ok {
			msgText := upd.Message.Body.Text
			return strings.HasPrefix(msgText, text)
		}
		return false
	}
}

func (f *Filter) And(filters ...FilterFunc) FilterFunc {
	return func(ctx *UpdateContext) bool {
		for _, filter := range filters {
			if !filter(ctx) {
				return false
			}
		}
		return true
	}
}

func (f *Filter) Or(filters ...FilterFunc) FilterFunc {
	return func(ctx *UpdateContext) bool {
		for _, filter := range filters {
			if filter(ctx) {
				return true
			}
		}
		return false
	}
}
