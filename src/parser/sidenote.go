package parser

import (
	"fmt"
	"html"
	"regexp"
	"strings"
	"sync/atomic"
)

func error(msg string) {
	fmt.Errorf("Error at %s", msg)
}

type TokenType int

const (
	leftMeta = iota
	rightMeta
	sidenote
	sidenoteBlock
	chartBlock
	literal
	toc
	eof
)

type Token struct {
	typ TokenType
	val string
}

type Lexer struct {
	start   int
	current int
	input   string
	tokens  []Token
}

func (l *Lexer) isAtEnd() bool {
	return l.current > len(l.input)
}

func (l *Lexer) scanTokens() []Token {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}

	l.tokens = append(l.tokens, Token{typ: eof})
	return l.tokens
}

func (l *Lexer) scanToken() {
	c := l.advance()
	switch c {
	case '{':
		if l.accept('%') {
			l.tokens = append(l.tokens, Token{typ: leftMeta, val: l.value()})
		}
	default:
		error("Unexpected character")
	}
}

func (l *Lexer) advance() byte {
	return l.input[l.current]
}

func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}

	return l.input[l.current]
}

func (l *Lexer) accept(expected byte) bool {
	if l.isAtEnd() {
		return false
	}
	if l.input[l.current] != expected {
		return false
	}

	l.current++
	return true
}

func (l *Lexer) value() string {
	return l.input[l.start:l.current]
}

func newLexer(input string) *Lexer {
	return &Lexer{
		start:   0,
		current: 0,
		input:   input,
		tokens:  make([]Token, 0),
	}
}

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

func postprocessSidenotes(content []byte) []byte {
	// Handle {sidenote: abc} and {sidenote} blocks
	parsedBytes := make([]byte, 0)

	inSidenote := false
	inCodeblock := false

	sidenoteOpener := "{sidenote"

	for i, b := range content {
		// Check if we get a sidenote
		if b == '{' && i+len(sidenoteOpener)-1 < len(content) && string(content[i:i+len(sidenoteOpener)]) == sidenoteOpener {
			if inSidenote {
				// Already in a sidenote, just append the character
				parsedBytes = append(parsedBytes, b)
				continue
			}
			inSidenote = true
			continue
		}
		i += 20
	}
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

	markerClasses := "sidenote-marker"
	sidenoteClasses := "sidenote float-right clear-right ml-[50px] mr-[-300px] w-[250px] -mt-12 mb-2 border-l-4 border-primary p-2 text-xs transition-all duration-300 ease-in-out hover:bg-base-300"

	// Process sidenotes in reverse order so positions don't shift
	for i := len(sidenoteMatches) - 1; i >= 0; i-- {
		match := sidenoteMatches[i]

		var replacement []byte
		if match.IsInline {
			marker := fmt.Sprintf(`<span class="%s" data-sidenote-id="%d">%d</span>`, markerClasses, match.ID, match.ID)
			sidenote := fmt.Sprintf(`<aside class="%s" id="sidenote-%d">%s</aside>`, sidenoteClasses, match.ID, html.EscapeString(match.Content))

			replacement = []byte(marker + sidenote)
		} else {
			sidenote := fmt.Sprintf(`<aside class="%s" id="sidenote-%d">%s</aside>`, sidenoteClasses, match.ID, html.EscapeString(match.Content))
			replacement = []byte(sidenote)
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
