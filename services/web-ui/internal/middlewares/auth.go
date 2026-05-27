package middlewares

import (
	"bytes"
	"context"
	"encoding/gob"
	"net/http"
	"time"
	"web-ui/internal/api"
	"web-ui/internal/models"
	"web-ui/pkg/cache"

	"github.com/gorilla/sessions"
)

type ResponseBuffer struct {
	http.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

func (rb *ResponseBuffer) WriteHeader(statusCode int) {
	rb.statusCode = statusCode
}

func (rb *ResponseBuffer) Write(b []byte) (int, error) {
	return rb.body.Write(b)
}

type Middlewares struct {
	store      sessions.Store
	userCache  *cache.MemCache[models.User]
	coreClient *api.CoreClient
}

func NewMiddlewares(store sessions.Store, client *api.CoreClient, userCache *cache.MemCache[models.User]) *Middlewares {
	gob.Register(api.Session{})
	return &Middlewares{store: store, userCache: userCache, coreClient: client}
}

func (m *Middlewares) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userSession, _ := m.store.Get(r, "user-session")
		session, ok := userSession.Values["session"].(api.Session)
		if !ok {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				w.WriteHeader(http.StatusUnauthorized)
			} else {
				http.Redirect(w, r, "/login", http.StatusSeeOther)
			}
			return
		}

		user, found := m.userCache.Get(session.SessionID)
		if !found {
			u, err := m.coreClient.GetCurrentUser(r.Context(), &session)
			if err != nil {
				if r.Header.Get("HX-Request") == "true" {
					w.Header().Set("HX-Redirect", "/logout")
					w.WriteHeader(http.StatusOK)
				} else {
					http.Redirect(w, r, "/login", http.StatusSeeOther)
				}
				return
			}
			user = *u
			m.userCache.Set(session.SessionID, user, 15*time.Minute)
		}

		ctx := context.WithValue(r.Context(), "session", &session)
		ctx = context.WithValue(ctx, "user", &user)
		r = r.WithContext(ctx)

		buf := &bytes.Buffer{}
		recorder := &ResponseBuffer{
			ResponseWriter: w,
			body:           buf,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		if session.Updated {
			maxAge := int(time.Until(session.ExiresAt.Add(30 * 24 * time.Hour)).Seconds())
			if maxAge < 0 {
				maxAge = -1
			}
			userSession.Options = &sessions.Options{
				Path: "/",
				MaxAge: maxAge,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			userSession.Values["session"] = session
			userSession.Save(r, w)
		}

		w.WriteHeader(recorder.statusCode)
		_, _ = w.Write(buf.Bytes())
	})
}
