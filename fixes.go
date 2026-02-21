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
	// Common mojibake signatures: Ã followed by a character in range 0x80-0xBF
	for i, r := range text {
		if r == 'Ã' && i+1 < len(text) {
			next := text[i+1]
			if next >= 0x80 && next <= 0xBF {
				return true
			}
			// UTF-8 multi-byte second byte showing as visible char
			if next == '©' || next == '®' || next == '™' ||
				next == '°' || next == '±' || next == '²' ||
				next == '³' || next == '¼' || next == '½' {
				return true
			}
		}
		// Detect Windows-1252 read as Latin-1
		if r == 'â' {
			return true
		}
	}
	return false
}

// decodeMojibake reverses Latin-1 misinterpretation of UTF-8.
// This reinterprets each rune as its Latin-1 byte value and re-decodes as UTF-8.
func decodeMojibake(text string) string {
	// Convert string to raw Latin-1 bytes
	var rawBytes []byte
	for _, r := range text {
		if r < 0x100 {
			rawBytes = append(rawBytes, byte(r))
		} else {
			// Not a Latin-1 character; encode as UTF-8
			buf := make([]byte, utf8.UTFMax)
			n := utf8.EncodeRune(buf, r)
			rawBytes = append(rawBytes, buf[:n]...)
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

// CommonMojibakePatterns returns a map of common mojibake sequences to their correct UTF-8.
// Useful for quick lookups or educational purposes.
func CommonMojibakePatterns() map[string]string {
	return map[string]string{
		"SÃ£o":    "São",
		"clichÃ©": "cliché",
		"cafÃ©":   "café",
		"rÃ©sumÃ©": "résumé",
		"naÃ¯ve":  "naïve",
		"â€™":     "\u2019", // right single quote
		"â€œ":     "\u201C", // left double quote
		"â€":      "\u201D", // right double quote
		"â€"":     "\u2014", // em dash
		"â€"":     "\u2013", // en dash
		"Â·":      "·",
		"Â©":      "©",
		"Â®":      "®",
		"â„¢":     "™",
	}
}

// QuickFix applies a fast dictionary lookup for the most common mojibake patterns.
// Faster than the full Fix() for known patterns but less comprehensive.
func QuickFix(text string) string {
	for broken, fixed := range CommonMojibakePatterns() {
		text = strings.ReplaceAll(text, broken, fixed)
	}
	return text
}
