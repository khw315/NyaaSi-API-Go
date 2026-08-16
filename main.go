package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/khw315/NyaaSi-API-Python/pkg/scraper"
)

type Server struct {
	nyaaAPI    *scraper.NyaaSiAPI
	sukebeiAPI *scraper.NyaaSiAPI
}

func NewServer() *Server {
	return &Server{
		nyaaAPI:    scraper.NewNyaaSiAPI(false),
		sukebeiAPI: scraper.NewNyaaSiAPI(true),
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error writing JSON response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, detail string) {
	writeJSON(w, status, map[string]string{"detail": detail})
}

func buildSearchRequest(r *http.Request, isSukebei bool) scraper.SearchRequest {
	q := r.URL.Query().Get("q")
	category := r.URL.Query().Get("category")
	subCategory := r.URL.Query().Get("sub_category")
	sortStr := r.URL.Query().Get("sort")
	orderStr := r.URL.Query().Get("order")
	pageStr := r.URL.Query().Get("page")

	req := scraper.SearchRequest{
		Term: q,
		Page: 1,
	}

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			req.Page = p
		}
	}

	catObj := scraper.MapCategory(category, subCategory, isSukebei)
	if catObj != nil {
		req.Category = catObj
	}

	if sortStr != "" {
		sortMap := map[string]scraper.Sort{
			"comments":  scraper.SortComments,
			"size":      scraper.SortSize,
			"id":        scraper.SortDate,
			"date":      scraper.SortDate,
			"seeders":   scraper.SortSeeders,
			"leechers":  scraper.SortLeechers,
			"downloads": scraper.SortDownloads,
		}
		if s, ok := sortMap[strings.ToLower(sortStr)]; ok {
			req.SortedBy = s
		}
	}

	if orderStr != "" {
		if strings.ToLower(orderStr) == "desc" {
			req.Ordering = scraper.OrderDescending
		} else if strings.ToLower(orderStr) == "asc" {
			req.Ordering = scraper.OrderAscending
		}
	}

	return req
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Welcome to the (Unofficial) Nyaa API",
		"version": "1.0.0",
		"docs":    "/docs",
		"license": "GPL-3.0 License",
		"github":  "https://github.com/khw315/NyaaSi-API-Python",
	})
}

func (s *Server) handleSearchNyaa(w http.ResponseWriter, r *http.Request) {
	req := buildSearchRequest(r, false)
	results, err := s.nyaaAPI.Search(req)
	if err != nil || len(results) == 0 {
		writeError(w, http.StatusNotFound, "No torrent/magnet link found")
		return
	}
	writeJSON(w, http.StatusOK, results.ToDict())
}

func (s *Server) handleSearchSukebei(w http.ResponseWriter, r *http.Request) {
	req := buildSearchRequest(r, true)
	results, err := s.sukebeiAPI.Search(req)
	if err != nil || len(results) == 0 {
		writeError(w, http.StatusNotFound, "No torrent/magnet link found")
		return
	}
	writeJSON(w, http.StatusOK, results.ToDict())
}

func (s *Server) handleGetNyaaID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("torrent_id")
	torrentID, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusNotFound, "Torrent not found")
		return
	}

	info, err := s.nyaaAPI.GetTorrentInfo(torrentID)
	if err != nil {
		if errors.Is(err, scraper.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Torrent not found")
			return
		}
		writeError(w, http.StatusNotFound, "Torrent not found")
		return
	}

	writeJSON(w, http.StatusOK, info.ToDict())
}

func (s *Server) handleGetSukebeiID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("torrent_id")
	torrentID, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusNotFound, "Torrent not found")
		return
	}

	info, err := s.sukebeiAPI.GetTorrentInfo(torrentID)
	if err != nil {
		if errors.Is(err, scraper.ErrNotFound) {
			writeError(w, http.StatusNotFound, "Torrent not found")
			return
		}
		writeError(w, http.StatusNotFound, "Torrent not found")
		return
	}

	writeJSON(w, http.StatusOK, info.ToDict())
}

func (s *Server) handleSearchNyaaUser(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("user_name")
	if userName == "" {
		writeError(w, http.StatusNotFound, "No torrent/magnet link found")
		return
	}

	req := buildSearchRequest(r, false)
	req.User = userName

	results, err := s.nyaaAPI.Search(req)
	if err != nil || len(results) == 0 {
		writeError(w, http.StatusNotFound, "No torrent/magnet link found")
		return
	}
	writeJSON(w, http.StatusOK, results.ToDict())
}

func (s *Server) handleSearchSukebeiUser(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("user_name")
	if userName == "" {
		writeError(w, http.StatusNotFound, "No torrent/magnet link found")
		return
	}

	req := buildSearchRequest(r, true)
	req.User = userName

	results, err := s.sukebeiAPI.Search(req)
	if err != nil || len(results) == 0 {
		writeError(w, http.StatusNotFound, "No torrent/magnet link found")
		return
	}
	writeJSON(w, http.StatusOK, results.ToDict())
}

func main() {
	server := NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", server.handleHome)
	mux.HandleFunc("GET /nyaa", server.handleSearchNyaa)
	mux.HandleFunc("GET /sukebei", server.handleSearchSukebei)
	mux.HandleFunc("GET /nyaa/id/{torrent_id}", server.handleGetNyaaID)
	mux.HandleFunc("GET /sukebei/id/{torrent_id}", server.handleGetSukebeiID)
	mux.HandleFunc("GET /nyaa/user/{user_name}", server.handleSearchNyaaUser)
	mux.HandleFunc("GET /sukebei/user/{user_name}", server.handleSearchSukebeiUser)

	port := "88"
	log.Printf("Starting NyaaSi API (Go) server on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
