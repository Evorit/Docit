package parser

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type Document struct {
	Title       string
	Content     string
	RawContent  string
	FilePath    string
	Headings    []Heading
	FrontMatter map[string]string
}

type Heading struct {
	Level   int
	Text    string
	Anchor  string
}

func Parse(mdPath string) (*Document, error) {
	content, err := os.ReadFile(mdPath)
	if err != nil {
		return nil, err
	}

	doc := &Document{
		FilePath:   mdPath,
		RawContent: string(content),
	}

	doc.FrontMatter, doc.RawContent = parseFrontMatter(doc.RawContent)

	extensions := parser.CommonExtensions | parser.Attributes
	p := parser.NewWithExtensions(extensions)
	htmlContent := markdown.ToHTML([]byte(doc.RawContent), p, nil)
	doc.Content = string(htmlContent)

	doc.Headings = extractHeadings(doc.RawContent)
	doc.Title = extractTitle(doc.RawContent, doc.FrontMatter)

	return doc, nil
}

func parseFrontMatter(content string) (map[string]string, string) {
	frontMatter := make(map[string]string)
	if !strings.HasPrefix(content, "---") {
		return frontMatter, content
	}

	lines := strings.Split(content, "\n")
	if len(lines) < 3 {
		return frontMatter, content
	}

	var bodyLines []string
	inFrontMatter := false
	for i, line := range lines {
		if i == 0 && strings.TrimSpace(line) == "---" {
			inFrontMatter = true
			continue
		}
		if i > 0 && strings.TrimSpace(line) == "---" {
			bodyLines = lines[i+1:]
			break
		}
		if inFrontMatter {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				frontMatter[key] = value
			}
		}
	}

	if bodyLines == nil {
		return frontMatter, content
	}
	return frontMatter, strings.Join(bodyLines, "\n")
}

func extractHeadings(content string) []Heading {
	var headings []Heading
	lines := strings.Split(content, "\n")
	headingRegex := regexp.MustCompile(`^(#{1,6})\s+(.+)$`)

	for _, line := range lines {
		matches := headingRegex.FindStringSubmatch(line)
		if len(matches) == 3 {
			level := len(matches[1])
			text := strings.TrimSpace(matches[2])
			anchor := generateAnchor(text)
			headings = append(headings, Heading{
				Level:  level,
				Text:   text,
				Anchor: anchor,
			})
		}
	}
	return headings
}

func extractTitle(content string, frontMatter map[string]string) string {
	if title, ok := frontMatter["title"]; ok && title != "" {
		return title
	}

	lines := strings.Split(content, "\n")
	headingRegex := regexp.MustCompile(`^#\s+(.+)$`)
	for _, line := range lines {
		matches := headingRegex.FindStringSubmatch(line)
		if len(matches) == 2 {
			return strings.TrimSpace(matches[1])
		}
	}

	return ""
}

func generateAnchor(text string) string {
	anchor := strings.ToLower(text)
	anchor = strings.ReplaceAll(anchor, " ", "-")
	re := regexp.MustCompile(`[^\w\-]`)
	anchor = re.ReplaceAllString(anchor, "")
	return anchor
}

func GetRelativePath(basePath, fullPath string) string {
	rel, err := filepath.Rel(basePath, fullPath)
	if err != nil {
		return fullPath
	}
	return rel
}