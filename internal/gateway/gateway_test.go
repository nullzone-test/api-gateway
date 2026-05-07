package gateway

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGatewayRouting(t *testing.T) {
	gw := New()
	gw.Handle("GET", "/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	req := httptest.NewRequest("GET", "/health", nil)
	rec := httptest.NewRecorder()
	gw.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}

func TestGatewayNotFound(t *testing.T) {
	gw := New()
	req := httptest.NewRequest("GET", "/missing", nil)
	rec := httptest.NewRecorder()
	gw.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestMiddlewarePipeline(t *testing.T) {
	gw := New()
	var order []int

	gw.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, 1)
			next.ServeHTTP(w, r)
		})
	})
	gw.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, 2)
			next.ServeHTTP(w, r)
		})
	})

	gw.Handle("GET", "/test", func(w http.ResponseWriter, r *http.Request) {
		order = append(order, 3)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	gw.ServeHTTP(rec, req)

	if len(order) != 3 || order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Errorf("middleware pipeline order wrong: %v", order)
	}
}
