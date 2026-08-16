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

// Helper to construct category from main_id and sub_id
func ParseCategory(mainID, subID int, isSukebei bool) DefaultCategory {
	cat := DefaultCategory{
		IsSukebeiFlag: isSukebei,
		MainID:        mainID,
		SubID:         subID,
	}

	if isSukebei {
		switch mainID {
		case 1:
			cat.MainName = "Art"
			switch subID {
			case 1:
				cat.SubName = "Anime"
			case 2:
				cat.SubName = "Doujinshi"
			case 3:
				cat.SubName = "Games"
			case 4:
				cat.SubName = "Manga"
			case 5:
				cat.SubName = "Pictures"
			}
		case 2:
			cat.MainName = "Real Life"
			switch subID {
			case 1:
				cat.SubName = "Photobooks and Pictures"
			case 2:
				cat.SubName = "Videos"
			}
		}
	} else {
		switch mainID {
		case 1:
			cat.MainName = "Anime"
			switch subID {
			case 1:
				cat.SubName = "Anime Music Video"
			case 2:
				cat.SubName = "English-translated"
			case 3:
				cat.SubName = "Non-English-translated"
			case 4:
				cat.SubName = "Raw"
			}
		case 2:
			cat.MainName = "Audio"
			switch subID {
			case 1:
				cat.SubName = "Lossless"
			case 2:
				cat.SubName = "Lossy"
			}
		case 3:
			cat.MainName = "Literature"
			switch subID {
			case 1:
				cat.SubName = "English-translated"
			case 2:
				cat.SubName = "Non-English-translated"
			case 3:
				cat.SubName = "Raw"
			}
		case 4:
			cat.MainName = "Live Action"
			switch subID {
			case 1:
				cat.SubName = "English-translated"
			case 2:
				cat.SubName = "Idol/Promotional Video"
			case 3:
				cat.SubName = "Non-English-translated"
			case 4:
				cat.SubName = "Raw"
			}
		case 5:
			cat.MainName = "Pictures"
			switch subID {
			case 1:
				cat.SubName = "Graphics"
			case 2:
				cat.SubName = "Photos"
			}
		case 6:
			cat.MainName = "Software"
			switch subID {
			case 1:
				cat.SubName = "Applications"
			case 2:
				cat.SubName = "Games"
			}
		}
	}
	return cat
}

// Map string category and subcategory to category struct
func MapCategory(catStr, subCatStr string, isSukebei bool) *DefaultCategory {
	if catStr == "" {
		return nil
	}
	catLower := strings.ToLower(catStr)
	subLower := strings.ToLower(subCatStr)

	if isSukebei {
		switch catLower {
		case "art":
			switch subLower {
			case "anime":
				c := ParseCategory(1, 1, true)
				return &c
			case "doujinshi":
				c := ParseCategory(1, 2, true)
				return &c
			case "games":
				c := ParseCategory(1, 3, true)
				return &c
			case "manga":
				c := ParseCategory(1, 4, true)
				return &c
			case "pictures":
				c := ParseCategory(1, 5, true)
				return &c
			default:
				c := ParseCategory(1, 0, true)
				return &c
			}
		case "real", "real_life":
			switch subLower {
			case "photobooks":
				c := ParseCategory(2, 1, true)
				return &c
			case "videos":
				c := ParseCategory(2, 2, true)
				return &c
			default:
				c := ParseCategory(2, 0, true)
				return &c
			}
		}
	} else {
		switch catLower {
		case "anime":
			switch subLower {
			case "amv":
				c := ParseCategory(1, 1, false)
				return &c
			case "english":
				c := ParseCategory(1, 2, false)
				return &c
			case "non-english":
				c := ParseCategory(1, 3, false)
				return &c
			case "raw":
				c := ParseCategory(1, 4, false)
				return &c
			default:
				c := ParseCategory(1, 0, false)
				return &c
			}
		case "audio":
			switch subLower {
			case "lossless":
				c := ParseCategory(2, 1, false)
				return &c
			case "lossy":
				c := ParseCategory(2, 2, false)
				return &c
			default:
				c := ParseCategory(2, 0, false)
				return &c
			}
		case "literature":
			switch subLower {
			case "english":
				c := ParseCategory(3, 1, false)
				return &c
			case "non-english":
				c := ParseCategory(3, 2, false)
				return &c
			case "raw":
				c := ParseCategory(3, 3, false)
				return &c
			default:
				c := ParseCategory(3, 0, false)
				return &c
			}
		case "live_action":
			switch subLower {
			case "english":
				c := ParseCategory(4, 1, false)
				return &c
			case "idol_pv":
				c := ParseCategory(4, 2, false)
				return &c
			case "non-english":
				c := ParseCategory(4, 3, false)
				return &c
			case "raw":
				c := ParseCategory(4, 4, false)
				return &c
			default:
				c := ParseCategory(4, 0, false)
				return &c
			}
		case "pictures":
			switch subLower {
			case "graphics":
				c := ParseCategory(5, 1, false)
				return &c
			case "photos":
				c := ParseCategory(5, 2, false)
				return &c
			default:
				c := ParseCategory(5, 0, false)
				return &c
			}
		case "software":
			switch subLower {
			case "applications":
				c := ParseCategory(6, 1, false)
				return &c
			case "games":
				c := ParseCategory(6, 2, false)
				return &c
			default:
				c := ParseCategory(6, 0, false)
				return &c
			}
		}
	}
	return nil
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
