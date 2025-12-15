package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"vk-parser/internal/config"
	"vk-parser/internal/ocr"
	"vk-parser/internal/parser"
	"vk-parser/internal/vk"
	"vk-parser/internal/webhook"
)

func main() {
	cfg := config.Load()
	vkClient := vk.NewClient(cfg.VKToken, cfg.VKGroupID)
	ocrClient := ocr.NewTesseract(cfg.TesseractPath)
	sender := webhook.NewSender(cfg.ScheduleServiceURL + "/api/v1/changes")

	vkParser := parser.NewParser(vkClient, ocrClient, sender, cfg)

	ctx, cancel := context.WithCancel(context.Background())

	log.Println("Start parsing...")
	log.Printf("Parse interval: %v minutes", cfg.CheckInterval.Minutes())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		cancel()
	}()

	vkParser.Run(ctx)

	log.Println("Parser stoped")
}
