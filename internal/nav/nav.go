package nav

import (
	"path/filepath"
	"strings"

	"docit/internal/config"
	"docit/internal/parser"
)

type Navigation struct {
	Sidebar   []SidebarGroup
	TableOfContents []HeadingItem
	PrevLink  *LinkItem
	NextLink  *LinkItem
}

type SidebarGroup struct {
	Title string
	Links []LinkItem
}

type LinkItem struct {
	Text string
	URL  string
	Active bool
}

type HeadingItem struct {
	Level   int
	Text    string
	Anchor  string
}

func BuildNavigation(docs []*parser.Document, cfg *config.Config, currentPath string) *Navigation {
	nav := &Navigation{
		Sidebar: buildSidebar(docs, cfg),
	}

	currentRelPath := parser.GetRelativePath(cfg.SourceDir, currentPath)

	nav.PrevLink, nav.NextLink = buildPrevNext(docs, currentRelPath)

	return nav
}

func buildSidebar(docs []*parser.Document, cfg *config.Config) []SidebarGroup {
	var groups []SidebarGroup

	if len(cfg.Sidebar) > 0 {
		for _, sg := range cfg.Sidebar {
			group := SidebarGroup{
				Title: sg.Title,
				Links: make([]LinkItem, 0),
			}
			for _, item := range sg.Items {
				link := findDocByFilename(docs, item)
				if link != nil {
					relPath := parser.GetRelativePath(cfg.SourceDir, link.FilePath)
					htmlPath := strings.TrimSuffix(relPath, ".md") + ".html"
					group.Links = append(group.Links, LinkItem{
						Text: link.Title,
						URL:  htmlPath,
					})
				} else {
					group.Links = append(group.Links, LinkItem{
						Text: item,
						URL:  "#",
					})
				}
			}
			groups = append(groups, group)
		}
	} else {
		group := SidebarGroup{
			Title: "文档",
			Links: make([]LinkItem, 0),
		}
		for _, doc := range docs {
			relPath := parser.GetRelativePath(cfg.SourceDir, doc.FilePath)
			htmlPath := strings.TrimSuffix(relPath, ".md") + ".html"
			group.Links = append(group.Links, LinkItem{
				Text: doc.Title,
				URL:  htmlPath,
			})
		}
		groups = append(groups, group)
	}

	return groups
}

func findDocByFilename(docs []*parser.Document, filename string) *parser.Document {
	for _, doc := range docs {
		if strings.Contains(doc.FilePath, filename) {
			return doc
		}
	}
	return nil
}

func buildPrevNext(docs []*parser.Document, currentPath string) (*LinkItem, *LinkItem) {
	var paths []string
	for _, doc := range docs {
		rel := parser.GetRelativePath("docs", doc.FilePath)
		paths = append(paths, rel)
	}

	currentIdx := -1
	for i, p := range paths {
		if p == currentPath {
			currentIdx = i
			break
		}
	}

	var prev, next *LinkItem

	if currentIdx > 0 {
		doc := docs[currentIdx-1]
		rel := parser.GetRelativePath("docs", doc.FilePath)
		htmlPath := strings.TrimSuffix(rel, ".md") + ".html"
		prev = &LinkItem{
			Text: "← " + doc.Title,
			URL:  htmlPath,
		}
	}

	if currentIdx >= 0 && currentIdx < len(docs)-1 {
		doc := docs[currentIdx+1]
		rel := parser.GetRelativePath("docs", doc.FilePath)
		htmlPath := strings.TrimSuffix(rel, ".md") + ".html"
		next = &LinkItem{
			Text: doc.Title + " →",
			URL:  htmlPath,
		}
	}

	return prev, next
}

func BuildTOC(headings []parser.Heading) []HeadingItem {
	var toc []HeadingItem
	for _, h := range headings {
		if h.Level >= 2 && h.Level <= 3 {
			toc = append(toc, HeadingItem{
				Level:  h.Level,
				Text:   h.Text,
				Anchor: h.Anchor,
			})
		}
	}
	return toc
}

func GetActiveLink(links []LinkItem, currentPath string) []LinkItem {
	var result []LinkItem
	for _, link := range links {
		result = append(result, LinkItem{
			Text:   link.Text,
			URL:    link.URL,
			Active: strings.Contains(currentPath, filepath.Base(link.URL)),
		})
	}
	return result
}