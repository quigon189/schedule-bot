package bot

import (
	"context"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

type Dispatcher struct {
	bot *maxbot.Api
	router *Router	
	polling bool
}

func NewDispatcher(bot *maxbot.Api) *Dispatcher {
	return &Dispatcher{
		bot: bot,
		router: NewRouter(),
	}
}

func (d *Dispatcher) IncludeRouter(router *Router) {
	d.router.IncludeRouter(router)
}

func (d *Dispatcher) StartPolling(ctx context.Context) error {
	d.polling = true

	for update := range d.bot.GetUpdates(ctx) {
		go  d.processUpdate(update)
	}

	return nil
}

func (d *Dispatcher) processUpdate(update schemes.UpdateInterface) {

}
