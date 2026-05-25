package middlewares

import (
	"bytes"
	"context"
	"encoding/gob"
	"log"
	"net/http"
	"time"
	"web-ui/internal/api"

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
	store sessions.Store
}

func NewMiddlewares(store sessions.Store) *Middlewares {
	gob.Register(api.Session{})
	return &Middlewares{store: store}
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

		log.Printf("Auth middleware: time %v session %+v", time.Now(), session)

		ctx := context.WithValue(r.Context(), "session", &session)
		r = r.WithContext(ctx)

		buf := &bytes.Buffer{}
		recorder := &ResponseBuffer{
			ResponseWriter: w,
			body: buf,
			statusCode: http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		log.Printf("Auth middleware: time %v session after handler %+v", time.Now(), session)

		if session.Updated {
			userSession.Values["session"] = session
			userSession.Save(r, w)
		}

		w.WriteHeader(recorder.statusCode)
		_, _ = w.Write(buf.Bytes())
	})
}
