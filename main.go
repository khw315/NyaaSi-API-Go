package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	docs "github.com/khw315/NyaaSi-API-Go/docs"
	"github.com/khw315/NyaaSi-API-Go/pkg/scraper"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

const (
	errNoTorrentFound  = "No torrent/magnet link found"
	errTorrentNotFound = "Torrent not found"
)

// @title           Nyaa-API
// @version         1.0.0
// @description     (Unofficial) Nyaa & Sukebei API built with Go
// @license.name    GPL-3.0 License
// @license.url     https://github.com/khw315/NyaaSi-API-Go/blob/master/LICENSE
// @BasePath        /

type Server struct {
	nyaaAPI    *scraper.NyaaSiAPI
	sukebeiAPI *scraper.NyaaSiAPI
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
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

// handleHome godoc
// @Summary      Home
// @Description  Home Route: Returns app details and health status.
// @Tags         Home
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       / [get]
func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"message": "Welcome to the (Unofficial) Nyaa API",
		"version": "1.0.0",
		"docs":    "/docs",
		"license": "GPL-3.0 License",
		"github":  "https://github.com/khw315/NyaaSi-API-Go",
	})
}

// handleSearchNyaa godoc
// @Summary      Nyaa Category Search
// @Description  Search nyaa.si and return results as JSON
// @Tags         Search
// @Produce      json
// @Param        q             query     string  false  "The search query for torrents."  example(sword art online)
// @Param        category      query     string  false  "Filter torrents by category (e.g., anime, audio, literature, pictures, software)."  example(anime)
// @Param        sub_category  query     string  false  "Filter torrents by subcategory (e.g., english, raw, non-english)."  example(english)
// @Param        sort          query     string  false  "Sort torrents by attribute (comments, size, date, seeders, leechers, downloads)."  example(seeders)
// @Param        order         query     string  false  "Order in which torrents should be sorted (asc or desc)."  example(desc)
// @Param        page          query     int     false  "Page number for pagination."  default(1)
// @Success      200           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]string
// @Router       /nyaa [get]
func (s *Server) handleSearchNyaa(w http.ResponseWriter, r *http.Request) {
	req := buildSearchRequest(r, false)
	results, err := s.nyaaAPI.Search(req)
	if err != nil || len(results) == 0 {
		writeError(w, http.StatusNotFound, errNoTorrentFound)
		return
	}
	writeJSON(w, http.StatusOK, results.ToDict())
}

// handleSearchSukebei godoc
// @Summary      Sukebei Category Search
// @Description  Search sukebei.nyaa.si and return results as JSON
// @Tags         Search
// @Produce      json
// @Param        q             query     string  false  "The search query for torrents."  example(MIDA-512)
// @Param        category      query     string  false  "Filter torrents by category (e.g., art, real)."  example(real)
// @Param        sub_category  query     string  false  "Filter torrents by subcategory (e.g., anime, doujinshi, games, manga, pictures, photobooks, videos)."  example(videos)
// @Param        sort          query     string  false  "Sort torrents by attribute (comments, size, date, seeders, leechers, downloads)."  example(downloads)
// @Param        order         query     string  false  "Order in which torrents should be sorted (asc or desc)."  example(desc)
// @Param        page          query     int     false  "Page number for pagination."  default(1)
// @Success      200           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]string
// @Router       /sukebei [get]
func (s *Server) handleSearchSukebei(w http.ResponseWriter, r *http.Request) {
	req := buildSearchRequest(r, true)
	results, err := s.sukebeiAPI.Search(req)
	if err != nil || len(results) == 0 {
		writeError(w, http.StatusNotFound, errNoTorrentFound)
		return
	}
	writeJSON(w, http.StatusOK, results.ToDict())
}

// handleGetNyaaID godoc
// @Summary      Nyaa Id Search
// @Description  Fetch full details for a specific Nyaa torrent.
// @Tags         ID Search
// @Produce      json
// @Param        torrent_id  path      int  true  "The unique numeric ID of the torrent on nyaa.si."  example(12345)
// @Success      200         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]string
// @Router       /nyaa/id/{torrent_id} [get]
func (s *Server) handleGetNyaaID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("torrent_id")
	torrentID, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusNotFound, errTorrentNotFound)
		return
	}

	info, err := s.nyaaAPI.GetTorrentInfo(torrentID)
	if err != nil {
		writeError(w, http.StatusNotFound, errTorrentNotFound)
		return
	}

	writeJSON(w, http.StatusOK, info.ToDict())
}

// handleGetSukebeiID godoc
// @Summary      Sukebei Id Search
// @Description  Fetch full details for a specific Sukebei torrent.
// @Tags         ID Search
// @Produce      json
// @Param        torrent_id  path      int  true  "The unique numeric ID of the torrent on sukebei.nyaa.si."  example(4505820)
// @Success      200         {object}  map[string]interface{}
// @Failure      404         {object}  map[string]string
// @Router       /sukebei/id/{torrent_id} [get]
func (s *Server) handleGetSukebeiID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("torrent_id")
	torrentID, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusNotFound, errTorrentNotFound)
		return
	}

	info, err := s.sukebeiAPI.GetTorrentInfo(torrentID)
	if err != nil {
		writeError(w, http.StatusNotFound, errTorrentNotFound)
		return
	}

	writeJSON(w, http.StatusOK, info.ToDict())
}

// handleSearchNyaaUser godoc
// @Summary      Nyaa User Search
// @Description  Search for torrents uploaded by a specific user on nyaa.si.
// @Tags         User Search
// @Produce      json
// @Param        user_name     path      string  true   "The Nyaa user name whose uploads you want to search."  example(HorribleSubs)
// @Param        q             query     string  false  "The search query for torrents."
// @Param        category      query     string  false  "Filter torrents by category."
// @Param        sub_category  query     string  false  "Filter torrents by subcategory."
// @Param        sort          query     string  false  "Sort torrents by attribute."
// @Param        order         query     string  false  "Sorting order."
// @Param        page          query     int     false  "Page number."  default(1)
// @Success      200           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]string
// @Router       /nyaa/user/{user_name} [get]
func (s *Server) handleSearchNyaaUser(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("user_name")
	if userName == "" {
		writeError(w, http.StatusNotFound, errNoTorrentFound)
		return
	}

	req := buildSearchRequest(r, false)
	req.User = userName

	results, err := s.nyaaAPI.Search(req)
	if err != nil || len(results) == 0 {
		writeError(w, http.StatusNotFound, errNoTorrentFound)
		return
	}
	writeJSON(w, http.StatusOK, results.ToDict())
}

// handleSearchSukebeiUser godoc
// @Summary      Sukebei User Search
// @Description  Search for torrents uploaded by a specific user on sukebei.nyaa.si.
// @Tags         User Search
// @Produce      json
// @Param        user_name     path      string  true   "The Sukebei user name whose uploads you want to search."  example(offkab)
// @Param        q             query     string  false  "The search query for torrents."
// @Param        category      query     string  false  "Filter torrents by category."
// @Param        sub_category  query     string  false  "Filter torrents by subcategory."
// @Param        sort          query     string  false  "Sort torrents by attribute."
// @Param        order         query     string  false  "Sorting order."
// @Param        page          query     int     false  "Page number."  default(1)
// @Success      200           {object}  map[string]interface{}
// @Failure      404           {object}  map[string]string
// @Router       /sukebei/user/{user_name} [get]
func (s *Server) handleSearchSukebeiUser(w http.ResponseWriter, r *http.Request) {
	userName := r.PathValue("user_name")
	if userName == "" {
		writeError(w, http.StatusNotFound, errNoTorrentFound)
		return
	}

	req := buildSearchRequest(r, true)
	req.User = userName

	results, err := s.sukebeiAPI.Search(req)
	if err != nil || len(results) == 0 {
		writeError(w, http.StatusNotFound, errNoTorrentFound)
		return
	}
	writeJSON(w, http.StatusOK, results.ToDict())
}

func main() {
	// Set dynamic empty host so Swagger UI sends relative requests matching any port/domain
	docs.SwaggerInfo.Host = ""

	server := NewServer()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", server.handleHome)
	mux.HandleFunc("GET /nyaa", server.handleSearchNyaa)
	mux.HandleFunc("GET /sukebei", server.handleSearchSukebei)
	mux.HandleFunc("GET /nyaa/id/{torrent_id}", server.handleGetNyaaID)
	mux.HandleFunc("GET /sukebei/id/{torrent_id}", server.handleGetSukebeiID)
	mux.HandleFunc("GET /nyaa/user/{user_name}", server.handleSearchNyaaUser)
	mux.HandleFunc("GET /sukebei/user/{user_name}", server.handleSearchSukebeiUser)

	// Swagger UI docs endpoint
	mux.Handle("GET /docs/", httpSwagger.WrapHandler)
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/docs/index.html", http.StatusMovedPermanently)
	})

	port := "88"
	log.Printf("Starting NyaaSi API (Go) server on port %s...", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), corsMiddleware(mux)); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
