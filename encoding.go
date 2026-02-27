package goftfy

import (
	"strings"
	"unicode/utf8"
)

// fixEncoding is the core mojibake fixer.
// Mojibake happens when UTF-8 bytes are decoded as Latin-1 (ISO-8859-1)
// and then re-encoded. We detect and reverse this.
func fixEncoding(text string) string {
	if utf8.ValidString(text) && !looksLikeMojibake(text) {
		return text
	}
	// Try to recover UTF-8 from Latin-1 mojibake
	result := decodeMojibake(text)
	if result != text && utf8.ValidString(result) {
		return result
	}
	return text
}

// looksLikeMojibake uses heuristics to detect common mojibake patterns.
func looksLikeMojibake(text string) bool {
	// Operate on runes (not bytes). The previous implementation mixed byte indexes
	// from range with text[i+1], which is unsafe for non-ASCII.
	rs := []rune(text)
	for i := 0; i < len(rs); i++ {
		r := rs[i]

		// Common UTF-8->Latin-1 mojibake signature: "Ã" then a rune in U+0080..U+00BF
		// (often shows up as "Ã©", "Ã±", "Ã£", etc.).
		if r == 'Ã' && i+1 < len(rs) {
			next := rs[i+1]
			if next >= 0x80 && next <= 0xBF {
				return true
			}
			switch next {
			case '©', '®', '™', '°', '±', '²', '³', '¼', '½':
				return true
			}
		}

		// Common Windows-1252 mojibake sequences often start with â / Â.
		if r == 'â' || r == 'Â' {
			return true
		}
	}
	return false
}

// decodeMojibake reverses Latin-1 misinterpretation of UTF-8.
// This reinterprets each rune as its Latin-1 byte value and re-decodes as UTF-8.
func decodeMojibake(text string) string {
	// Convert string to raw Latin-1 bytes.
	// Pre-size to byte length as a reasonable upper bound for most mojibake strings.
	rawBytes := make([]byte, 0, len(text))
	for _, r := range text {
		if r < 0x100 {
			rawBytes = append(rawBytes, byte(r))
		} else {
			// Not a Latin-1 character; append its UTF-8 encoding.
			rawBytes = utf8.AppendRune(rawBytes, r)
		}
	}

	if utf8.Valid(rawBytes) {
		candidate := string(rawBytes)
		// Make sure we actually improved things
		if countNonASCII(candidate) < countNonASCII(text) {
			return candidate
		}
	}
	return text
}

func countNonASCII(s string) int {
	count := 0
	for _, r := range s {
		if r > 127 {
			count++
		}
	}
	return count
}

// commonMojibakePatternsOrdered is the deterministic replacement order for QuickFix.
var commonMojibakePatternsOrdered = []struct{ broken, fixed string }{
	{"SÃ£o", "São"},
	{"cafÃ©", "café"},
	{"clichÃ©", "cliché"},
	{"rÃ©sumÃ©", "résumé"},
	{"naÃ¯ve", "naïve"},

	// Common Windows-1252 punctuation mojibake
	{"â€™", "\u2019"}, // right single quotation mark
	{"â€˜", "\u2018"}, // left single quotation mark
	{"â€œ", "\u201C"}, // left double quotation mark
	{"â€�", "\u201D"}, // right double quotation mark
	{"â€”", "\u2014"}, // em dash
	{"â€“", "\u2013"}, // en dash
	{"â€¦", "\u2026"}, // ellipsis

	// Misc
	{"Â·", "·"},
	{"Â©", "©"},
	{"Â®", "®"},
	{"â„¢", "™"},
}

var commonMojibakePatternsMap = func() map[string]string {
	m := make(map[string]string, len(commonMojibakePatternsOrdered))
	for _, p := range commonMojibakePatternsOrdered {
		m[p.broken] = p.fixed
	}
	return m
}()

// CommonMojibakePatterns returns a map of common mojibake sequences to their correct UTF-8.
// Useful for quick lookups or educational purposes.
//
// The returned map is a copy to prevent callers from mutating package state.
func CommonMojibakePatterns() map[string]string {
	out := make(map[string]string, len(commonMojibakePatternsMap))
	for k, v := range commonMojibakePatternsMap {
		out[k] = v
	}
	return out
}

// QuickFix applies a fast dictionary lookup for the most common mojibake patterns.
// Faster than the full Fix() for known patterns but less comprehensive.
func QuickFix(text string) string {
	for _, p := range commonMojibakePatternsOrdered {
		text = strings.ReplaceAll(text, p.broken, p.fixed)
	}
	return text
}