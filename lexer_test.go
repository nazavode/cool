package cool

import (
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"log"
	"path"
	"testing"
)

var testFiles = []struct {
	file   string
	tokens []Token
}{
	{"testdata/pathologicalstrings.cool", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	{"testdata/nestedcomment.cool", nil},
	{"testdata/s04.test.cool", []Token{INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL}},
	{"testdata/s05.test.cool", []Token{STRINGLITERAL}},
	{"testdata/s14.test.cool", []Token{OBJECTID, OBJECTID}},
	{"testdata/s16.test.cool", []Token{INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL}},
	{"testdata/s25.test.cool", []Token{OBJECTID, OBJECTID, OBJECTID, OBJECTID, TYPEID, TYPEID, TYPEID, TYPEID}},
	{"testdata/s26.test.cool", []Token{INTEGERLITERAL, OBJECTID}},
	{"testdata/stringcomment.cool", []Token{STRINGLITERAL, STRINGLITERAL}},
	{"testdata/twice_512_nested_comments.cl.cool", []Token{OBJECTID}},
	{"testdata/wq0607-c1.cool", []Token{STRINGLITERAL}},
	{"testdata/longstring_escapedbackslashes.cool", []Token{STRINGLITERAL, STRINGLITERAL}},
	{"testdata/s19.test.cool", []Token{INVALID_TOKEN, OBJECTID, INVALID_TOKEN}},
	{"testdata/s31.test.cool", []Token{INVALID_TOKEN}},
	{"testdata/s32.test.cool", []Token{INVALID_TOKEN}},
	{"testdata/s33.test.cool", []Token{OBJECTID, INVALID_TOKEN}},
	{"testdata/s34.test.cool", []Token{OBJECTID, INVALID_TOKEN}},
	{"testdata/wq0607-c1.cool", []Token{STRINGLITERAL}},
	{"testdata/wq0607-c2.cool", []Token{STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL}},
	{"testdata/null_in_code.cl.cool", []Token{OBJECTID, OBJECTID, OBJECTID, OBJECTID, ASSIGNGT, INVALID_TOKEN, LTMINUS, RPAREN}},
	{"testdata/null_in_string_unescaped_newline.cl.cool", []Token{INVALID_TOKEN, OBJECTID, PLUS}},
	{"testdata/longcomment.cool", []Token{TYPEID, OBJECTID, TYPEID, OBJECTID, TYPEID, OBJECTID, TYPEID, OBJECTID, TYPEID, TYPEID, OBJECTID, OBJECTID, OBJECTID, SEMICOLON, OBJECTID}},
	// TODO currently failing
	// {"testdata/wq0607-c3.cool", []Token{STRINGLITERAL, STRINGLITERAL, INVALID_TOKEN, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, STRINGLITERAL, INVALID_TOKEN, STRINGLITERAL, INVALID_TOKEN, INVALID_TOKEN}},
	// {"testdata/wq0607-c4.cool", []Token{}},
	// {"testdata/null_in_string.cl.cool", []Token{INVALID_TOKEN}},
	// {"testdata/null_in_string_followed_by_tokens.cl.cool", []Token{INVALID_TOKEN, OBJECTID, PLUS}},
}

var testSnippets = []struct {
	name   string
	source string
	tokens []Token
}{
	{"LineComment", "class -- this is a class\n  class", []Token{CLASS, CLASS}},
	{"BlockComment0", "(**)", nil},
	{"BlockComment1", "(*(**)*)", nil},
	{"BlockComment2", "(*(*(**)*)*)", nil},
	{"BlockComment3", "(*(*(**)*)(*(**)*)*)", nil},
	{"BlockComment4", "class (*else(*then(*if*)*)class(*(*loop*)pool*)case*) class", []Token{CLASS, CLASS}},
	{"BlockComment5", "class(************)class", []Token{CLASS, CLASS}},
	{"BlockComment6", "class(***)class", []Token{CLASS, CLASS}},
	{"Empty", "", nil},
	{"Identifiers", "object Type oBJECT", []Token{OBJECTID, TYPEID, OBJECTID}},
	{"IntegerLiterals", "0 000 0000 01234567890", []Token{INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL}},
	{"StringLiterals", "\"\" \" \" \" foo \"", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
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
	{"ValidChars", "+/-*=<.~,;:()@{}", []Token{PLUS, DIV, MINUS, MULT, ASSIGN, LT, DOT, TILDE, COMMA, SEMICOLON, COLON, LPAREN, RPAREN, ATSIGN, LBRACE, RBRACE}},
	{"NullInCode", "class\000if", []Token{CLASS, INVALID_TOKEN, IF}},
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

func TestLexerSnippets(t *testing.T) {
	for _, tt := range testSnippets {
		t.Run(tt.name, func(t *testing.T) {
			got := scan(tt.source)
			if diff := cmp.Diff(tt.tokens, got); diff != "" {
				t.Errorf("lex mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLexerFiles(t *testing.T) {
	for _, tt := range testFiles {
		name := path.Base(tt.file)
		data, err := ioutil.ReadFile(tt.file)
		if err != nil {
			log.Fatalln(err)
		}
		source := string(data)
		t.Run(name, func(t *testing.T) {
			got := scan(source)
			if diff := cmp.Diff(tt.tokens, got); diff != "" {
				t.Errorf("lex mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
