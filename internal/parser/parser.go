package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

// MarkdownToJSON converts markdown with frontmatter to JSON
// bodyKey specifies the JSON property name for the markdown body content
// pretty determines if the JSON should be formatted with indentation
func MarkdownToJSON(markdown []byte, bodyKey string, pretty bool) ([]byte, error) {
	frontmatter, body, err := parseFrontmatter(markdown)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Create output map with frontmatter fields
	output := make(map[string]any)
	for k, v := range frontmatter {
		output[k] = v
	}

	// Add body with custom key
	output[bodyKey] = body

	// Marshal to JSON
	if pretty {
		return json.MarshalIndent(output, "", "  ")
	}
	return json.Marshal(output)
}

// parseFrontmatter extracts YAML frontmatter and body from markdown
func parseFrontmatter(markdown []byte) (map[string]any, string, error) {
	content := string(markdown)
	lines := strings.Split(content, "\n")

	// Check if content starts with frontmatter delimiter
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		// No frontmatter, entire content is body
		return make(map[string]any), strings.TrimSpace(content), nil
	}

	// Find closing delimiter
	closingIndex := -1
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			closingIndex = i
			break
		}
	}

	if closingIndex == -1 {
		// No closing delimiter found, treat as no frontmatter
		return make(map[string]any), strings.TrimSpace(content), nil
	}

	// Extract frontmatter YAML
	frontmatterLines := lines[1:closingIndex]
	frontmatterYAML := strings.Join(frontmatterLines, "\n")

	// Parse YAML frontmatter
	var frontmatter map[string]any
	if len(strings.TrimSpace(frontmatterYAML)) > 0 {
		if err := yaml.Unmarshal([]byte(frontmatterYAML), &frontmatter); err != nil {
			return nil, "", fmt.Errorf("failed to parse YAML frontmatter: %w", err)
		}
	} else {
		frontmatter = make(map[string]any)
	}

	// Extract body (everything after closing delimiter)
	var bodyLines []string
	if closingIndex+1 < len(lines) {
		bodyLines = lines[closingIndex+1:]
	}
	body := strings.TrimSpace(strings.Join(bodyLines, "\n"))

	return frontmatter, body, nil
}

// JSONToMarkdown converts JSON back to markdown with frontmatter
// bodyKey specifies which JSON property contains the markdown body
func JSONToMarkdown(jsonData []byte, bodyKey string) ([]byte, error) {
	var data map[string]any
	if err := json.Unmarshal(jsonData, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Extract body content
	body, ok := data[bodyKey]
	if !ok {
		return nil, fmt.Errorf("body key %q not found in JSON", bodyKey)
	}

	bodyStr, ok := body.(string)
	if !ok {
		return nil, fmt.Errorf("body value is not a string")
	}

	// Remove body from data to get frontmatter
	delete(data, bodyKey)

	var buf bytes.Buffer

	// Write frontmatter if present
	if len(data) > 0 {
		yamlData, err := yaml.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal frontmatter to YAML: %w", err)
		}

		buf.WriteString("---\n")
		buf.Write(yamlData)
		buf.WriteString("---\n\n")
	}

	// Write body
	buf.WriteString(bodyStr)

	return buf.Bytes(), nil
}

// IsJSONObject checks if the input looks like a JSON object
func IsJSONObject(data []byte) bool {
	trimmed := bytes.TrimSpace(data)
	if len(trimmed) == 0 {
		return false
	}

	// Check if it starts with {
	if trimmed[0] != '{' {
		return false
	}

	// Verify it's valid JSON by attempting to unmarshal
	var js map[string]any
	return json.Unmarshal(trimmed, &js) == nil
}
