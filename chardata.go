package goftfy

// CharInfo holds information about a Unicode character's context.
type CharInfo struct {
	Rune        rune
	Category    string
	IsProblematic bool
	Suggestion  rune
}

// AnalyzeString returns per-character analysis of potentially problematic chars.
func AnalyzeString(text string) []CharInfo {
	var result []CharInfo
	for _, r := range text {
		info := analyzeRune(r)
		if info.IsProblematic {
			result = append(result, info)
		}
	}
	return result
}

func analyzeRune(r rune) CharInfo {
	info := CharInfo{Rune: r}
	switch {
	case r >= 0xD800 && r <= 0xDFFF:
		info.Category = "surrogate"
		info.IsProblematic = true
		info.Suggestion = '\uFFFD'
	case r >= 0x00 && r < 0x20 && r != '\t' && r != '\n' && r != '\r':
		info.Category = "control_C0"
		info.IsProblematic = true
	case r >= 0x7F && r <= 0x9F:
		info.Category = "control_C1"
		info.IsProblematic = true
	case r == '\uFFFD':
		info.Category = "replacement_char"
		info.IsProblematic = true
	case isMojibakeChar(r):
		info.Category = "likely_mojibake"
		info.IsProblematic = true
	default:
		info.Category = "ok"
	}
	return info
}

func isMojibakeChar(r rune) bool {
	// Characters that appear frequently in UTF-8 read as Latin-1
	mojibakeIndicators := []rune{
		'Ã', 'â', 'Â', 'ï', 'Å', 'Ä', 'Ö', 'Ü',
	}
	for _, m := range mojibakeIndicators {
		if r == m {
			return true
		}
	}
	return false
}

// HasReplacementChars reports whether the string contains the Unicode replacement character.
func HasReplacementChars(text string) bool {
	return strings.Contains(text, "\uFFFD")
}

// HasSurrogates reports whether the string contains unpaired UTF-16 surrogates.
func HasSurrogates(text string) bool {
	for _, r := range text {
		if r >= 0xD800 && r <= 0xDFFF {
			return true
		}
	}
	return false
}
