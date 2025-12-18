package bot

import (
	"context"
	"fmt"
	"log"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Dispatcher struct {
	bot     *maxbot.Api
	router  *Router
	polling bool
}

func NewDispatcher(bot *maxbot.Api) *Dispatcher {
	return &Dispatcher{
		bot:    bot,
		router: NewRouter(),
	}
}

func (d *Dispatcher) IncludeRouter(router *Router) {
	d.router.IncludeRouter(router)
}

func (d *Dispatcher) StartPolling(ctx context.Context) {
	d.polling = true

	for update := range d.bot.GetUpdates(ctx) {
		go d.processUpdate(ctx, update)
	}
}

func (d *Dispatcher) processUpdate(ctx context.Context, update schemes.UpdateInterface) {
	uctx := &UpdateContext{
		Context: ctx,
		Update:  update,
		Bot:     d.bot,
	}

	var updateType string
	switch update.(type) {
	case *schemes.MessageCreatedUpdate:
		updateType = "message"
	case *schemes.MessageCallbackUpdate:
		updateType = "callback_query"
	default:
		return
	}

	logText := fmt.Sprintf("Update type: %s, ", updateType)

	handlers := d.router.getHandlers(updateType)

	for _, registration := range handlers {
		if registration.Filter(uctx) {
			err := registration.Handler(uctx)
			if err != nil {
				logText += fmt.Sprintf("handlers error: %v", err)
			} else {
				logText += "processed"
			}
			break
		}
	}

	log.Println(logText)
}
