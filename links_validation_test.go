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

func TestCreateSetsCreatedAt(t *testing.T) {
	router := testRouter(t)

	body := `{"original_url":"https://example.com","short_name":"ctime"}`
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
	if created["created_at"] == nil || created["created_at"] == "" {
		t.Fatalf("expected created_at in response: %v", created)
	}
}

func TestUpdateShortNameConflict(t *testing.T) {
	router := testRouter(t)

	for _, payload := range []string{
		`{"original_url":"https://a.com","short_name":"one"}`,
		`{"original_url":"https://b.com","short_name":"two"}`,
	} {
		req := httptest.NewRequest(http.MethodPost, "/api/links", bytes.NewBufferString(payload))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		if rec.Code != http.StatusCreated {
			t.Fatalf("seed: %d body=%s", rec.Code, rec.Body.String())
		}
	}

	body := `{"original_url":"https://c.com","short_name":"one"}`
	req := httptest.NewRequest(http.MethodPut, "/api/links/2", bytes.NewBufferString(body))
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
	errs, _ := resp["errors"].(map[string]any)
	if errs["short_name"] != "short name already in use" {
		t.Fatalf("errors: %v", resp["errors"])
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
