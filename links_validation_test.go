package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateInvalidURL(t *testing.T) {
	router := testRouter(t)

	body := `{"original_url":"not-a-url","short_name":"abc"}`
	req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want %d, got %d body=%s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
	}
	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	errs, ok := resp["errors"].(map[string]any)
	if !ok {
		t.Fatalf("errors: %v", resp)
	}
	msg, _ := errs["original_url"].(string)
	if !strings.Contains(msg, "failed on the 'url' tag") {
		t.Fatalf("original_url message: %q", msg)
	}
}

func TestCreateInvalidJSON(t *testing.T) {
	router := testRouter(t)

	req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(`{`))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("want %d, got %d", http.StatusBadRequest, rec.Code)
	}
	var resp map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatal(err)
	}
	if resp["error"] != "invalid request" {
		t.Fatalf("error: %v", resp["error"])
	}
}

func TestCreateShortNameTooShort(t *testing.T) {
	router := testRouter(t)

	body := `{"original_url":"https://example.com","short_name":"ab"}`
	req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want %d, got %d body=%s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
	}
}

func TestUpdateInvalidURL(t *testing.T) {
	router := testRouter(t)

	seed := `{"original_url":"https://example.com","short_name":"seed1"}`
	req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(seed))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("seed: %d body=%s", rec.Code, rec.Body.String())
	}

	body := `{"original_url":"bad","short_name":"seed2"}`
	req = httptest.NewRequest(http.MethodPut, "/api/links/1", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("want %d, got %d body=%s", http.StatusUnprocessableEntity, rec.Code, rec.Body.String())
	}
}
