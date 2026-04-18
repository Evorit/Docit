# Docit

A lightweight static documentation generator written in Go.

## Overview

Docit transforms your Markdown files into clean, responsive static HTML documentation with minimal JavaScript. It's designed for developers who want simple, fast documentation without the bloat.

## Features

- **Markdown Support** - Full markdown parsing with code highlighting
- **Clean Design** - Minimal, native CSS
- **Responsive** - Works on desktop and mobile
- **Fast** - Generates static HTML, no runtime required
- **Minimal JS** - Only ~50 lines of optional enhancement JS
- **Navigation** - Auto-generated sidebar and table of contents
- **Customizable** - Configurable theme colors and layout

## Quick Start

### Installation

build from source:

```bash
git clone https://github.com/Evorit/Docit.git
cd Docit
go build -o docit ./cmd/docit
```

### Initialize a Project

```bash
docit init
```

This creates:
- `docs/index.md` - Your first documentation page
- `docit.yaml` - Configuration file

### Build Documentation

```bash
docit
```

Output will be in the `dist/` directory.

## Configuration

Create or edit `docit.yaml`:

```yaml
title: "My Docs"
description: "Documentation for my project"
sourceDir: "docs"
outputDir: "dist"
theme:
  primaryColor: "#3b82f6"
nav:
  - text: "Home"
    link: "index.html"
  - text: "Guide"
    link: "guide.html"
sidebar:
  - title: "Getting Started"
    items:
      - "index.md"
      - "guide.md"
```

### Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `title` | Site title | "Documentation" |
| `description` | Site description | "" |
| `sourceDir` | Source markdown directory | "docs" |
| `outputDir` | Output HTML directory | "dist" |
| `theme.primaryColor` | Primary accent color | "#3b82f6" |

## Writing Docs

### Front Matter

Use YAML front matter to set page metadata:

```markdown
---
title: My Page
description: Page description for SEO
---

# My Page Content
```

### Markdown Features

Docit supports standard Markdown:
- Headings, paragraphs, lists
- **Bold** and *italic* text
- [Links](https://example.com)
- Images
- Code blocks with syntax highlighting
- Tables
- Blockquotes

## Philosophy

Docit follows the principle of **minimal JavaScript**:
- All core functionality works without JS
- JS is only used for progressive enhancement (smooth scroll, copy button)
- No client-side frameworks required
- Pure static HTML output

## License

BSD-3-Clause license

see LICENSE file for details.

## Contributing

Contributions are welcome! Please open an issue or submit a PR.