# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

## [v0.1.0] - Initial Release

### Added
- `Fix()` — one-call fix with sensible defaults
- `FixWithOptions()` — fine-grained control over which fixes apply
- `FixLines()`, `FixSlice()`, `FixMap()` — batch helpers
- Mojibake detection and reversal (UTF-8 misread as Latin-1)
- HTML entity decoding
- Line break normalization
- Control character removal
- Curly quote straightening (opt-in)
- ANSI terminal escape removal (opt-in)
- Surrogate character handling
- `IsValid()`, `Explain()`, `CountProblems()`, `AnalyzeString()`
- `QuickFix()` with common pattern dictionary
- Zero external dependencies
