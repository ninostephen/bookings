package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// NoSurf adds CSRF protection to all post requests.
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProducation,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

// LoadMiddleware loads and saves session of every request
func LoadMiddleware(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
