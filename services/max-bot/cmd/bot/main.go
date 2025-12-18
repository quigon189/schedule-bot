package main

import (
	"context"
	"fmt"
	"log"
	"max-bot/utils/bot"
	"os"
	"os/signal"
	"syscall"

	maxbot "github.com/max-messenger/max-bot-api-client-go"
	"github.com/max-messenger/max-bot-api-client-go/schemes"
)

func main() {
	token := os.Getenv("MAX_BOT_TOKEN")
	if token == "" {
		log.Fatal("MAX_BOT_TOKEN not set")
	}

	api, err := maxbot.New(token)
	if err != nil {
		log.Fatalf("Failed to init max bot: %v", err)
	}
	_, err = api.Bots.GetBot(context.Background())
	if err != nil {
		log.Fatalf("Failed to get bot info: %v", err)
	}
	dp := bot.NewDispatcher(api)

	echoRouter := bot.NewRouter()

	echoRouter.Use(func(next bot.HandlerFunc) bot.HandlerFunc {
		return func(ctx *bot.UpdateContext) error {
			log.Printf("New update from %d", ctx.Update.GetUserID())
			return next(ctx)
		}
	})

	echoRouter.Message(func(uc *bot.UpdateContext) error {
		if upd, ok := uc.Update.(*schemes.MessageCreatedUpdate); ok {
			return uc.Answer(upd.GetText())
		}
		return fmt.Errorf("invalide update type")
	},
		func(uc *bot.UpdateContext) bool {
			return true
		},
	)

	dp.IncludeRouter(echoRouter)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		cancel()
	}()

	dp.StartPolling(ctx)
}
