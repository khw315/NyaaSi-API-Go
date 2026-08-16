package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khw315/NyaaSi-API-Go/pkg/scraper"
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

func TestBuildSearchRequest_AllOptions(t *testing.T) {
	sortOptions := []string{"comments", "size", "id", "date", "seeders", "leechers", "downloads"}
	for _, sortOpt := range sortOptions {
		url := fmt.Sprintf("/nyaa?q=test&category=anime&sub_category=english&sort=%s&order=asc&page=3", sortOpt)
		req := httptest.NewRequest(http.MethodGet, url, nil)
		searchReq := buildSearchRequest(req, false)

		if searchReq.Term != "test" {
			t.Errorf("Expected term test, got %s", searchReq.Term)
		}
		if searchReq.Page != 3 {
			t.Errorf("Expected page 3, got %d", searchReq.Page)
		}
		if searchReq.Ordering != scraper.OrderAscending {
			t.Errorf("Expected asc ordering, got %s", searchReq.Ordering)
		}
		if searchReq.SortedBy == "" {
			t.Errorf("Expected sortedBy to be set for %s", sortOpt)
		}
	}

	// Test descending order and invalid page
	reqDesc := httptest.NewRequest(http.MethodGet, "/nyaa?order=desc&page=abc", nil)
	searchReqDesc := buildSearchRequest(reqDesc, false)
	if searchReqDesc.Ordering != scraper.OrderDescending {
		t.Errorf("Expected desc ordering, got %s", searchReqDesc.Ordering)
	}
	if searchReqDesc.Page != 1 {
		t.Errorf("Expected default page 1, got %d", searchReqDesc.Page)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusBadRequest, "bad request detail")
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

func TestServerHandlers_NotFound(t *testing.T) {
	// Setup mock scraper API returning 404
	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockBackend.Close()

	server := &Server{
		nyaaAPI:    scraper.NewNyaaSiAPIWithBaseURL(mockBackend.URL, false),
		sukebeiAPI: scraper.NewNyaaSiAPIWithBaseURL(mockBackend.URL, true),
	}

	// 1. handleSearchNyaa (404)
	req1 := httptest.NewRequest(http.MethodGet, "/nyaa?q=notfound", nil)
	w1 := httptest.NewRecorder()
	server.handleSearchNyaa(w1, req1)
	if w1.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for handleSearchNyaa, got %d", w1.Code)
	}

	// 2. handleSearchSukebei (404)
	req2 := httptest.NewRequest(http.MethodGet, "/sukebei?q=notfound", nil)
	w2 := httptest.NewRecorder()
	server.handleSearchSukebei(w2, req2)
	if w2.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for handleSearchSukebei, got %d", w2.Code)
	}

	// 3. handleGetNyaaID (404)
	req3 := httptest.NewRequest(http.MethodGet, "/nyaa/id/99999", nil)
	req3.SetPathValue("torrent_id", "99999")
	w3 := httptest.NewRecorder()
	server.handleGetNyaaID(w3, req3)
	if w3.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for handleGetNyaaID, got %d", w3.Code)
	}

	// 4. handleGetSukebeiID (404)
	req4 := httptest.NewRequest(http.MethodGet, "/sukebei/id/99999", nil)
	req4.SetPathValue("torrent_id", "99999")
	w4 := httptest.NewRecorder()
	server.handleGetSukebeiID(w4, req4)
	if w4.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for handleGetSukebeiID, got %d", w4.Code)
	}

	// Invalid ID string
	req3Bad := httptest.NewRequest(http.MethodGet, "/nyaa/id/abc", nil)
	req3Bad.SetPathValue("torrent_id", "abc")
	w3Bad := httptest.NewRecorder()
	server.handleGetNyaaID(w3Bad, req3Bad)
	if w3Bad.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for bad Nyaa ID string, got %d", w3Bad.Code)
	}

	req4Bad := httptest.NewRequest(http.MethodGet, "/sukebei/id/abc", nil)
	req4Bad.SetPathValue("torrent_id", "abc")
	w4Bad := httptest.NewRecorder()
	server.handleGetSukebeiID(w4Bad, req4Bad)
	if w4Bad.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for bad Sukebei ID string, got %d", w4Bad.Code)
	}

	// 5. handleSearchNyaaUser (404)
	req5 := httptest.NewRequest(http.MethodGet, "/nyaa/user/testuser", nil)
	req5.SetPathValue("user_name", "testuser")
	w5 := httptest.NewRecorder()
	server.handleSearchNyaaUser(w5, req5)
	if w5.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for handleSearchNyaaUser, got %d", w5.Code)
	}

	// 6. handleSearchSukebeiUser (404)
	req6 := httptest.NewRequest(http.MethodGet, "/sukebei/user/testuser", nil)
	req6.SetPathValue("user_name", "testuser")
	w6 := httptest.NewRecorder()
	server.handleSearchSukebeiUser(w6, req6)
	if w6.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for handleSearchSukebeiUser, got %d", w6.Code)
	}

	// Empty username
	req5Empty := httptest.NewRequest(http.MethodGet, "/nyaa/user/", nil)
	w5Empty := httptest.NewRecorder()
	server.handleSearchNyaaUser(w5Empty, req5Empty)
	if w5Empty.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for empty Nyaa username, got %d", w5Empty.Code)
	}

	req6Empty := httptest.NewRequest(http.MethodGet, "/sukebei/user/", nil)
	w6Empty := httptest.NewRecorder()
	server.handleSearchSukebeiUser(w6Empty, req6Empty)
	if w6Empty.Code != http.StatusNotFound {
		t.Errorf("Expected 404 for empty Sukebei username, got %d", w6Empty.Code)
	}
}

func TestServerHandlers_Success(t *testing.T) {
	mockListHTML := `<html><body>
	<table class="torrent-list">
		<tbody>
			<tr class="success">
				<td><a href="/?c=1_2"></a></td>
				<td><a href="/view/100" title="Item 100">Item 100</a></td>
				<td><a href="/download/100.torrent">DL</a><a href="magnet:?xt=urn:btih:100">M</a></td>
				<td>10 MiB</td>
				<td data-timestamp="1600000000">Date</td>
				<td>5</td>
				<td>1</td>
				<td>10</td>
			</tr>
		</tbody>
	</table>
	</body></html>`

	mockDetailHTML := `<html><body>
	<div class="panel panel-success">
		<div class="panel-heading"><h3 class="panel-title">Item 100 Title</h3></div>
		<div class="panel-body">
			<div class="row"><div class="col-md-5"><a href="/?c=1_2">Cat</a></div></div>
		</div>
		<div class="panel-footer">
			<a href="/download/100.torrent">DL</a>
			<a href="magnet:?xt=urn:btih:100">M</a>
		</div>
	</div>
	</body></html>`

	mockBackend := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.URL.Path == "/view/100" {
			fmt.Fprintln(w, mockDetailHTML)
			return
		}
		fmt.Fprintln(w, mockListHTML)
	}))
	defer mockBackend.Close()

	server := &Server{
		nyaaAPI:    scraper.NewNyaaSiAPIWithBaseURL(mockBackend.URL, false),
		sukebeiAPI: scraper.NewNyaaSiAPIWithBaseURL(mockBackend.URL, true),
	}

	// 1. handleSearchNyaa
	req1 := httptest.NewRequest(http.MethodGet, "/nyaa?q=test", nil)
	w1 := httptest.NewRecorder()
	server.handleSearchNyaa(w1, req1)
	if w1.Code != http.StatusOK {
		t.Errorf("Expected 200 for handleSearchNyaa, got %d", w1.Code)
	}

	// 2. handleSearchSukebei
	req2 := httptest.NewRequest(http.MethodGet, "/sukebei?q=test", nil)
	w2 := httptest.NewRecorder()
	server.handleSearchSukebei(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("Expected 200 for handleSearchSukebei, got %d", w2.Code)
	}

	// 3. handleGetNyaaID
	req3 := httptest.NewRequest(http.MethodGet, "/nyaa/id/100", nil)
	req3.SetPathValue("torrent_id", "100")
	w3 := httptest.NewRecorder()
	server.handleGetNyaaID(w3, req3)
	if w3.Code != http.StatusOK {
		t.Errorf("Expected 200 for handleGetNyaaID, got %d", w3.Code)
	}

	// 4. handleGetSukebeiID
	req4 := httptest.NewRequest(http.MethodGet, "/sukebei/id/100", nil)
	req4.SetPathValue("torrent_id", "100")
	w4 := httptest.NewRecorder()
	server.handleGetSukebeiID(w4, req4)
	if w4.Code != http.StatusOK {
		t.Errorf("Expected 200 for handleGetSukebeiID, got %d", w4.Code)
	}

	// 5. handleSearchNyaaUser
	req5 := httptest.NewRequest(http.MethodGet, "/nyaa/user/someuser", nil)
	req5.SetPathValue("user_name", "someuser")
	w5 := httptest.NewRecorder()
	server.handleSearchNyaaUser(w5, req5)
	if w5.Code != http.StatusOK {
		t.Errorf("Expected 200 for handleSearchNyaaUser, got %d", w5.Code)
	}

	// 6. handleSearchSukebeiUser
	req6 := httptest.NewRequest(http.MethodGet, "/sukebei/user/someuser", nil)
	req6.SetPathValue("user_name", "someuser")
	w6 := httptest.NewRecorder()
	server.handleSearchSukebeiUser(w6, req6)
	if w6.Code != http.StatusOK {
		t.Errorf("Expected 200 for handleSearchSukebeiUser, got %d", w6.Code)
	}
}
