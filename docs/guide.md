---
title: Guide
description: User guide for Docit
---

# User Guide

This guide explains how to use Docit to create your documentation.

## Installation

Download the binary for your platform from the releases page.

## Configuration

Create a `docit.yaml` file in your project root:

```yaml
title: "My Documentation"
sourceDir: "docs"
outputDir: "dist"
theme:
  primaryColor: "#3b82f6"
```

## Writing Docs

Add markdown files to the `docs` directory. Each file will become a page.

### Front Matter

You can use YAML front matter to set page title and description:

```yaml
---
title: My Page
description: This is my page description
---
```

### Markdown Features

Docit supports standard markdown including:

- **Bold** and *italic* text
- [Links](https://example.com)
- Images
- Code blocks
- Tables

## Building

Run the build command:

```bash
docit
```

This will generate static HTML files in the `dist` directory.