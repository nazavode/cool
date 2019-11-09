package cool

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

var lextests = []struct {
	name   string
	source string
	tokens []Token
}{
	{"Empty", "", nil},
	{"Identifiers", "object Type oBJECT", []Token{OBJECTID, TYPEID, OBJECTID}},
	{"IntegerLiterals", "0 000 0000 01234567890", []Token{INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL}},
	{"StringLiterals", "\"\" \" \" \" foo \"", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	{"LineComment", "class -- this is a class\n  class", []Token{CLASS, CLASS}},
	{"BlockComment0", "(**)", nil},
	{"BlockComment1", "(*(**)*)", nil},
	{"BlockComment2", "(*(*(**)*)*)", nil},
	{"BlockComment3", "(*(*(**)*)(*(**)*)*)", nil},
	{"BlockComment4", "(*else(*then(*if*)*)class(*(*loop*)pool*)case*)", nil},
	{"Whitespace", "    \t\t \f \v \r\r\r\n\n      ", nil},
	{"BoolLiterals", "true false tRUE fALSE True False", []Token{BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, TYPEID, TYPEID}},
}

func scan(source string) []Token {
	var lex Lexer
	lex.Init(source)
	var tokens []Token
	for cur := lex.Next(); cur != EOI; cur = lex.Next() {
		tokens = append(tokens, cur)
	}
	return tokens
}

func TestLexer(t *testing.T) {
	for _, tt := range lextests {
		t.Run(tt.name, func(t *testing.T) {
			got := scan(tt.source)
			if diff := cmp.Diff(tt.tokens, got); diff != "" {
				t.Errorf("lex mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
