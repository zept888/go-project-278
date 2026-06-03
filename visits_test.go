package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func seedLink(t *testing.T, router *gin.Engine, shortName string) {
	t.Helper()
	body := fmt.Sprintf(`{"original_url":"https://example.com/%s","short_name":%q}`, shortName, shortName)
	req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed link: status %d body=%s", rec.Code, rec.Body.String())
	}
}

func TestRedirectRecordsVisit(t *testing.T) {
	router := testRouter(t)
	seedLink(t, router, "hit")

	req := httptest.NewRequest(http.MethodGet, "/r/hit", nil)
	req.Header.Set("User-Agent", "curl/8.5.0")
	req.Header.Set("Referer", "https://ref.example/")
	req.RemoteAddr = "172.18.0.1:12345"
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("redirect: want %d, got %d", http.StatusFound, rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/link_visits", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list visits: want %d, got %d body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}
	var visits []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &visits); err != nil {
		t.Fatal(err)
	}
	if len(visits) != 1 {
		t.Fatalf("want 1 visit, got %d", len(visits))
	}
	v := visits[0]
	if v["link_id"].(float64) != 1 {
		t.Fatalf("link_id: %v", v["link_id"])
	}
	if v["status"].(float64) != 302 {
		t.Fatalf("status: %v", v["status"])
	}
	if v["user_agent"] != "curl/8.5.0" {
		t.Fatalf("user_agent: %v", v["user_agent"])
	}
}

func TestRedirectNotFound(t *testing.T) {
	router := testRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/r/missing", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("want %d, got %d", http.StatusNotFound, rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/link_visits", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	var visits []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &visits); err != nil {
		t.Fatal(err)
	}
	if len(visits) != 0 {
		t.Fatalf("want no visits, got %d", len(visits))
	}
}

func TestLinkVisitsPagination(t *testing.T) {
	router := testRouter(t)
	seedLink(t, router, "p")

	for i := 0; i < 11; i++ {
		req := httptest.NewRequest(http.MethodGet, "/r/p", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusFound {
			t.Fatalf("redirect %d: status %d", i, rec.Code)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/api/link_visits?range=[0,10]", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d body=%s", rec.Code, rec.Body.String())
	}
	if got := rec.Header().Get("Content-Range"); got != "link_visits 0-10/11" {
		t.Fatalf("Content-Range: got %q", got)
	}
	var page []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if len(page) != 10 {
		t.Fatalf("want 10 visits, got %d", len(page))
	}

	req = httptest.NewRequest(http.MethodGet, "/api/link_visits", nil)
	req.Header.Set("Range", "[10, 20]")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if got := rec.Header().Get("Content-Range"); got != "link_visits 10-20/11" {
		t.Fatalf("Content-Range: got %q", got)
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if len(page) != 1 {
		t.Fatalf("want 1 visit, got %d", len(page))
	}
}

func TestLinkVisitsPaginationInvalidRange(t *testing.T) {
	router := testRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/link_visits?range=bad", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want 400, got %d", rec.Code)
	}
}
