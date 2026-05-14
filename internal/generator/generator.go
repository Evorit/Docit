package generator

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	"docit/internal/config"
	"docit/internal/nav"
	"docit/internal/parser"
)

type Generator struct {
	cfg   *config.Config
	docs  []*parser.Document
}

func New(cfg *config.Config, docs []*parser.Document) *Generator {
	return &Generator{
		cfg:  cfg,
		docs: docs,
	}
}

func (g *Generator) Generate() error {
	if err := os.RemoveAll(g.cfg.OutputDir); err != nil {
		return err
	}

	if err := os.MkdirAll(g.cfg.OutputDir, 0755); err != nil {
		return err
	}

	cssContent := g.generateCSS()
	if err := os.WriteFile(filepath.Join(g.cfg.OutputDir, "style.css"), []byte(cssContent), 0644); err != nil {
		return err
	}

	for _, doc := range g.docs {
		if err := g.generatePage(doc); err != nil {
			return err
		}
	}

	indexPath := filepath.Join(g.cfg.OutputDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		if err := g.generateIndex(); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) generatePage(doc *parser.Document) error {
	navigation := nav.BuildNavigation(g.docs, g.cfg, doc.FilePath)
	toc := nav.BuildTOC(doc.Headings)

	html := g.renderHTML(doc, navigation, toc)
	outputPath := g.cfg.GetOutputPath(doc.FilePath)

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(outputPath, []byte(html), 0644)
}

func (g *Generator) generateIndex() error {
	doc := &parser.Document{
		Title:       g.cfg.Title,
		Content:     "<p>Welcome to " + html.EscapeString(g.cfg.Title) + "</p>",
		FrontMatter: make(map[string]string),
	}

	navigation := nav.BuildNavigation(g.docs, g.cfg, "")
	toc := make([]nav.HeadingItem, 0)

	html := g.renderHTML(doc, navigation, toc)
	return os.WriteFile(filepath.Join(g.cfg.OutputDir, "index.html"), []byte(html), 0644)
}

func (g *Generator) renderHTML(doc *parser.Document, navigation *nav.Navigation, toc []nav.HeadingItem) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>%s</title>
	%s
	%s
</head>
<body>
	<header class="header">
		<div class="header-inner">
			<a href="index.html" class="logo">%s</a>
			%s
			<div class="theme-toggle">
				<button data-mode="light" title="Light">
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="5"/><line x1="12" y1="1" x2="12" y2="3"/><line x1="12" y1="21" x2="12" y2="23"/><line x1="4.22" y1="4.22" x2="5.64" y2="5.64"/><line x1="18.36" y1="18.36" x2="19.78" y2="19.78"/><line x1="1" y1="12" x2="3" y2="12"/><line x1="21" y1="12" x2="23" y2="12"/><line x1="4.22" y1="19.78" x2="5.64" y2="18.36"/><line x1="18.36" y1="5.64" x2="19.78" y2="4.22"/></svg>
				</button>
				<button data-mode="dark" title="Dark">
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>
				</button>
				<button data-mode="system" title="System">
					<svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="3" width="20" height="14" rx="2" ry="2"/><line x1="8" y1="21" x2="16" y2="21"/><line x1="12" y1="17" x2="12" y2="21"/></svg>
				</button>
			</div>
		</div>
	</header>

	<div class="layout">
		<aside class="sidebar">
			<nav class="sidebar-nav">
				%s
			</nav>
		</aside>

		<main class="content">
			<article class="doc">
				%s
				%s
			</article>
		</main>

		%s
	</div>

	%s
</body>
</html>`,
		g.escapeHTML(doc.Title),
		g.generateMeta(doc),
		g.loadCSS(),
		g.escapeHTML(g.cfg.Title),
		g.renderNav(navigation),
		g.renderSidebar(navigation),
		g.renderContent(doc),
		g.renderPrevNext(navigation),
		g.renderTOC(toc),
		g.generateMinimalJS(),
	)
}

func (g *Generator) generateMeta(doc *parser.Document) string {
	var meta string
	if desc, ok := doc.FrontMatter["description"]; ok {
		meta += fmt.Sprintf(`<meta name="description" content="%s">`, g.escapeHTML(desc))
	}
	meta += fmt.Sprintf(`<link rel="stylesheet" href="style.css">`)
	return meta
}

func (g *Generator) loadCSS() string {
	return ""
}

func (g *Generator) escapeHTML(s string) string {
	return html.EscapeString(s)
}

func (g *Generator) renderNav(navigation *nav.Navigation) string {
	if len(g.cfg.Nav) == 0 {
		return ""
	}

	var items []string
	for _, item := range g.cfg.Nav {
		items = append(items, fmt.Sprintf(`<a href="%s" class="nav-link">%s</a>`,
			g.escapeHTML(item.Link), g.escapeHTML(item.Text)))
	}

	return fmt.Sprintf(`<nav class="top-nav">%s</nav>`, strings.Join(items, ""))
}

func (g *Generator) renderSidebar(navigation *nav.Navigation) string {
	if len(navigation.Sidebar) == 0 {
		return ""
	}

	var groups []string
	for _, group := range navigation.Sidebar {
		var links []string
		for _, link := range group.Links {
			activeClass := ""
			if link.Active {
				activeClass = " active"
			}
			links = append(links, fmt.Sprintf(`<a href="%s" class="sidebar-link%s">%s</a>`,
				g.escapeHTML(link.URL), activeClass, g.escapeHTML(link.Text)))
		}
		groups = append(groups, fmt.Sprintf(`<div class="sidebar-group">
			<h3 class="sidebar-group-title">%s</h3>
			<div class="sidebar-links">%s</div>
		</div>`, g.escapeHTML(group.Title), strings.Join(links, "")))
	}

	return strings.Join(groups, "")
}

func (g *Generator) renderContent(doc *parser.Document) string {
	content := strings.ReplaceAll(doc.Content, "<a href=\"", `<a href="`)
	content = strings.ReplaceAll(content, `href="docs/`, `href="`)
	content = strings.ReplaceAll(content, `href="./`, `href="`)
	content = strings.ReplaceAll(content, `src="docs/`, `src="`)
	content = strings.ReplaceAll(content, `src="./`, `src="`)

	headings := doc.Headings
	for _, h := range headings {
		old := fmt.Sprintf(`<h%d>%s</h%d>`, h.Level, h.Text, h.Level)
		new := fmt.Sprintf(`<h%d id="%s">%s</h%d>`, h.Level, h.Anchor, h.Text, h.Level)
		content = strings.ReplaceAll(content, old, new)
	}

	return content
}

func (g *Generator) renderPrevNext(navigation *nav.Navigation) string {
	var parts []string

	if navigation.PrevLink != nil {
		parts = append(parts, fmt.Sprintf(`<a href="%s" class="page-nav prev">%s</a>`,
			g.escapeHTML(navigation.PrevLink.URL), g.escapeHTML(navigation.PrevLink.Text)))
	}

	if navigation.NextLink != nil {
		parts = append(parts, fmt.Sprintf(`<a href="%s" class="page-nav next">%s</a>`,
			g.escapeHTML(navigation.NextLink.URL), g.escapeHTML(navigation.NextLink.Text)))
	}

	if len(parts) == 0 {
		return ""
	}

	return fmt.Sprintf(`<div class="page-nav-wrapper">%s</div>`, strings.Join(parts, ""))
}

func (g *Generator) renderTOC(toc []nav.HeadingItem) string {
	if len(toc) == 0 {
		return ""
	}

	var items []string
	for _, item := range toc {
		indent := ""
		if item.Level == 3 {
			indent = " class=\"toc-level-3\""
		}
		items = append(items, fmt.Sprintf(`<a href="#%s"%s>%s</a>`,
			item.Anchor, indent, g.escapeHTML(item.Text)))
	}

	return fmt.Sprintf(`<aside class="toc">
		<h4 class="toc-title">On this page</h4>
		<nav class="toc-nav">%s</nav>
	</aside>`, strings.Join(items, ""))
}

func (g *Generator) generateMinimalJS() string {
	return `<script>
(function() {
	const sidebar = document.querySelector('.sidebar');
	const header = document.querySelector('.header');
	let lastScroll = 0;

	window.addEventListener('scroll', function() {
		const currentScroll = window.pageYOffset;
		if (currentScroll <= 0) {
			header.classList.remove('hidden');
			return;
		}
		if (currentScroll > lastScroll && currentScroll > 60) {
			header.classList.add('hidden');
		} else {
			header.classList.remove('hidden');
		}
		lastScroll = currentScroll;
	});

	const tocLinks = document.querySelectorAll('.toc-nav a');
	tocLinks.forEach(function(link) {
		link.addEventListener('click', function(e) {
			const id = this.getAttribute('href').substring(1);
			const el = document.getElementById(id);
			if (el) {
				e.preventDefault();
				const top = el.getBoundingClientRect().top + window.pageYOffset - 80;
				window.scrollTo({ top: top, behavior: 'smooth' });
			}
		});
	});

	const copyButton = document.createElement('button');
	copyButton.className = 'copy-button';
	copyButton.textContent = 'Copy';
	document.querySelectorAll('pre').forEach(function(pre) {
		const btn = copyButton.cloneNode(true);
		btn.addEventListener('click', function() {
			const code = pre.querySelector('code');
			if (code) {
				navigator.clipboard.writeText(code.textContent).then(function() {
					btn.textContent = 'Copied!';
					setTimeout(function() { btn.textContent = 'Copy'; }, 2000);
				});
			}
		});
		pre.style.position = 'relative';
		pre.appendChild(btn);
	});
function applyTheme(mode) {
		const buttons = document.querySelectorAll('.theme-toggle button');
		buttons.forEach(function(btn) {
			btn.classList.toggle('active', btn.dataset.mode === mode);
		});
		if (mode === 'system') {
			document.documentElement.removeAttribute('data-theme');
			const systemTheme = getSystemTheme();
			document.documentElement.setAttribute('data-theme', systemTheme);
			localStorage.removeItem('theme');
		} else {
			document.documentElement.setAttribute('data-theme', mode);
			localStorage.setItem('theme', mode);
		}
	}

	function getSystemTheme() {
		return window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
	}

	function initTheme() {
		const savedTheme = localStorage.getItem('theme');
		if (savedTheme) {
			applyTheme(savedTheme);
		} else {
			const systemTheme = getSystemTheme();
			document.documentElement.setAttribute('data-theme', systemTheme);
			const systemBtn = document.querySelector('.theme-toggle button[data-mode="system"]');
			if (systemBtn) systemBtn.classList.add('active');
		}
	}

	initTheme();

	const themeButtons = document.querySelectorAll('.theme-toggle button');
	themeButtons.forEach(function(btn) {
		btn.addEventListener('click', function() {
			applyTheme(this.dataset.mode);
		});
	});

	if (window.matchMedia) {
		window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', function(e) {
			if (!localStorage.getItem('theme')) {
				document.documentElement.setAttribute('data-theme', e.matches ? 'dark' : 'light');
			}
		});
	}
})();
</script>`
}

func (g *Generator) generateCSS() string {
	primaryColor := g.cfg.Theme.PrimaryColor
	if primaryColor == "" {
		primaryColor = "#3b82f6"
	}

	return fmt.Sprintf(`:root {
	--primary: %s;
	--primary-dark: #1d4ed8;
	--text: #1f2937;
	--text-light: #6b7280;
	--bg: #ffffff;
	--bg-secondary: #f9fafb;
	--border: #e5e7eb;
	--code-bg: #f3f4f6;
	--sidebar-width: 280px;
	--toc-width: 220px;
	--header-height: 64px;
}

* {
	margin: 0;
	padding: 0;
	box-sizing: border-box;
}

html {
	scroll-behavior: smooth;
	scroll-padding-top: calc(var(--header-height) + 20px);
}

body {
	font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
	line-height: 1.6;
	color: var(--text);
	background: var(--bg);
}

a {
	color: var(--primary);
	text-decoration: none;
}

a:hover {
	text-decoration: underline;
}

.header {
	position: fixed;
	top: 0;
	left: 0;
	right: 0;
	height: var(--header-height);
	background: var(--bg);
	border-bottom: 1px solid var(--border);
	z-index: 100;
	transition: transform 0.3s ease;
}

.header.hidden {
	transform: translateY(-100%%);
}

.header-inner {
	max-width: 1400px;
	margin: 0 auto;
	padding: 0 24px;
	height: 100%%;
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.logo {
	font-size: 1.25rem;
	font-weight: 600;
	color: var(--text);
}

.top-nav {
	display: flex;
	gap: 24px;
}

.nav-link {
	color: var(--text-light);
	font-size: 0.9rem;
}

.nav-link:hover {
	color: var(--primary);
}

.layout {
	display: flex;
	max-width: 1400px;
	margin: 0 auto;
	padding-top: var(--header-height);
	min-height: 100vh;
}

.sidebar {
	width: var(--sidebar-width);
	flex-shrink: 0;
	position: fixed;
	top: var(--header-height);
	left: 0;
	bottom: 0;
	overflow-y: auto;
	padding: 24px;
	background: var(--bg-secondary);
	border-right: 1px solid var(--border);
}

.sidebar-group {
	margin-bottom: 24px;
}

.sidebar-group-title {
	font-size: 0.75rem;
	font-weight: 600;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	color: var(--text-light);
	margin-bottom: 12px;
}

.sidebar-links {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.sidebar-link {
	display: block;
	padding: 8px 12px;
	border-radius: 6px;
	font-size: 0.9rem;
	color: var(--text);
}

.sidebar-link:hover {
	background: var(--border);
	text-decoration: none;
}

.sidebar-link.active {
	background: var(--primary);
	color: white;
}

.content {
	flex: 1;
	margin-left: var(--sidebar-width);
	padding: 32px 48px;
	max-width: calc(100%% - var(--sidebar-width) - var(--toc-width));
}

.doc {
	max-width: 800px;
}

.doc h1 {
	font-size: 2.25rem;
	font-weight: 700;
	margin-bottom: 16px;
	line-height: 1.2;
}

.doc h2 {
	font-size: 1.75rem;
	font-weight: 600;
	margin-top: 48px;
	margin-bottom: 16px;
	padding-bottom: 8px;
	border-bottom: 1px solid var(--border);
}

.doc h3 {
	font-size: 1.375rem;
	font-weight: 600;
	margin-top: 32px;
	margin-bottom: 12px;
}

.doc h4, .doc h5, .doc h6 {
	font-size: 1.125rem;
	font-weight: 600;
	margin-top: 24px;
	margin-bottom: 12px;
}

.doc p {
	margin-bottom: 16px;
}

.doc ul, .doc ol {
	margin-bottom: 16px;
	padding-left: 24px;
}

.doc li {
	margin-bottom: 8px;
}

.doc code {
	font-family: "SFMono-Regular", Consolas, "Liberation Mono", Menlo, monospace;
	font-size: 0.875em;
	background: var(--code-bg);
	padding: 2px 6px;
	border-radius: 4px;
}

.doc pre {
	background: var(--code-bg);
	padding: 16px;
	border-radius: 8px;
	overflow-x: auto;
	margin-bottom: 16px;
	position: relative;
}

.doc pre code {
	background: none;
	padding: 0;
	font-size: 0.875rem;
	line-height: 1.5;
}

.copy-button {
	position: absolute;
	top: 8px;
	right: 8px;
	padding: 4px 12px;
	font-size: 0.75rem;
	background: var(--bg);
	border: 1px solid var(--border);
	border-radius: 4px;
	cursor: pointer;
	opacity: 0;
	transition: opacity 0.2s;
}

.doc pre:hover .copy-button {
	opacity: 1;
}

.copy-button:hover {
	background: var(--border);
}

.doc blockquote {
	border-left: 4px solid var(--primary);
	padding-left: 16px;
	margin: 16px 0;
	color: var(--text-light);
	font-style: italic;
}

.doc table {
	width: 100%%;
	border-collapse: collapse;
	margin-bottom: 16px;
}

.doc th, .doc td {
	border: 1px solid var(--border);
	padding: 12px;
	text-align: left;
}

.doc th {
	background: var(--bg-secondary);
	font-weight: 600;
}

.doc img {
	max-width: 100%%;
	height: auto;
	border-radius: 8px;
	margin: 16px 0;
}

.doc hr {
	border: none;
	border-top: 1px solid var(--border);
	margin: 32px 0;
}

.toc {
	width: var(--toc-width);
	flex-shrink: 0;
	position: fixed;
	top: var(--header-height);
	right: 0;
	bottom: 0;
	overflow-y: auto;
	padding: 24px;
	border-left: 1px solid var(--border);
	background: var(--bg);
}

.toc-title {
	font-size: 0.75rem;
	font-weight: 600;
	text-transform: uppercase;
	letter-spacing: 0.05em;
	color: var(--text-light);
	margin-bottom: 12px;
}

.toc-nav {
	display: flex;
	flex-direction: column;
	gap: 4px;
}

.toc-nav a {
	font-size: 0.85rem;
	color: var(--text-light);
	padding: 4px 8px;
	border-left: 2px solid transparent;
}

.toc-nav a:hover {
	color: var(--primary);
	border-left-color: var(--primary);
	text-decoration: none;
}

.toc-level-3 {
	padding-left: 16px !important;
}

.page-nav-wrapper {
	display: flex;
	justify-content: space-between;
	margin-top: 48px;
	padding-top: 24px;
	border-top: 1px solid var(--border);
}

.page-nav {
	display: inline-block;
	padding: 12px 20px;
	border: 1px solid var(--border);
	border-radius: 8px;
	font-size: 0.9rem;
}

.page-nav:hover {
	border-color: var(--primary);
	color: var(--primary);
	text-decoration: none;
}

.page-nav.prev {
	margin-right: auto;
}

.page-nav.next {
	margin-left: auto;
}

@media (max-width: 1200px) {
	.toc {
		display: none;
	}
	.content {
		max-width: calc(100%% - var(--sidebar-width));
	}
}

@media (max-width: 768px) {
	.sidebar {
		display: none;
	}
	.content {
		margin-left: 0;
		max-width: 100%%;
		padding: 24px 16px;
	}
	.top-nav {
		display: none;
	}
	.theme-toggle {
		display: none;
	}
}

.theme-toggle {
	display: flex;
	gap: 4px;
	padding: 4px;
	background: var(--bg-secondary);
	border: 1px solid var(--border);
	border-radius: 8px;
}

.theme-toggle button {
	padding: 6px 10px;
	font-size: 0.75rem;
	background: transparent;
	border: none;
	border-radius: 4px;
	cursor: pointer;
	color: var(--text-light);
	transition: all 0.2s;
	display: flex;
	align-items: center;
	justify-content: center;
}

.theme-toggle button svg {
	width: 16px;
	height: 16px;
}

.theme-toggle button:hover {
	background: var(--border);
}

.theme-toggle button.active {
	background: var(--primary);
	color: white;
}

[data-theme="dark"] body,
[data-theme="dark"] .header,
[data-theme="dark"] .toc {
	--text: #e5e7eb;
	--text-light: #9ca3af;
	--bg: #1a1a1a;
	--bg-secondary: #262626;
	--border: #404040;
	--code-bg: #262626;
}
`, primaryColor)
}