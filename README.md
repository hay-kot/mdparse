# mdparse

Bidirectional converter between Markdown (with YAML frontmatter) and JSON.

## Features

- **Auto-detection**: Automatically detects input format (Markdown or JSON)
- **Bidirectional**: Convert Markdown → JSON and JSON → Markdown
- **YAML Frontmatter**: Preserves frontmatter as JSON properties
- **Flexible Input**: Accepts stdin, file paths, or literal content
- **Customizable**: Configure body property name and JSON formatting

## Installation

```bash
go install github.com/hay-kot/mdparse@latest
```

## Usage

```bash
# Markdown to JSON from stdin
echo "# Hello" | mdparse

# Markdown file to JSON with pretty printing
mdparse document.md --pretty

# JSON to markdown (round-trip)
mdparse document.md | mdparse

# Custom body property name
mdparse --body-key content document.md

# Literal markdown to JSON
mdparse "# Title"
```

### Options

- `--body-key, -b`: JSON property name for markdown body content (default: "$body")
- `--pretty, -p`: Pretty-print JSON output with indentation
- `--log-level`: Set log level (debug, info, warn, error, fatal, panic)
- `--log-file`: Path to log file (optional)

## Examples

### Markdown to JSON

Input (`document.md`):
```markdown
---
title: My Document
author: John Doe
---

# Introduction

This is the content.
```

Output:
```bash
mdparse document.md --pretty
```

```json
{
  "title": "My Document",
  "author": "John Doe",
  "$body": "# Introduction\n\nThis is the content."
}
```

### JSON to Markdown

Input:
```json
{
  "title": "My Document",
  "author": "John Doe",
  "$body": "# Introduction\n\nThis is the content."
}
```

Output:
```bash
echo '{"title":"My Document","author":"John Doe","$body":"# Introduction"}' | mdparse
```

```markdown
---
author: John Doe
title: My Document
---

# Introduction
```

## License

MIT
