package scraper

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/PuerkitoBio/goquery"
)

var (
	ErrNotFound = errors.New("torrent not found")
)

type NyaaSiAPI struct {
	isSukebei  bool
	baseURL    string
	httpClient *http.Client
}

func NewNyaaSiAPI(isSukebei bool) *NyaaSiAPI {
	domain := "nyaa.si"
	if isSukebei {
		domain = "sukebei.nyaa.si"
	}
	return &NyaaSiAPI{
		isSukebei: isSukebei,
		baseURL:   fmt.Sprintf("https://%s", domain),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (api *NyaaSiAPI) Search(req SearchRequest) (SearchResult, error) {
	queryParams := url.Values{}

	if req.Term != "" {
		queryParams.Set("q", req.Term)
	}

	if req.Category != nil {
		cat := req.Category
		queryParams.Set("c", fmt.Sprintf("%d_%d", cat.GetMainCategoryID(), cat.GetSubCategoryID()))
	}

	if req.Filter != FilterNone {
		queryParams.Set("f", strconv.Itoa(int(req.Filter)))
	}

	if req.User != "" {
		queryParams.Set("u", req.User)
	}

	if req.Page > 0 {
		queryParams.Set("p", strconv.Itoa(req.Page))
	}

	if req.Ordering != "" {
		queryParams.Set("o", string(req.Ordering))
	}

	if req.SortedBy != "" {
		queryParams.Set("s", string(req.SortedBy))
	}

	reqURL := fmt.Sprintf("%s/?%s", api.baseURL, queryParams.Encode())

	httpReq, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) NyaaSi-API-Go")

	resp, err := api.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return SearchResult{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return ParseTorrentList(doc, api.isSukebei)
}

func (api *NyaaSiAPI) GetTorrentInfo(torrentID int) (*TorrentInfo, error) {
	reqURL := fmt.Sprintf("%s/view/%d", api.baseURL, torrentID)

	httpReq, err := http.NewRequest(http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) NyaaSi-API-Go")

	resp, err := api.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP request failed with status code %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	return ParseTorrentInfo(doc, api.isSukebei)
}
