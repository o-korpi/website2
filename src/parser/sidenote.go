package parser

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"sync/atomic"
)

// Global counter for sidenote IDs
var sidenoteCounter int64

type CodeBlockRange struct {
	Start, End int
}

type SidenoteMatch struct {
	Start, End   int
	Content      string
	IsInline     bool
	ID           int64
	OriginalText []byte
}

func findCodeBlockRanges(content []byte) []CodeBlockRange {
	var ranges []CodeBlockRange

	// Find code blocks
	patterns := []string{
		`<pre[^>]*><code[^>]*>.*?</code></pre>`, // Standard code blocks
		`<code[^>]*>.*?</code>`,                 // Inline code blocks
	}

	for _, pattern := range patterns {
		regex := regexp.MustCompile("(?s)" + pattern)
		matches := regex.FindAllIndex(content, -1)

		for _, match := range matches {
			ranges = append(ranges, CodeBlockRange{
				Start: match[0],
				End:   match[1],
			})
		}
	}

	return ranges
}

// isInCodeBlock checks if a position is within any code block range
func isInCodeBlock(pos int, codeBlocks []CodeBlockRange) bool {
	for _, block := range codeBlocks {
		if pos >= block.Start && pos <= block.End {
			return true
		}
	}
	return false
}

func ProcessSidenotes(content []byte) []byte {
	// Reset counter
	atomic.StoreInt64(&sidenoteCounter, 0)

	// Find all code blocks
	codeBlocks := findCodeBlockRanges(content)

	// Find all sidenotes
	sidenoteMatches := findAllSidenotes(content, codeBlocks)

	// Process sidenotes in reverse order so positions don't shift
	for i := len(sidenoteMatches) - 1; i >= 0; i-- {
		match := sidenoteMatches[i]

		var replacement []byte
		if match.IsInline {
			markerClasses := "sidenote-marker"
			marker := fmt.Sprintf(`<span class="%s" data-sidenote-id="%d">%d</span>`, markerClasses, match.ID, match.ID)

			sidenoteClasses := "sidenote inline float-right clear-right ml-[50px] mr-[-300px] w-[250px] mb-1 glass border border-l border-l-4 border-primary p-1 text-sm rounded shadow transition-all duration-300 ease-in-out"
			sidenote := fmt.Sprintf(`<aside class="%s" id="sidenote-%d">%s</aside>`, sidenoteClasses, match.ID, html.EscapeString(match.Content))

			replacement = []byte(marker + sidenote)
		} else {
			replacement = []byte(fmt.Sprintf(
				`<aside class="float-right clear-right w-64 ml-4 -mr-80 mb-4 bg-base-200 border-l-4 border-base-300 p-4 text-sm text-base-content rounded-lg shadow-lg transition-all duration-300 hover:bg-base-300 hover:border-primary" id="sidenote-%d">%s</aside>`,
				match.ID, html.EscapeString(match.Content)))
		}

		// Replace the match in the content
		content = append(content[:match.Start], append(replacement, content[match.End:]...)...)
	}

	return content
}

func findAllSidenotes(content []byte, codeBlocks []CodeBlockRange) []SidenoteMatch {
	var matches []SidenoteMatch

	// Find inline sidenotes
	inlineRegex := regexp.MustCompile(`\{sidenote\s+([^}]+)\}`)
	inlineMatches := inlineRegex.FindAllSubmatchIndex(content, -1)

	for _, match := range inlineMatches {
		start, end := match[0], match[1]
		contentStart, contentEnd := match[2], match[3]

		if isInCodeBlock(start, codeBlocks) {
			continue
		}

		// Generate unique ID
		id := atomic.AddInt64(&sidenoteCounter, 1)
		sidenoteContent := strings.TrimSpace(string(content[contentStart:contentEnd]))

		matches = append(matches, SidenoteMatch{
			Start:        start,
			End:          end,
			Content:      sidenoteContent,
			IsInline:     true,
			ID:           id,
			OriginalText: content[start:end],
		})

		fmt.Printf("DEBUG: Added inline sidenote %d: %s\n", id, sidenoteContent)
	}

	// Find block sidenotes
	blockRegex := regexp.MustCompile(`(?s)\{sidenote\}(.*?)\{/sidenote\}`)
	blockMatches := blockRegex.FindAllSubmatchIndex(content, -1)

	for _, match := range blockMatches {
		start, end := match[0], match[1]
		contentStart, contentEnd := match[2], match[3]

		if isInCodeBlock(start, codeBlocks) {
			continue
		}

		id := atomic.AddInt64(&sidenoteCounter, 1)
		sidenoteContent := strings.TrimSpace(string(content[contentStart:contentEnd]))

		matches = append(matches, SidenoteMatch{
			Start:        start,
			End:          end,
			Content:      sidenoteContent,
			IsInline:     false,
			ID:           id,
			OriginalText: content[start:end],
		})
	}

	return matches
}
