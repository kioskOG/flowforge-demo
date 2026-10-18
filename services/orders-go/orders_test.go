package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTotal(t *testing.T) {
	got, err := Total(3, 250)
	if err != nil || got != 750 {
		t.Fatalf("Total(3, 250) = %d, %v; want 750, nil", got, err)
	}
	if _, err := Total(-1, 5); err == nil {
		t.Fatal("Total(-1, 5) should fail")
	}
}

func TestHandler(t *testing.T) {
	rec := httptest.NewRecorder()
	handler(rec, httptest.NewRequest(http.MethodGet, "/total?qty=2&price=100", nil))
	if rec.Code != http.StatusOK || rec.Body.String() != "200\n" {
		t.Fatalf("got %d %q", rec.Code, rec.Body.String())
	}
}
