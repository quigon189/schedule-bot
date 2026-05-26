package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"web-ui/internal/api"
	"web-ui/internal/config"
	"web-ui/internal/models"
	"web-ui/internal/router"
	cache "web-ui/pkg"

	"github.com/gorilla/sessions"
)

func main() {
	cfg := config.Load()

	cookieStore := sessions.NewCookieStore([]byte(cfg.CookieSecret))
	cookieStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   cfg.SessionMaxAge,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	coreClient := api.NewCoreClient(cfg.CoreURL, time.Duration(cfg.CoreTimeout)*time.Second, 30 * time.Second)

	userCache := cache.NewMemCache[models.User](5*time.Minute)

	// csrfMiddleware := csrf.Protect(
	// 	[]byte("secret-key-from-config"),
	// 	csrf.Secure(false), // true для https
	// 	csrf.Path("/"),
	// 	csrf.RequestHeader("X-CSRF-Token"),
	// )

	r := router.NewRouter(coreClient, cookieStore, userCache)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r.Handler(),
	}

	log.Printf("Server web-ui started on :%s", cfg.Port)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server stoped")
}
