package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLinksPagination(t *testing.T) {
	router := testRouter(t)

	for i := 1; i <= 11; i++ {
		body := fmt.Sprintf(`{"original_url":"https://example.com/%d","short_name":"s%d"}`, i, i)
		req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("seed %d: status %d body=%s", i, rec.Code, rec.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/links?range=[0,10]", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if got := rec.Header().Get("Content-Range"); got != "links 0-10/11" {
		t.Fatalf("Content-Range: got %q", got)
	}
	var page []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if len(page) != 10 {
		t.Fatalf("want 10 items, got %d", len(page))
	}

	req = httptest.NewRequest(http.MethodGet, "/api/links?range=[5,10]", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if got := rec.Header().Get("Content-Range"); got != "links 5-10/11" {
		t.Fatalf("Content-Range: got %q", got)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if len(page) != 5 {
		t.Fatalf("want 5 items, got %d", len(page))
	}
	if id, _ := page[0]["id"].(float64); id != 6 {
		t.Fatalf("first id want 6, got %v", page[0]["id"])
	}
}

func TestLinksPaginationInvalidRange(t *testing.T) {
	router := testRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/links?range=bad", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
}
