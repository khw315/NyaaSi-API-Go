package scraper

import (
	"fmt"
	"strings"
	"time"
)

type TorrentState string

const (
	StateNormal  TorrentState = "normal"
	StateRemake  TorrentState = "remake"
	StateTrusted TorrentState = "trusted"
)

type Filter int

const (
	FilterNone Filter = iota
	FilterNoRemakes
	FilterTrustedOnly
)

type Ordering string

const (
	OrderAscending  Ordering = "asc"
	OrderDescending Ordering = "desc"
)

type Sort string

const (
	SortComments  Sort = "comments"
	SortSize      Sort = "size"
	SortDate      Sort = "id"
	SortSeeders   Sort = "seeders"
	SortLeechers  Sort = "leechers"
	SortDownloads Sort = "downloads"
)

// Category representation
type Category interface {
	IsSukebei() bool
	GetMainCategoryID() int
	GetSubCategoryID() int
	GetMainCategoryName() string
	GetSubCategoryName() string
	String() string
}

type DefaultCategory struct {
	IsSukebeiFlag bool
	MainID        int
	SubID         int
	MainName      string
	SubName       string
}

func (c DefaultCategory) IsSukebei() bool {
	return c.IsSukebeiFlag
}

func (c DefaultCategory) GetMainCategoryID() int {
	return c.MainID
}

func (c DefaultCategory) GetSubCategoryID() int {
	return c.SubID
}

func (c DefaultCategory) GetMainCategoryName() string {
	return c.MainName
}

func (c DefaultCategory) GetSubCategoryName() string {
	return c.SubName
}

func (c DefaultCategory) String() string {
	if c.MainName != "" && c.SubName != "" {
		return fmt.Sprintf("%s - %s", c.MainName, c.SubName)
	}
	if c.SubName != "" {
		return c.SubName
	}
	return c.MainName
}

const (
	catEnglish    = "English-translated"
	catNonEnglish = "Non-English-translated"
	catRaw        = "Raw"
	catGames      = "Games"
	subNonEnglish = "non-english"
)

func parseSukebeiCategory(mainID, subID int) (mainName, subName string) {
	switch mainID {
	case 1:
		mainName = "Art"
		switch subID {
		case 1:
			subName = "Anime"
		case 2:
			subName = "Doujinshi"
		case 3:
			subName = catGames
		case 4:
			subName = "Manga"
		case 5:
			subName = "Pictures"
		}
	case 2:
		mainName = "Real Life"
		switch subID {
		case 1:
			subName = "Photobooks and Pictures"
		case 2:
			subName = "Videos"
		}
	}
	return
}

func parseNyaaCategory(mainID, subID int) (mainName, subName string) {
	switch mainID {
	case 1:
		mainName = "Anime"
		switch subID {
		case 1:
			subName = "Anime Music Video"
		case 2:
			subName = catEnglish
		case 3:
			subName = catNonEnglish
		case 4:
			subName = catRaw
		}
	case 2:
		mainName = "Audio"
		switch subID {
		case 1:
			subName = "Lossless"
		case 2:
			subName = "Lossy"
		}
	case 3:
		mainName = "Literature"
		switch subID {
		case 1:
			subName = catEnglish
		case 2:
			subName = catNonEnglish
		case 3:
			subName = catRaw
		}
	case 4:
		mainName = "Live Action"
		switch subID {
		case 1:
			subName = catEnglish
		case 2:
			subName = "Idol/Promotional Video"
		case 3:
			subName = catNonEnglish
		case 4:
			subName = catRaw
		}
	case 5:
		mainName = "Pictures"
		switch subID {
		case 1:
			subName = "Graphics"
		case 2:
			subName = "Photos"
		}
	case 6:
		mainName = "Software"
		switch subID {
		case 1:
			subName = "Applications"
		case 2:
			subName = catGames
		}
	}
	return
}

// Helper to construct category from main_id and sub_id
func ParseCategory(mainID, subID int, isSukebei bool) DefaultCategory {
	cat := DefaultCategory{
		IsSukebeiFlag: isSukebei,
		MainID:        mainID,
		SubID:         subID,
	}

	if isSukebei {
		cat.MainName, cat.SubName = parseSukebeiCategory(mainID, subID)
	} else {
		cat.MainName, cat.SubName = parseNyaaCategory(mainID, subID)
	}

	return cat
}

func mapSukebeiCategory(catLower, subLower string) *DefaultCategory {
	switch catLower {
	case "art":
		subMap := map[string]int{
			"anime":     1,
			"doujinshi": 2,
			"games":     3,
			"manga":     4,
			"pictures":  5,
		}
		subID := subMap[subLower]
		c := ParseCategory(1, subID, true)
		return &c
	case "real", "real_life":
		subMap := map[string]int{
			"photobooks": 1,
			"videos":     2,
		}
		subID := subMap[subLower]
		c := ParseCategory(2, subID, true)
		return &c
	}
	return nil
}

func mapNyaaCategory(catLower, subLower string) *DefaultCategory {
	switch catLower {
	case "anime":
		subMap := map[string]int{
			"amv":         1,
			"english":     2,
			subNonEnglish: 3,
			"raw":         4,
		}
		c := ParseCategory(1, subMap[subLower], false)
		return &c
	case "audio":
		subMap := map[string]int{
			"lossless": 1,
			"lossy":    2,
		}
		c := ParseCategory(2, subMap[subLower], false)
		return &c
	case "literature":
		subMap := map[string]int{
			"english":     1,
			subNonEnglish: 2,
			"raw":         3,
		}
		c := ParseCategory(3, subMap[subLower], false)
		return &c
	case "live_action":
		subMap := map[string]int{
			"english":     1,
			"idol_pv":     2,
			subNonEnglish: 3,
			"raw":         4,
		}
		c := ParseCategory(4, subMap[subLower], false)
		return &c
	case "pictures":
		subMap := map[string]int{
			"graphics": 1,
			"photos":   2,
		}
		c := ParseCategory(5, subMap[subLower], false)
		return &c
	case "software":
		subMap := map[string]int{
			"applications": 1,
			"games":        2,
		}
		c := ParseCategory(6, subMap[subLower], false)
		return &c
	}
	return nil
}

// Map string category and subcategory to category struct
func MapCategory(catStr, subCatStr string, isSukebei bool) *DefaultCategory {
	if catStr == "" {
		return nil
	}
	catLower := strings.ToLower(catStr)
	subLower := strings.ToLower(subCatStr)

	if isSukebei {
		return mapSukebeiCategory(catLower, subLower)
	}
	return mapNyaaCategory(catLower, subLower)
}

// SearchRequest holds options for querying Nyaa/Sukebei
type SearchRequest struct {
	Term     string
	Category *DefaultCategory
	Filter   Filter
	User     string
	Page     int
	Ordering Ordering
	SortedBy Sort
}

// TorrentPreview represents an item in the search results table
type TorrentPreview struct {
	ID           int          `json:"id"`
	TorrentState TorrentState `json:"torrent_state"`
	Category     DefaultCategory
	Title        string `json:"title"`
	CommentCount int    `json:"comment_count"`
	DownloadLink string `json:"download_link"`
	MagnetLink   string `json:"magnet_link"`
	Size         string `json:"size"`
	Date         time.Time
	Seeders      int `json:"seeders"`
	Leechers     int `json:"leechers"`
	Completed    int `json:"completed"`
}

func (tp TorrentPreview) ToDict() map[string]interface{} {
	domain := "nyaa.si"
	if tp.Category.IsSukebei() {
		domain = "sukebei.nyaa.si"
	}
	link := fmt.Sprintf("https://%s/view/%d", domain, tp.ID)

	return map[string]interface{}{
		"category":  tp.Category.String(),
		"title":     tp.Title,
		"link":      link,
		"torrent":   tp.DownloadLink,
		"magnet":    tp.MagnetLink,
		"size":      tp.Size,
		"time":      tp.Date.Format("2006-01-02 15:04"),
		"seeders":   tp.Seeders,
		"leechers":  tp.Leechers,
		"downloads": tp.Completed,
	}
}

type SearchResult []TorrentPreview

func (sr SearchResult) ToDict() map[string]interface{} {
	data := make([]map[string]interface{}, 0, len(sr))
	for _, item := range sr {
		data = append(data, item.ToDict())
	}
	return map[string]interface{}{
		"count": len(sr),
		"data":  data,
	}
}

type TorrentComment struct {
	CommentID int       `json:"comment_id"`
	Username  string    `json:"username"`
	IsTrusted bool      `json:"is_trusted"`
	Avatar    string    `json:"avatar"`
	Date      time.Time `json:"date"`
	Text      string    `json:"text"`
}

type TorrentInfo struct {
	Title        string           `json:"title"`
	Description  string           `json:"description"`
	Category     *DefaultCategory `json:"category"`
	Size         string           `json:"size"`
	Date         *time.Time       `json:"date"`
	Uploader     string           `json:"uploader"`
	TorrentState TorrentState     `json:"torrent_state"`
	Seeders      int              `json:"seeders"`
	Leechers     int              `json:"leechers"`
	Completed    int              `json:"completed"`
	Information  string           `json:"information"`
	Hash         string           `json:"hash"`
	DownloadLink string           `json:"download_link"`
	MagnetLink   string           `json:"magnet_link"`
	Comments     []TorrentComment `json:"comments"`
}

func (ti TorrentInfo) ToDict() map[string]interface{} {
	catName := ""
	if ti.Category != nil {
		catName = ti.Category.String()
	}

	dateStr := ""
	if ti.Date != nil {
		dateStr = ti.Date.Format("2006-01-02 15:04:05")
	}

	commentsList := make([]map[string]interface{}, 0, len(ti.Comments))
	for _, c := range ti.Comments {
		cDateStr := ""
		if !c.Date.IsZero() {
			cDateStr = c.Date.Format("2006-01-02 15:04:05")
		}
		commentsList = append(commentsList, map[string]interface{}{
			"user": c.Username,
			"date": cDateStr,
			"text": c.Text,
		})
	}

	return map[string]interface{}{
		"title":         ti.Title,
		"description":   ti.Description,
		"category":      catName,
		"size":          ti.Size,
		"date":          dateStr,
		"uploader":      ti.Uploader,
		"seeders":       ti.Seeders,
		"leechers":      ti.Leechers,
		"completed":     ti.Completed,
		"information":   ti.Information,
		"hash":          ti.Hash,
		"download_link": ti.DownloadLink,
		"magnet_link":   ti.MagnetLink,
		"comments":      commentsList,
	}
}
