package scraper

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const attrDataTimestamp = "data-timestamp"

var (
	categoryURLRegexp = regexp.MustCompile(`/\?c=([0-9]+)_([0-9]+)`)
	viewURLRegexp     = regexp.MustCompile(`/view/([0-9]+)`)
)

func parseCategoryFromURL(href string, isSukebei bool) DefaultCategory {
	matches := categoryURLRegexp.FindStringSubmatch(href)
	if len(matches) == 3 {
		mainID, _ := strconv.Atoi(matches[1])
		subID, _ := strconv.Atoi(matches[2])
		return ParseCategory(mainID, subID, isSukebei)
	}
	return DefaultCategory{IsSukebeiFlag: isSukebei}
}

func parseTorrentIDFromURL(href string) int {
	matches := viewURLRegexp.FindStringSubmatch(href)
	if len(matches) == 2 {
		id, _ := strconv.Atoi(matches[1])
		return id
	}
	return 0
}

func parseTimestamp(tsStr string) time.Time {
	ts, err := strconv.ParseInt(tsStr, 10, 64)
	if err != nil {
		return time.Time{}
	}
	return time.Unix(ts, 0).UTC()
}

func ParseTorrentList(doc *goquery.Document, isSukebei bool) (SearchResult, error) {
	table := doc.Find("table.torrent-list > tbody")
	if table.Length() == 0 {
		html, _ := doc.Html()
		if strings.Contains(strings.ToLower(html), "no torrents found") {
			return SearchResult{}, nil
		}
		return nil, errors.New("cannot find torrent-list table")
	}

	var results SearchResult
	table.Find("tr").Each(func(i int, row *goquery.Selection) {
		state := StateNormal
		if row.HasClass("danger") {
			state = StateRemake
		} else if row.HasClass("success") {
			state = StateTrusted
		}

		cells := row.Find("td")
		if cells.Length() < 8 {
			return
		}

		// Category
		catLink, _ := cells.Eq(0).Find("a").Attr("href")
		cat := parseCategoryFromURL(catLink, isSukebei)

		// Title & ID
		titleCell := cells.Eq(1)
		titleLink := titleCell.Find("a:not(.comments)").First()
		title := strings.TrimSpace(titleLink.Text())
		href, _ := titleLink.Attr("href")
		torrentID := parseTorrentIDFromURL(href)

		// Comments
		commentTag := titleCell.Find("a.comments").First()
		commentCount := 0
		if commentTag.Length() > 0 {
			commentCount, _ = strconv.Atoi(strings.TrimSpace(commentTag.Text()))
		}

		// Links
		links := cells.Eq(2).Find("a")
		base := "https://nyaa.si"
		if isSukebei {
			base = "https://sukebei.nyaa.si"
		}

		downloadHref, _ := links.Eq(0).Attr("href")
		if !strings.HasPrefix(downloadHref, "http") {
			downloadHref = base + downloadHref
		}
		magnetHref, _ := links.Eq(1).Attr("href")

		// Metadata
		size := strings.TrimSpace(cells.Eq(3).Text())
		dateTS, _ := cells.Eq(4).Attr(attrDataTimestamp)
		date := parseTimestamp(dateTS)

		seeders, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(5).Text()))
		leechers, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(6).Text()))
		completed, _ := strconv.Atoi(strings.TrimSpace(cells.Eq(7).Text()))

		results = append(results, TorrentPreview{
			ID:           torrentID,
			TorrentState: state,
			Category:     cat,
			Title:        title,
			CommentCount: commentCount,
			DownloadLink: downloadHref,
			MagnetLink:   magnetHref,
			Size:         size,
			Date:         date,
			Seeders:      seeders,
			Leechers:     leechers,
			Completed:    completed,
		})
	})

	return results, nil
}

func parsePanelBodyCell(rows *goquery.Selection, rowIdx, colIdx int) *goquery.Selection {
	if rowIdx < rows.Length() {
		cols := rows.Eq(rowIdx).Find("div.col-md-5")
		if colIdx < cols.Length() {
			return cols.Eq(colIdx)
		}
	}
	return nil
}

func parsePanelBodyLeftCol(rows *goquery.Selection, info *TorrentInfo, isSukebei bool) {
	if cell := parsePanelBodyCell(rows, 0, 0); cell != nil {
		if catLinks := cell.Find("a"); catLinks.Length() >= 2 {
			href, _ := catLinks.Eq(1).Attr("href")
			cat := parseCategoryFromURL(href, isSukebei)
			info.Category = &cat
		}
	}
	if cell := parsePanelBodyCell(rows, 1, 0); cell != nil {
		if uploaderTag := cell.Find("a").First(); uploaderTag.Length() > 0 {
			info.Uploader = strings.TrimSpace(uploaderTag.Text())
		}
	}
	if cell := parsePanelBodyCell(rows, 2, 0); cell != nil {
		info.Information = strings.TrimSpace(cell.Text())
	}
	if cell := parsePanelBodyCell(rows, 3, 0); cell != nil {
		info.Size = strings.TrimSpace(cell.Text())
	}
	if cell := parsePanelBodyCell(rows, 4, 0); cell != nil {
		info.Hash = strings.TrimSpace(cell.Find("kbd").Text())
	}
}

func parsePanelBodyRightCol(rows *goquery.Selection, info *TorrentInfo) {
	if cell := parsePanelBodyCell(rows, 0, 1); cell != nil {
		tsStr, _ := cell.Attr(attrDataTimestamp)
		t := parseTimestamp(tsStr)
		info.Date = &t
	}
	if cell := parsePanelBodyCell(rows, 1, 1); cell != nil {
		val, _ := strconv.Atoi(strings.TrimSpace(cell.Find("span").Text()))
		info.Seeders = val
	}
	if cell := parsePanelBodyCell(rows, 2, 1); cell != nil {
		val, _ := strconv.Atoi(strings.TrimSpace(cell.Find("span").Text()))
		info.Leechers = val
	}
	if cell := parsePanelBodyCell(rows, 3, 1); cell != nil {
		val, _ := strconv.Atoi(strings.TrimSpace(cell.Text()))
		info.Completed = val
	}
}

func parsePanelBody(body *goquery.Selection, info *TorrentInfo, isSukebei bool) {
	rows := body.Find("div.row")
	parsePanelBodyLeftCol(rows, info, isSukebei)
	parsePanelBodyRightCol(rows, info)
}

func parsePanelFooter(footer *goquery.Selection, info *TorrentInfo, isSukebei bool) {
	base := "https://nyaa.si"
	if isSukebei {
		base = "https://sukebei.nyaa.si"
	}

	dlTag := footer.Find("a[href^=/download/]").First()
	if dlTag.Length() > 0 {
		dlHref, _ := dlTag.Attr("href")
		if !strings.HasPrefix(dlHref, "http") {
			dlHref = base + dlHref
		}
		info.DownloadLink = dlHref
	}

	magnetTag := footer.Find("a[href^=magnet:]").First()
	if magnetTag.Length() == 0 {
		magnetTag = footer.Find("a.card-footer-item").First()
	}
	if magnetTag.Length() > 0 {
		info.MagnetLink, _ = magnetTag.Attr("href")
	}
}

func parseTorrentComments(doc *goquery.Document, info *TorrentInfo) {
	commentDiv := doc.Find("div#comments")
	if commentDiv.Length() == 0 {
		return
	}

	commentDiv.Find("div.comment-panel").Each(func(i int, s *goquery.Selection) {
		userLink := s.Find("a[href^=/user/]").First()
		username := strings.TrimSpace(userLink.Text())

		var cDate time.Time
		tsTag := s.Find("small[" + attrDataTimestamp + "]").First()
		if tsTag.Length() > 0 {
			tsStr, _ := tsTag.Attr(attrDataTimestamp)
			cDate = parseTimestamp(tsStr)
		}

		contentDiv := s.Find("div.comment-content").First()
		text := strings.TrimSpace(contentDiv.Text())

		info.Comments = append(info.Comments, TorrentComment{
			Username: username,
			Date:     cDate,
			Text:     text,
		})
	})
}

func ParseTorrentInfo(doc *goquery.Document, isSukebei bool) (*TorrentInfo, error) {
	var panel *goquery.Selection
	doc.Find("div.panel").Each(func(i int, s *goquery.Selection) {
		if s.Find("div.panel-footer.clearfix").Length() > 0 {
			panel = s
		}
	})

	if panel == nil {
		return nil, errors.New("cannot find main panel")
	}

	info := &TorrentInfo{}

	if panel.HasClass("panel-danger") {
		info.TorrentState = StateRemake
	} else if panel.HasClass("panel-success") {
		info.TorrentState = StateTrusted
	} else {
		info.TorrentState = StateNormal
	}

	info.Title = strings.TrimSpace(panel.Find("div.panel-heading .panel-title").First().Text())

	parsePanelBody(panel.Find("div.panel-body"), info, isSukebei)
	parsePanelFooter(panel.Find("div.panel-footer.clearfix"), info, isSukebei)

	if descDiv := doc.Find("div#torrent-description"); descDiv.Length() > 0 {
		info.Description = strings.TrimSpace(descDiv.Text())
	}

	parseTorrentComments(doc, info)

	return info, nil
}
