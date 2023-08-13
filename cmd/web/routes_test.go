package main

import (
	"fmt"
	"github.com/ddytert/bookings/internal/config"
	"github.com/go-chi/chi/v5"
	"testing"
)

func TestRoutes(t *testing.T) {
	ac := config.AppConfig{}
	mux := routes(&ac)

	switch v := mux.(type) {
	case *chi.Mux:
	// do nothing, test passed
	default:
		t.Error(fmt.Sprintf("Type is not http.Handler, but is %T", v))
	}
}
