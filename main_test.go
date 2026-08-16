package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHome(t *testing.T) {
	server := NewServer()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	server.handleHome(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got %v", resp.Status)
	}

	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if body["version"] != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", body["version"])
	}
}

func TestCorsMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	mw := corsMiddleware(handler)

	// Test OPTIONS preflight
	reqOptions := httptest.NewRequest(http.MethodOptions, "/nyaa", nil)
	wOptions := httptest.NewRecorder()
	mw.ServeHTTP(wOptions, reqOptions)
	if wOptions.Code != http.StatusOK {
		t.Errorf("Expected 200 for OPTIONS, got %d", wOptions.Code)
	}
	if wOptions.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("Expected CORS origin header")
	}

	// Test GET request
	reqGet := httptest.NewRequest(http.MethodGet, "/nyaa", nil)
	wGet := httptest.NewRecorder()
	mw.ServeHTTP(wGet, reqGet)
	if wGet.Code != http.StatusOK {
		t.Errorf("Expected 200 for GET, got %d", wGet.Code)
	}
}

func TestBuildSearchRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/nyaa?q=test&category=anime&sub_category=english&sort=seeders&order=desc&page=2", nil)
	searchReq := buildSearchRequest(req, false)

	if searchReq.Term != "test" {
		t.Errorf("Expected term test, got %s", searchReq.Term)
	}
	if searchReq.Page != 2 {
		t.Errorf("Expected page 2, got %d", searchReq.Page)
	}
	if searchReq.Category == nil {
		t.Error("Expected category to be set")
	}
}

func TestInvalidIDs(t *testing.T) {
	server := NewServer()

	// Invalid ID for handleGetNyaaID
	reqNyaa := httptest.NewRequest(http.MethodGet, "/nyaa/id/invalid", nil)
	wNyaa := httptest.NewRecorder()
	server.handleGetNyaaID(wNyaa, reqNyaa)
	if wNyaa.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for invalid Nyaa ID, got %d", wNyaa.Code)
	}

	// Invalid ID for handleGetSukebeiID
	reqSukebei := httptest.NewRequest(http.MethodGet, "/sukebei/id/invalid", nil)
	wSukebei := httptest.NewRecorder()
	server.handleGetSukebeiID(wSukebei, reqSukebei)
	if wSukebei.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for invalid Sukebei ID, got %d", wSukebei.Code)
	}
}

func TestEmptyUser(t *testing.T) {
	server := NewServer()

	// Empty username
	reqNyaa := httptest.NewRequest(http.MethodGet, "/nyaa/user/", nil)
	wNyaa := httptest.NewRecorder()
	server.handleSearchNyaaUser(wNyaa, reqNyaa)
	if wNyaa.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for empty Nyaa user, got %d", wNyaa.Code)
	}

	reqSukebei := httptest.NewRequest(http.MethodGet, "/sukebei/user/", nil)
	wSukebei := httptest.NewRecorder()
	server.handleSearchSukebeiUser(wSukebei, reqSukebei)
	if wSukebei.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for empty Sukebei user, got %d", wSukebei.Code)
	}
}
