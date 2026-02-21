# goftfy — Fixes Text For You (Go)

[![Go Reference](https://pkg.go.dev/badge/github.com/njchilds90/goftfy.svg)](https://pkg.go.dev/github.com/njchilds90/goftfy)
[![Go Report Card](https://goreportcard.com/badge/github.com/njchilds90/goftfy)](https://goreportcard.com/report/github.com/njchilds90/goftfy)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**goftfy** is a Go port of the popular Python [`ftfy`](https://github.com/rspeer/python-ftfy) library by Robyn Speer. It fixes broken Unicode text — mojibake, garbled characters, HTML entities in the wrong place, control characters, and other text encoding artifacts.

No external dependencies. Zero allocations for already-clean text.

---

## Why goftfy?

When you scrape websites, process legacy databases, consume third-party APIs, or feed text through LLM pipelines, you inevitably encounter garbled text:

| Broken | Fixed |
|--------|-------|
| `SÃ£o Paulo` | `São Paulo` |
| `cafÃ©` | `café` |
| `rÃ©sumÃ©` | `résumé` |
| `AT&amp;T` | `AT&T` |
| `naÃ¯ve` | `naïve` |
| `â€™` | `'` (right single quote) |

Python's `ftfy` is the gold standard for this in Python. **goftfy** brings the same power to Go with zero dependencies.

---

## Installation
```bash
go get github.com/njchilds90/goftfy
```

---

## Quick Start
```go
package main

import (
    "fmt"
    "github.com/njchilds90/goftfy"
)

func main() {
    fmt.Println(goftfy.Fix("SÃ£o Paulo"))  // São Paulo
    fmt.Println(goftfy.Fix("AT&amp;T"))    // AT&T
    fmt.Println(goftfy.Fix("cafÃ©"))       // café
}
```

---

## API Reference

### Core
```go
// Fix applies all default fixes.
goftfy.Fix(text string) string

// FixWithOptions applies only the specified fixes.
goftfy.FixWithOptions(text string, opts Options) string

// DefaultOptions returns the recommended option set.
goftfy.DefaultOptions() Options
```

### Batch
```go
// FixLines fixes each line independently.
goftfy.FixLines(text string) string

// FixSlice fixes every string in a slice.
goftfy.FixSlice(texts []string) []string

// FixMap fixes every value in a map[string]string.
goftfy.FixMap(m map[string]string) map[string]string
```

### Analysis
```go
// IsValid reports whether text needs no fixing.
goftfy.IsValid(text string) bool

// Explain returns a human-readable summary of what was fixed.
goftfy.Explain(original, fixed string) string

// CountProblems estimates the number of encoding artifacts.
goftfy.CountProblems(text string) int

// AnalyzeString returns per-character diagnostic info.
goftfy.AnalyzeString(text string) []CharInfo

// HasReplacementChars checks for U+FFFD.
goftfy.HasReplacementChars(text string) bool

// HasSurrogates checks for unpaired UTF-16 surrogates.
goftfy.HasSurrogates(text string) bool
```

### Quick utilities
```go
// QuickFix uses a fast pattern dictionary for common mojibake.
goftfy.QuickFix(text string) string

// CommonMojibakePatterns returns the built-in pattern map.
goftfy.CommonMojibakePatterns() map[string]string
```

---

## Options
```go
opts := goftfy.Options{
    FixEncoding:           true,   // Fix mojibake (UTF-8 read as Latin-1)
    FixHTMLEntities:       true,   // Decode &amp; &lt; &#8217; etc.
    FixLineBreaks:         true,   // Normalize \r\n, \r → \n
    FixSurrogates:         true,   // Replace unpaired surrogates with U+FFFD
    FixControlChars:       true,   // Strip C0/C1 control chars
    FixCurlyQuotes:        false,  // Straighten " " ' ' → " "  ' '
    NormalizationForm:     "NFC",  // Unicode normalization (or "")
    RemoveTerminalEscapes: false,  // Strip ANSI escape codes
}
```

---

## Use Cases

**Data pipelines / ETL** — Fix text fields before inserting into a database.

**AI agent pipelines** — Clean scraped web text before sending to an LLM.

**API response processing** — Fix encoding artifacts from third-party APIs.

**Log processing** — Strip ANSI escape codes and control characters from logs.

**CSV/spreadsheet import** — Fix encoding artifacts in imported data.

---

## Comparison with Python ftfy

| Feature | Python ftfy | goftfy |
|---------|-------------|--------|
| Mojibake (UTF-8 → Latin-1) | ✅ | ✅ |
| HTML entity decoding | ✅ | ✅ |
| Line break normalization | ✅ | ✅ |
| Control character removal | ✅ | ✅ |
| Curly quote straightening | ✅ | ✅ |
| ANSI escape removal | ✅ | ✅ |
| Unicode normalization | ✅ | ✅ (NFC best-effort) |
| Zero dependencies | ✅ | ✅ |
| Batch processing helpers | ❌ | ✅ |

---

## Contributing

Contributions welcome! Please open an issue or pull request.

1. Fork the repository
2. Create your feature branch: `git checkout -b feature/my-feature`
3. Add tests for new functionality
4. Submit a pull request

---

## License

MIT — see [LICENSE](LICENSE)

---

## Credits

Inspired by [ftfy](https://github.com/rspeer/python-ftfy) by Robyn Speer.
