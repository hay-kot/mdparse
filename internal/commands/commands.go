package commands

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v3"

	"github.com/hay-kot/mdparse/internal/parser"
)

// ParseCmd implements the parse command
type ParseCmd struct {
	flags      *Flags
	parseFlags *ParseFlags
}

// NewParseCmd creates a new parse command
func NewParseCmd(flags *Flags, parseFlags *ParseFlags) *ParseCmd {
	return &ParseCmd{
		flags:      flags,
		parseFlags: parseFlags,
	}
}

// Run executes the parse command
func (cmd *ParseCmd) Run(ctx context.Context, c *cli.Command) error {
	source, err := cmd.readInput(c)
	if err != nil {
		return err
	}

	// Detect format and convert
	var output []byte
	if parser.IsJSONObject(source) {
		// JSON to Markdown
		output, err = parser.JSONToMarkdown(source, cmd.parseFlags.BodyKey)
		if err != nil {
			return fmt.Errorf("failed to convert JSON to markdown: %w", err)
		}
	} else {
		// Markdown to JSON
		output, err = parser.MarkdownToJSON(source, cmd.parseFlags.BodyKey, cmd.parseFlags.Pretty)
		if err != nil {
			return fmt.Errorf("failed to convert markdown to JSON: %w", err)
		}
	}

	if !bytes.HasSuffix(output, []byte("\n")) {
		output = append(output, '\n')
	}

	fmt.Print(string(output))
	return nil
}

func (cmd *ParseCmd) readInput(c *cli.Command) ([]byte, error) {
	if c.Args().Len() > 0 {
		arg := c.Args().First()

		// Check if argument looks like a file path
		if isLikelyPath(arg) {
			if _, err := os.Stat(arg); err == nil {
				data, err := os.ReadFile(arg)
				if err != nil {
					return nil, fmt.Errorf("failed to read file: %w", err)
				}
				return data, nil
			}
			// Looked like a path but doesn't exist
			return nil, fmt.Errorf("file not found: %s", arg)
		}

		// Treat as literal content (markdown or JSON)
		return []byte(arg), nil
	}

	// Read from stdin
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdin: %w", err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("no input provided")
	}

	return data, nil
}
