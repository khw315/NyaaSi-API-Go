package scraper

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func TestParseTorrentList_Empty(t *testing.T) {
	html := `<html><body><div class="container"><p>No torrents found</p></div></body></html>`
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		t.Fatalf("Unexpected error creating document: %v", err)
	}

	results, err := ParseTorrentList(doc, false)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

func TestParseCategory(t *testing.T) {
	cat := ParseCategory(1, 2, false)
	if cat.MainName != "Anime" || cat.SubName != "English-translated" {
		t.Errorf("Unexpected category: %v", cat.String())
	}
	if cat.String() != "Anime - English-translated" {
		t.Errorf("Expected 'Anime - English-translated', got '%s'", cat.String())
	}
}

func TestMapCategory(t *testing.T) {
	cat := MapCategory("anime", "english", false)
	if cat == nil {
		t.Fatal("Expected non-nil category")
	}
	if cat.String() != "Anime - English-translated" {
		t.Errorf("Expected 'Anime - English-translated', got '%s'", cat.String())
	}

	sukCat := MapCategory("art", "manga", true)
	if sukCat == nil {
		t.Fatal("Expected non-nil sukebei category")
	}
	if sukCat.String() != "Art - Manga" {
		t.Errorf("Expected 'Art - Manga', got '%s'", sukCat.String())
	}
}
