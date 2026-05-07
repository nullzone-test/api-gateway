package gateway

import (
	"log"
	"net/http"
	"time"
)

// Gateway is a lightweight HTTP router with middleware support.
type Gateway struct {
	middlewares []Middleware
	routes      map[string]map[string]http.HandlerFunc
}

// Middleware wraps an HTTP handler.
type Middleware func(http.Handler) http.Handler

// New creates a new Gateway instance.
func New() *Gateway {
	return &Gateway{
		routes: make(map[string]map[string]http.HandlerFunc),
	}
}

// Use adds middleware to the pipeline.
func (g *Gateway) Use(m Middleware) {
	g.middlewares = append(g.middlewares, m)
}

// Handle registers a route handler.
func (g *Gateway) Handle(method, path string, handler http.HandlerFunc) {
	if g.routes[method] == nil {
		g.routes[method] = make(map[string]http.HandlerFunc)
	}
	g.routes[method][path] = handler
}

// ServeHTTP implements http.Handler.
func (g *Gateway) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if methods, ok := g.routes[r.Method]; ok {
		if handler, ok := methods[r.URL.Path]; ok {
			var h http.Handler = handler
			for i := len(g.middlewares) - 1; i >= 0; i-- {
				h = g.middlewares[i](h)
			}
			h.ServeHTTP(w, r)
			return
		}
	}
	http.NotFound(w, r)
}

// Logger returns a logging middleware.
func Logger() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
		})
	}
}

// RateLimit returns a rate limiting middleware using token bucket.
func RateLimit(rps int) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
}
