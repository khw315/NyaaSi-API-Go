package scraper

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPI_Search(t *testing.T) {
	mockHTML := `<html><body>
	<table class="torrent-list">
		<tbody>
			<tr class="success">
				<td><a href="/?c=1_2" title="Anime - English-translated"></a></td>
				<td>
					<a href="/view/12345" title="Sample Torrent Title">Sample Torrent Title</a>
				</td>
				<td>
					<a href="/download/12345.torrent"><i class="fa-download"></i></a>
					<a href="magnet:?xt=urn:btih:1234567890abcdef"><i class="fa-magnet"></i></a>
				</td>
				<td>100.0 MiB</td>
				<td data-timestamp="1600000000">2020-09-13</td>
				<td>10</td>
				<td>5</td>
				<td>100</td>
			</tr>
		</tbody>
	</table>
	</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") == "notfound" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Query().Get("q") == "error" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockHTML)
	}))
	defer server.Close()

	api := NewNyaaSiAPIWithBaseURL(server.URL, false)

	cat := ParseCategory(1, 2, false)
	req := SearchRequest{
		Term:     "test",
		Category: &cat,
		Filter:   FilterTrustedOnly,
		User:     "testuser",
		Page:     1,
		Ordering: OrderDescending,
		SortedBy: SortSeeders,
	}

	results, err := api.Search(req)
	if err != nil {
		t.Fatalf("Unexpected error in Search: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
	if results[0].Title != "Sample Torrent Title" {
		t.Errorf("Expected title 'Sample Torrent Title', got '%s'", results[0].Title)
	}

	dict := results.ToDict()
	if dict["count"] != 1 {
		t.Errorf("Expected dict count 1, got %v", dict["count"])
	}

	notFoundResults, err := api.Search(SearchRequest{Term: "notfound"})
	if err != nil {
		t.Fatalf("Unexpected error on 404: %v", err)
	}
	if len(notFoundResults) != 0 {
		t.Errorf("Expected empty results on 404, got %d", len(notFoundResults))
	}

	_, err500 := api.Search(SearchRequest{Term: "error"})
	if err500 == nil {
		t.Error("Expected error on status 500")
	}
}

func TestAPI_GetTorrentInfo(t *testing.T) {
	mockDetailHTML := `<html><body>
	<div class="panel panel-success">
		<div class="panel-heading"><h3 class="panel-title">Detailed Torrent Title</h3></div>
		<div class="panel-body">
			<div class="row">
				<div class="col-md-5"><a href="#">Category</a><a href="/?c=1_2">Anime - English</a></div>
				<div class="col-md-5" data-timestamp="1600000000">Date</div>
			</div>
			<div class="row">
				<div class="col-md-5"><a href="#">UploaderName</a></div>
				<div class="col-md-5"><span>50</span></div>
			</div>
			<div class="row">
				<div class="col-md-5">Info text</div>
				<div class="col-md-5"><span>10</span></div>
			</div>
			<div class="row">
				<div class="col-md-5">500 MiB</div>
				<div class="col-md-5">200</div>
			</div>
			<div class="row">
				<div class="col-md-5"><kbd>1234567890ABCDEF</kbd></div>
			</div>
		</div>
		<div class="panel-footer">
			<a href="/download/12345.torrent">Download</a>
			<a href="magnet:?xt=urn:btih:12345">Magnet</a>
		</div>
	</div>
	<div id="torrent-description">Full markdown description</div>
	<div id="comments">
		<div class="comment-panel">
			<div class="panel-body">
				<a href="/user/commenter" class="username">commenter</a>
				<small data-timestamp="1600000000"></small>
				<div class="comment-content">Great release!</div>
			</div>
		</div>
	</div>
	</body></html>`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/view/999" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if r.URL.Path == "/view/500" {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, mockDetailHTML)
	}))
	defer server.Close()

	api := NewNyaaSiAPIWithBaseURL(server.URL, false)

	info, err := api.GetTorrentInfo(12345)
	if err != nil {
		t.Fatalf("Unexpected error in GetTorrentInfo: %v", err)
	}
	if info.Title != "Detailed Torrent Title" {
		t.Errorf("Expected title 'Detailed Torrent Title', got '%s'", info.Title)
	}
	if info.Uploader != "UploaderName" {
		t.Errorf("Expected uploader 'UploaderName', got '%s'", info.Uploader)
	}
	if len(info.Comments) != 1 {
		t.Fatalf("Expected 1 comment, got %d", len(info.Comments))
	}
	if info.Comments[0].Text != "Great release!" {
		t.Errorf("Expected comment 'Great release!', got '%s'", info.Comments[0].Text)
	}

	dict := info.ToDict()
	if dict["title"] != "Detailed Torrent Title" {
		t.Errorf("Expected dict title 'Detailed Torrent Title', got %v", dict["title"])
	}

	_, err404 := api.GetTorrentInfo(999)
	if err404 != ErrNotFound {
		t.Errorf("Expected ErrNotFound on 404, got %v", err404)
	}

	_, err500 := api.GetTorrentInfo(500)
	if err500 == nil {
		t.Error("Expected error on status 500")
	}
}

func TestModels_AllCategories(t *testing.T) {
	// Test Constructor
	apiNyaa := NewNyaaSiAPI(false)
	if apiNyaa.isSukebei {
		t.Error("Expected Nyaa API flag false")
	}
	apiSukebei := NewNyaaSiAPI(true)
	if !apiSukebei.isSukebei {
		t.Error("Expected Sukebei API flag true")
	}

	// Test Nyaa categories
	for mainID := 1; mainID <= 6; mainID++ {
		for subID := 1; subID <= 4; subID++ {
			cat := ParseCategory(mainID, subID, false)
			if cat.IsSukebei() {
				t.Error("Expected Nyaa category flag false")
			}
			_ = cat.GetMainCategoryName()
			_ = cat.GetSubCategoryName()
			_ = cat.GetMainCategoryID()
			_ = cat.GetSubCategoryID()
		}
	}

	// Test Sukebei categories
	for mainID := 1; mainID <= 2; mainID++ {
		for subID := 1; subID <= 5; subID++ {
			cat := ParseCategory(mainID, subID, true)
			if !cat.IsSukebei() {
				t.Error("Expected Sukebei category flag true")
			}
		}
	}

	// Test MapCategory for all branches
	nyaaCats := []string{"anime", "audio", "literature", "live_action", "pictures", "software"}
	for _, c := range nyaaCats {
		_ = MapCategory(c, "raw", false)
	}
	sukebeiCats := []string{"art", "real", "real_life"}
	for _, c := range sukebeiCats {
		_ = MapCategory(c, "videos", true)
	}

	if MapCategory("", "", false) != nil {
		t.Error("Expected nil for empty category string")
	}
	if MapCategory("invalid", "invalid", false) != nil {
		t.Error("Expected nil for invalid category string")
	}

	// Test DefaultCategory String formatting
	catMainOnly := DefaultCategory{MainName: "MainOnly"}
	if catMainOnly.String() != "MainOnly" {
		t.Errorf("Expected 'MainOnly', got '%s'", catMainOnly.String())
	}
	catSubOnly := DefaultCategory{SubName: "SubOnly"}
	if catSubOnly.String() != "SubOnly" {
		t.Errorf("Expected 'SubOnly', got '%s'", catSubOnly.String())
	}
}
