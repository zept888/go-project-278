package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const testBaseURL = "https://short.io"

func TestLinksCRUD(t *testing.T) {
	router := testRouter(t)

	createBody := `{"original_url":"https://example.com/long-url","short_name":"exmpl"}`
	req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(createBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("create: want %d, got %d body=%s", http.StatusCreated, rec.Code, rec.Body.String())
	}

	var created map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created["short_url"] != "https://short.io/r/exmpl" {
		t.Fatalf("unexpected short_url: %v", created["short_url"])
	}

	req = httptest.NewRequest(http.MethodGet, "/api/links/1", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("get: want %d, got %d", http.StatusOK, rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/links", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("list: want %d, got %d", http.StatusOK, rec.Code)
	}
	var list []map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &list); err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("list length: want 1, got %d", len(list))
	}

	updateBody := `{"original_url":"https://example.com/updated","short_name":"exmpl2"}`
	req = httptest.NewRequest(http.MethodPut, "/api/links/1", bytes.NewBufferString(updateBody))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("update: want %d, got %d body=%s", http.StatusOK, rec.Code, rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/r/exmpl2", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("redirect: want %d, got %d", http.StatusFound, rec.Code)
	}
	if rec.Header().Get("Location") != "https://example.com/updated" {
		t.Fatalf("redirect location: %q", rec.Header().Get("Location"))
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/links/1", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("delete: want %d, got %d", http.StatusNoContent, rec.Code)
	}
	if rec.Body.Len() != 0 {
		t.Fatalf("delete body must be empty")
	}
}

func TestCreateWithoutShortName(t *testing.T) {
	router := testRouter(t)

	body := `{"original_url":"https://example.com/auto"}`
	req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("want %d, got %d body=%s", http.StatusCreated, rec.Code, rec.Body.String())
	}
	var created map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatal(err)
	}
	if created["short_name"] == "" {
		t.Fatal("expected generated short_name")
	}
}

func TestShortNameConflict(t *testing.T) {
	router := testRouter(t)

	rec := httptest.NewRecorder()
	for _, payload := range []string{
		`{"original_url":"https://a.com","short_name":"dup"}`,
		`{"original_url":"https://b.com","short_name":"dup"}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
	}

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want %d, got %d body=%s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	errs, _ := resp["errors"].(map[string]any)
	if errs["short_name"] != "short name already in use" {
		t.Fatalf("errors: %v", resp["errors"])
	}
}

func TestCORS(t *testing.T) {
	router := testRouter(t)

	req := httptest.NewRequest(http.MethodOptions, "/api/links", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "GET")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("OPTIONS Allow-Origin: got %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}

	req = httptest.NewRequest(http.MethodGet, "/api/links?range=[0,10]", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Fatalf("GET Allow-Origin: got %q", rec.Header().Get("Access-Control-Allow-Origin"))
	}
	if !strings.Contains(rec.Header().Get("Access-Control-Expose-Headers"), "Content-Range") {
		t.Fatalf("Expose-Headers: got %q", rec.Header().Get("Access-Control-Expose-Headers"))
	}
}

func TestNotFound(t *testing.T) {
	router := testRouter(t)

	req := httptest.NewRequest(http.MethodGet, "/api/links/99", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("get: want %d, got %d", http.StatusNotFound, rec.Code)
	}

	req = httptest.NewRequest(http.MethodPut, "/api/links/99", bytes.NewBufferString(`{"original_url":"https://x.com","short_name":"missing"}`))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("update: want %d, got %d", http.StatusNotFound, rec.Code)
	}

	req = httptest.NewRequest(http.MethodDelete, "/api/links/99", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("delete: want %d, got %d", http.StatusNotFound, rec.Code)
	}
}
