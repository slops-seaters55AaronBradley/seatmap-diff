package diff

import (
	"fmt"
	"io"
	"strings"
)

// OutputFormat defines the output format for diff results.
type OutputFormat string

const (
	// FormatText outputs a human-readable text diff.
	FormatText OutputFormat = "text"
	// FormatJSON outputs a machine-readable JSON diff.
	FormatJSON OutputFormat = "json"
	// FormatMarkdown outputs a Markdown-formatted diff.
	FormatMarkdown OutputFormat = "markdown"
)

// ANSI color codes for terminal output.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Formatter writes a Result to an io.Writer in the specified format.
type Formatter struct {
	Format    OutputFormat
	Colorized bool
	Writer    io.Writer
}

// NewFormatter creates a new Formatter with the given options.
func NewFormatter(w io.Writer, format OutputFormat, colorized bool) *Formatter {
	return &Formatter{
		Format:    format,
		Colorized: colorized,
		Writer:    w,
	}
}

// Write outputs the diff result using the configured format.
func (f *Formatter) Write(result *Result) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(result)
	case FormatMarkdown:
		return f.writeMarkdown(result)
	default:
		return f.writeText(result)
	}
}

// writeText outputs a human-readable, optionally colorized diff.
func (f *Formatter) writeText(result *Result) error {
	if len(result.Changes) == 0 {
		_, err := fmt.Fprintln(f.Writer, "No differences found.")
		return err
	}

	for _, change := range result.Changes {
		var line string
		switch change.Type {
		case ChangeAdded:
			marker := f.colorize("+", colorGreen)
			line = fmt.Sprintf("%s [added]   %s: %v", marker, change.Path, change.NewValue)
		case ChangeRemoved:
			marker := f.colorize("-", colorRed)
			line = fmt.Sprintf("%s [removed] %s: %v", marker, change.Path, change.OldValue)
		case ChangeModified:
			marker := f.colorize("~", colorYellow)
			line = fmt.Sprintf("%s [changed] %s: %v -> %v", marker, change.Path, change.OldValue, change.NewValue)
		}
		if _, err := fmt.Fprintln(f.Writer, line); err != nil {
			return err
		}
	}

	summary := fmt.Sprintf("\nSummary: %d added, %d removed, %d modified",
		result.Stats.Added, result.Stats.Removed, result.Stats.Modified)
	_, err := fmt.Fprintln(f.Writer, f.colorize(summary, colorBold))
	return err
}

// writeMarkdown outputs a Markdown table representation of the diff.
func (f *Formatter) writeMarkdown(result *Result) error {
	if len(result.Changes) == 0 {
		_, err := fmt.Fprintln(f.Writer, "_No differences found._")
		return err
	}

	lines := []string{
		"| Type | Path | Old Value | New Value |",
		"|------|------|-----------|-----------|" ,
	}
	for _, change := range result.Changes {
		row := fmt.Sprintf("| %s | `%s` | %v | %v |",
			string(change.Type), change.Path, change.OldValue, change.NewValue)
		lines = append(lines, row)
	}
	_, err := fmt.Fprintln(f.Writer, strings.Join(lines, "\n"))
	return err
}

// writeJSON outputs the result as indented JSON.
func (f *Formatter) writeJSON(result *Result) error {
	data, err := result.MarshalJSON()
	if err != nil {
		return fmt.Errorf("formatting result as JSON: %w", err)
	}
	_, err = fmt.Fprintln(f.Writer, string(data))
	return err
}

// colorize wraps text in an ANSI color code if colorization is enabled.
func (f *Formatter) colorize(text, color string) string {
	if !f.Colorized {
		return text
	}
	return color + text + colorReset
}
