package parser

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"
	"vk-parser/internal/config"
	"vk-parser/internal/ocr"
	"vk-parser/internal/vk"
	"vk-parser/internal/webhook"

	"github.com/SevereCloud/vksdk/v3/object"
)

type Parser struct {
	vkClient *vk.Client
	ocr      *ocr.Tesseract
	sender   *webhook.Sender
	config   *config.Config

	processed map[int]bool
	mu        sync.RWMutex
}

func NewParser(vkClient *vk.Client, ocr *ocr.Tesseract, sender *webhook.Sender, cfg *config.Config) *Parser {
	return &Parser{
		vkClient:  vkClient,
		ocr:       ocr,
		sender:    sender,
		config:    cfg,
		processed: make(map[int]bool),
	}
}

func (p *Parser) Run(ctx context.Context) {
	ticker := time.NewTicker(p.config.CheckInterval)
	defer ticker.Stop()

	p.parsePosts()

	for {
		select {
		case <-ticker.C:
			p.parsePosts()
		case <-ctx.Done():
			return
		}
	}
}

func (p *Parser) parsePosts() {
	posts, err := p.vkClient.GetWallPosts(p.config.MaxPosts)
	if err != nil {
		log.Printf("Failed to get posts: %v", err)
		return
	}

	for _, post := range posts {
		if p.isProcessed(post.ID) {
			continue
		}

		p.processPost(post)	
	}
}

func (p *Parser) processPost(post object.WallWallpost) {
	urls := []string{}
	log.Printf("Processing post: %d", post.ID)
	for _, attachment := range post.Attachments {
		if attachment.Type == "photo" {
			photo := attachment.Photo
			urls = append(urls, photo.Sizes[len(photo.Sizes)-1].URL)
		}
	}

	for _, url := range urls {
		imageBytes, err := downloadImage(url)
		if err != nil {
			log.Printf("Failed to download image: %v", err)
			continue
		}

		text, err := p.ocr.ExtractTextFromImage(imageBytes)
		if err != nil {
			log.Printf("Failed to extract text: %v", err)
			continue
		}

		if date := checkPattern(text, p.config.Pattern); date != "" {
			data := map[string]any{
				"date":        date,
				"description": post.Text,
				"image_urls":  urls,
			}

			if err := p.sender.Send(data); err != nil {
				log.Printf("Failed to send parsed data: %v", err)
			} else {
				p.markProcessed(post.ID)
			}
			return
		}
	}
}

func (p *Parser) isProcessed(postID int) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	return p.processed[postID]
}

func (p *Parser) markProcessed(postID int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	log.Printf("Post with id %d processed", postID)
	p.processed[postID] = true
}

func downloadImage(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error %v", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func checkPattern(text, pattern string) string {
	re, err := regexp.Compile("Изменения в расписании на (\\d{2})\\.(\\d{2})\\.(\\d{4})")
	if err != nil {
		return ""
	}

	match := re.FindStringSubmatch(text)

	if match == nil {
		return ""
	}

	return fmt.Sprintf("%s-%s-%s", match[3], match[2], match[1])
}
