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
	{"KeywordClass", "class CLASS Class cLASS", []Token{CLASS, CLASS, CLASS, CLASS}},
	{"KeywordElse", "else ELSE Else eLSE", []Token{ELSE, ELSE, ELSE, ELSE}},
	{"KeywordIf", "if IF If iF", []Token{IF, IF, IF, IF}},
	{"KeywordFi", "fi FI Fi fI", []Token{FI, FI, FI, FI}},
	{"KeywordIn", "in IN In iN", []Token{IN, IN, IN, IN}},
	{"KeywordInherits", "inherits INHERITS Inherits iNHERITS", []Token{INHERITS, INHERITS, INHERITS, INHERITS}},
	{"KeywordIsvoid", "isvoid ISVOID Isvoid iSVOID", []Token{ISVOID, ISVOID, ISVOID, ISVOID}},
	{"KeywordLet", "let LET Let lET", []Token{LET, LET, LET, LET}},
	{"KeywordLoop", "loop LOOP Loop lOOP", []Token{LOOP, LOOP, LOOP, LOOP}},
	{"KeywordPool", "pool POOL Pool pOOL", []Token{POOL, POOL, POOL, POOL}},
	{"KeywordThen", "then THEN Then tHEN", []Token{THEN, THEN, THEN, THEN}},
	{"KeywordWhile", "while WHILE While wHILE", []Token{WHILE, WHILE, WHILE, WHILE}},
	{"KeywordCase", "case CASE Case cASE", []Token{CASE, CASE, CASE, CASE}},
	{"KeywordEsac", "esac ESAC Esac eSAC", []Token{ESAC, ESAC, ESAC, ESAC}},
	{"KeywordNew", "new NEW New nEW", []Token{NEW, NEW, NEW, NEW}},
	{"KeywordOf", "of OF Of oF", []Token{OF, OF, OF, OF}},
	{"KeywordNot", "not NOT Not nOT", []Token{NOT, NOT, NOT, NOT}},
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
