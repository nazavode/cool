// generated by Textmapper; DO NOT EDIT

package cool

import (
	"strings"
	"unicode/utf8"
)

// Lexer states.
const (
	StateInitial   = 0
	StateInComment = 1
)

// Lexer uses a generated DFA to scan through a utf-8 encoded input string. If
// the string starts with a BOM character, it gets skipped.
type Lexer struct {
	source string

	ch          rune // current character, -1 means EOI
	offset      int  // character offset
	tokenOffset int  // last token offset
	line        int  // current line number (1-based)
	tokenLine   int  // last token line
	scanOffset  int  // scanning offset
	value       interface{}

	State        int // lexer state, modifiable
	commentLevel int // number of open nested block comments
}

var bomSeq = "\xef\xbb\xbf"

// Init prepares the lexer l to tokenize source by performing the full reset
// of the internal state.
func (l *Lexer) Init(source string) {
	l.source = source

	l.ch = 0
	l.offset = 0
	l.tokenOffset = 0
	l.line = 1
	l.tokenLine = 1
	l.State = 0
	l.commentLevel = 0

	if strings.HasPrefix(source, bomSeq) {
		l.offset += len(bomSeq)
	}

	l.rewind(l.offset)
}

// Next finds and returns the next token in l.source. The source end is
// indicated by Token.EOI.
//
// The token text can be retrieved later by calling the Text() method.
func (l *Lexer) Next() Token {
restart:
	l.tokenLine = l.line
	l.tokenOffset = l.offset

	state := tmStateMap[l.State]
	backupRule := -1
	var backupOffset int
	for state >= 0 {
		var ch int
		if uint(l.ch) < tmRuneClassLen {
			ch = int(tmRuneClass[l.ch])
		} else if l.ch < 0 {
			state = int(tmLexerAction[state*tmNumClasses])
			continue
		} else {
			ch = 1
		}
		state = int(tmLexerAction[state*tmNumClasses+ch])
		if state > tmFirstRule {
			if state < 0 {
				state = (-1 - state) * 2
				backupRule = tmBacktracking[state]
				backupOffset = l.offset
				state = tmBacktracking[state+1]
			}
			if l.ch == '\n' {
				l.line++
			}

			// Scan the next character.
			// Note: the following code is inlined to avoid performance implications.
			l.offset = l.scanOffset
			if l.offset < len(l.source) {
				r, w := rune(l.source[l.offset]), 1
				if r >= 0x80 {
					// not ASCII
					r, w = utf8.DecodeRuneInString(l.source[l.offset:])
				}
				l.scanOffset += w
				l.ch = r
			} else {
				l.ch = -1 // EOI
			}
		}
	}

	rule := tmFirstRule - state
recovered:

	token := tmToken[rule]
	space := false
	switch rule {
	case 0:
		if backupRule >= 0 {
			rule = backupRule
			l.rewind(backupOffset)
		} else if l.offset == l.tokenOffset {
			l.rewind(l.scanOffset)
		}
		if rule != 0 {
			goto recovered
		}
	case 3: // whitespace: /[\n\r\t ]+/
		space = true
	case 4: // EnterBlockComment: /\(\*/
		space = true
		{
			l.enterBlockComment()
		}
	case 5: // invalid_token: /\*\)/
		space = true
	case 6: // ExitBlockComment: /\*\)/
		space = true
		{
			l.exitBlockComment()
		}
	case 7: // BlockComment: /([^*(]|\*+[^)]|\([^*])*/
		space = true
	case 8: // SingleLineComment: /\-\-.*/
		space = true
	}
	if space {
		goto restart
	}
	return token
}

// Pos returns the start and end positions of the last token returned by Next().
func (l *Lexer) Pos() (start, end int) {
	start = l.tokenOffset
	end = l.offset
	return
}

// Line returns the line number of the last token returned by Next().
func (l *Lexer) Line() int {
	return l.tokenLine
}

// Text returns the substring of the input corresponding to the last token.
func (l *Lexer) Text() string {
	return l.source[l.tokenOffset:l.offset]
}

// Value returns the value associated with the last returned token.
func (l *Lexer) Value() interface{} {
	return l.value
}

// rewind can be used in lexer actions to accept a portion of a scanned token, or to include
// more text into it.
func (l *Lexer) rewind(offset int) {
	if offset < l.offset {
		l.line -= strings.Count(l.source[offset:l.offset], "\n")
	} else {
		if offset > len(l.source) {
			offset = len(l.source)
		}
		l.line += strings.Count(l.source[l.offset:offset], "\n")
	}

	// Scan the next character.
	l.scanOffset = offset
	l.offset = offset
	if l.offset < len(l.source) {
		r, w := rune(l.source[l.offset]), 1
		if r >= 0x80 {
			// not ASCII
			r, w = utf8.DecodeRuneInString(l.source[l.offset:])
		}
		l.scanOffset += w
		l.ch = r
	} else {
		l.ch = -1 // EOI
	}
}

func (l *Lexer) enterBlockComment() {
	l.commentLevel++
	l.State = StateInComment
}

func (l *Lexer) exitBlockComment() {
	l.commentLevel--
	if l.commentLevel <= 0 {
		l.State = StateInitial
	}
}