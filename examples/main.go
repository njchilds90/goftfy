//go:build ignore
// +build ignore

// Example program demonstrating goftfy usage.
// Run with: go run examples/main.go
package main

import (
	"fmt"

	"github.com/njchilds90/goftfy"
)

func main() {
	// Basic fix
	fmt.Println("=== Basic Fix ===")
	broken := "SÃ£o Paulo is a cafÃ© city"
	fixed := goftfy.Fix(broken)
	fmt.Printf("Before: %s\n", broken)
	fmt.Printf("After:  %s\n\n", fixed)

	// HTML entities
	fmt.Println("=== HTML Entities ===")
	html := "AT&amp;T &lt;3 &quot;Go&quot;"
	fmt.Printf("Before: %s\n", html)
	fmt.Printf("After:  %s\n\n", goftfy.Fix(html))

	// Batch fix a slice (useful for AI agents processing scraped data)
	fmt.Println("=== Fix Slice (AI Agent batch processing) ===")
	docs := []string{
		"rÃ©sumÃ© for review",
		"clean text, no changes",
		"naÃ¯ve approach",
	}
	fixed2 := goftfy.FixSlice(docs)
	for i, doc := range fixed2 {
		fmt.Printf("[%d] %s\n", i, doc)
	}
	fmt.Println()

	// Explain what changed
	fmt.Println("=== Explain ===")
	original := "cafÃ©"
	result := goftfy.Fix(original)
	fmt.Println(goftfy.Explain(original, result))
	fmt.Println()

	// Check validity
	fmt.Println("=== IsValid ===")
	fmt.Printf("Is 'Hello world' valid? %v\n", goftfy.IsValid("Hello world"))
	fmt.Printf("Is 'SÃ£o' valid?        %v\n", goftfy.IsValid("SÃ£o"))
	fmt.Println()

	// Custom options
	fmt.Println("=== Custom Options (straighten quotes) ===")
	opts := goftfy.DefaultOptions()
	opts.FixCurlyQuotes = true
	quoted := "\u201CHello, World!\u201D"
	fmt.Printf("Before: %s\n", quoted)
	fmt.Printf("After:  %s\n\n", goftfy.FixWithOptions(quoted, opts))

	// Map fix — great for struct fields from DB or API
	fmt.Println("=== Fix Map ===")
	record := map[string]string{
		"city":    "SÃ£o Paulo",
		"country": "Brasil",
		"note":    "AT&amp;T office",
	}
	cleaned := goftfy.FixMap(record)
	for k, v := range cleaned {
		fmt.Printf("  %s: %s\n", k, v)
	}
}
