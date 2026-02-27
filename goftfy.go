// Package goftfy fixes broken Unicode text (mojibake, garbled characters,
// HTML entities, and other encoding artifacts). It is a Go port of the
// Python ftfy library by Robyn Speer.
//
// Basic usage:
//
//	fixed := goftfy.Fix("SÃ£o Paulo")
//	// Returns: "São Paulo"
//
// For more control use FixOptions:
//
//	fixed := goftfy.FixWithOptions("text", goftfy.Options{
//	    FixEncoding: true,
//	    FixHTMLEntities: true,
//	    FixLineBreaks: true,
//	})
package goftfy

import (
	"html"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/unicode/norm"
)

// Options controls which fixes are applied.
type Options struct {
	// FixEncoding fixes mojibake (UTF-8 text misread as Latin-1, etc.)
	FixEncoding bool
	// FixHTMLEntities decodes HTML entities like &amp; &lt; &#8217; etc.
	FixHTMLEntities bool
	// FixLineBreaks normalizes line endings to \n
	FixLineBreaks bool
	// FixSurrogates removes unpaired UTF-16 surrogates
	FixSurrogates bool
	// FixControlChars removes or replaces C0/C1 control characters
	FixControlChars bool
	// FixCurlyQuotes optionally straightens curly quotes to ASCII
	FixCurlyQuotes bool
	// NormalizationForm applies Unicode normalization (NFC, NFD, NFKC, NFKD) or "" for none
	NormalizationForm string
	// RemoveTerminalEscapes strips ANSI escape sequences
	RemoveTerminalEscapes bool
}

// DefaultOptions returns the recommended default options (mirrors ftfy defaults).
func DefaultOptions() Options {
	return Options{
		FixEncoding:           true,
		FixHTMLEntities:       true,
		FixLineBreaks:         true,
		FixSurrogates:         true,
		FixControlChars:       true,
		FixCurlyQuotes:        false,
		NormalizationForm:     "NFC",
		RemoveTerminalEscapes: false,
	}
}

// Fix applies all default fixes to the input string and returns the corrected text.
func Fix(text string) string {
	return FixWithOptions(text, DefaultOptions())
}

// FixWithOptions applies only the selected fixes from opts.
func FixWithOptions(text string, opts Options) string {
	if opts.RemoveTerminalEscapes {
		text = removeTerminalEscapes(text)
	}
	if opts.FixSurrogates {
		text = fixSurrogates(text)
	}
	if opts.FixEncoding {
		text = fixEncoding(text)
	}
	if opts.FixHTMLEntities {
		text = fixHTMLEntities(text)
	}
	if opts.FixLineBreaks {
		text = fixLineBreaks(text)
	}
	if opts.FixControlChars {
		text = fixControlChars(text)
	}
	if opts.FixCurlyQuotes {
		text = fixCurlyQuotes(text)
	}
	if opts.NormalizationForm != "" {
		text = normalize(text, opts.NormalizationForm)
	}
	return text
}

// Explain returns a human-readable description of what fixes were applied.
//
// Note: Explain() does not accept Options, so it infers applied stages by
// replaying the default pipeline in order and recording which stages changed
// the text. If the provided "fixed" string does not match what the default
// pipeline would produce, the explanation is marked as inferred.
func Explain(original, fixed string) string {
	if original == fixed {
		return "No changes needed."
	}

	// Infer changes in the same order as FixWithOptions uses.
	text := original
	var notes []string

	stage := func(name string, fn func(string) string) {
		newText := fn(text)
		if newText != text {
			notes = append(notes, name)
			text = newText
		}
	}

	opts := DefaultOptions()
	if opts.RemoveTerminalEscapes {
		stage("removed terminal escapes", removeTerminalEscapes)
	}
	if opts.FixSurrogates {
		stage("fixed surrogates", fixSurrogates)
	}
	if opts.FixEncoding {
		stage("fixed mojibake encoding", fixEncoding)
	}
	if opts.FixHTMLEntities {
		stage("decoded HTML entities", fixHTMLEntities)
	}
	if opts.FixLineBreaks {
		stage("normalized line breaks", fixLineBreaks)
	}
	if opts.FixControlChars {
		stage("removed control characters", fixControlChars)
	}
	if opts.FixCurlyQuotes {
		stage("straightened curly quotes", fixCurlyQuotes)
	}
	if opts.NormalizationForm != "" {
		stage("normalized unicode", func(s string) string { return normalize(s, opts.NormalizationForm) })
	}

	if len(notes) == 0 {
		return "Fixes applied: (unable to infer stages)."
	}

	msg := "Fixes applied: " + strings.Join(notes, ", ") + "."
	if text != fixed {
		msg += " (inferred using default options; provided fixed output differs)"
	}
	return msg
}

// IsValid reports whether the string is clean valid text needing no fixes.
func IsValid(text string) bool {
	return Fix(text) == text
}

// FixLines fixes each line of a multi-line string independently.
func FixLines(text string) string {
	lines := strings.Split(text, "\n")
	for i, line := range lines {
		lines[i] = Fix(line)
	}
	return strings.Join(lines, "\n")
}

// FixSlice fixes every string in a slice.
func FixSlice(texts []string) []string {
	result := make([]string, len(texts))
	for i, t := range texts {
		result[i] = Fix(t)
	}
	return result
}

// FixMap fixes every value in a map[string]string.
func FixMap(m map[string]string) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = Fix(v)
	}
	return result
}

// CountProblems returns the number of characters that appear to be encoding artifacts.
func CountProblems(text string) int {
	fixed := Fix(text)
	if text == fixed {
		return 0
	}
	// Approximate: count UTF-8 rune length difference. Never return negative.
	diff := utf8.RuneCountInString(text) - utf8.RuneCountInString(fixed)
	if diff < 0 {
		return 0
	}
	return diff
}

// normalize applies Unicode normalization (NFC, NFD, NFKC, NFKD).
func normalize(text, form string) string {
	switch strings.ToUpper(strings.TrimSpace(form)) {
	case "NFC":
		return norm.NFC.String(text)
	case "NFD":
		return norm.NFD.String(text)
	case "NFKC":
		return norm.NFKC.String(text)
	case "NFKD":
		return norm.NFKD.String(text)
	default:
		return text
	}
}

// ansiEscape matches ANSI terminal escape sequences.
var ansiEscape = regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]|\x1b[^[\\]`)

func removeTerminalEscapes(text string) string {
	return ansiEscape.ReplaceAllString(text, "")
}

func fixHTMLEntities(text string) string {
	// Only decode if it looks like HTML entities are present
	if !strings.Contains(text, "&") {
		return text
	}
	return html.UnescapeString(text)
}

func fixLineBreaks(text string) string {
	// Normalize \r\n and \r to \n
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")
	// Replace Unicode line/paragraph separators
	text = strings.ReplaceAll(text, "\u2028", "\n")
	text = strings.ReplaceAll(text, "\u2029", "\n")
	return text
}

func fixSurrogates(text string) string {
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		if r >= 0xD800 && r <= 0xDFFF {
			// unpaired surrogate — replace with replacement char
			b.WriteRune(unicode.ReplacementChar)
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func fixControlChars(text string) string {
	var b strings.Builder
	b.Grow(len(text))
	for _, r := range text {
		// Allow tab, newline, carriage return; strip other C0 and all C1 controls
		if r == '\t' || r == '\n' || r == '\r' {
			b.WriteRune(r)
		} else if r < 0x20 || (r >= 0x7F && r <= 0x9F) {
			// strip
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

var curlyQuoteReplacer = strings.NewReplacer(
	"\u2018", "'", // left single quotation mark
	"\u2019", "'", // right single quotation mark
	"\u201A", "'", // single low-9 quotation mark
	"\u201B", "'", // single high-reversed-9 quotation mark
	"\u201C", `"`, // left double quotation mark
	"\u201D", `"`, // right double quotation mark
	"\u201E", `"`, // double low-9 quotation mark
	"\u201F", `"`, // double high-reversed-9 quotation mark
	"\u2039", "<", // single left-pointing angle quotation mark
	"\u203A", ">", // single right-pointing angle quotation mark
	"\u00AB", `"`, // left-pointing double angle quotation mark
	"\u00BB", `"`, // right-pointing double angle quotation mark
)

func fixCurlyQuotes(text string) string {
	return curlyQuoteReplacer.Replace(text)
}