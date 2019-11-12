package cool

import (
	"bytes"
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"io/ioutil"
	"log"
	"testing"
	"text/template"
)

type SourceToken struct {
	Line     int    `json:"line"`
	Terminal Token  `json:"token"`
	Value    string `json:"source"`
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
	temp, err := template.New("golden").Parse("{{.}}.lexer.gold.json")
	if err != nil {
		panic(err)
	}
	for _, sourceFileName := range testFiles {
		t.Run(sourceFileName, func(t *testing.T) {
			// Source
			sourceBuf, err := ioutil.ReadFile(sourceFileName)
			if err != nil {
				log.Fatalln(err)
			}
			source := string(sourceBuf)
			sourceTokens := scanSource(source)
			// Golden
			var b bytes.Buffer
			err = temp.Execute(&b, sourceFileName)
			if err != nil {
				log.Fatalln(err)
			}
			goldFileName := b.String()
			goldBuf, err := ioutil.ReadFile(goldFileName)
			if err != nil {
				log.Fatalln(err)
			}
			var goldTokens []SourceToken
			err = json.Unmarshal(goldBuf, &goldTokens)
			if err != nil {
				log.Fatalln(err)
			}
			// Compare
			if diff := cmp.Diff(goldTokens, sourceTokens); diff != "" {
				t.Errorf("lex mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func scanSource(source string) []SourceToken {
	var lex Lexer
	lex.Init(source)
	var tokens []SourceToken
	for cur := lex.Next(); cur != EOI; cur = lex.Next() {
		tokens = append(tokens, SourceToken{lex.Line(), cur, lex.Text()})
	}
	return tokens
}

func scan(source string) []Token {
	var tokens []Token
	for _, t := range scanSource(source) {
		tokens = append(tokens, t.Terminal)
	}
	return tokens
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
	{"Identifier", "object Type oBJECT", []Token{OBJECTID, TYPEID, OBJECTID}},
	{"IntegerLiteral", "0 000 0000 01234567890", []Token{INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL, INTEGERLITERAL}},
	{"StringLiteral", "\"\" \" \" \" foo \"", []Token{STRINGLITERAL, STRINGLITERAL, STRINGLITERAL}},
	{"Whitespace", "    \t\t \f \v \r\r\r\n\n      ", nil},
	{"BoolLiteral", "true false tRUE fALSE True False", []Token{BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, BOOLLITERAL, TYPEID, TYPEID}},
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
	{"NullInCode", "class\000if", []Token{CLASS, INVALID_TOKEN, IF}},
	{"TerminalPlus", "+", []Token{PLUS}},
	{"TerminalDiv", "/", []Token{DIV}},
	{"TerminalMinus", "-", []Token{MINUS}},
	{"TerminalMult", "*", []Token{MULT}},
	{"TerminalAssign", "=", []Token{ASSIGN}},
	{"TerminalLt", "<", []Token{LT}},
	{"TerminalDot", ".", []Token{DOT}},
	{"TerminalTilde", "~", []Token{TILDE}},
	{"TerminalComma", ",", []Token{COMMA}},
	{"TerminalSemicolon", ";", []Token{SEMICOLON}},
	{"TerminalColon", ":", []Token{COLON}},
	{"TerminalLParen", "(", []Token{LPAREN}},
	{"TerminalRParen", ")", []Token{RPAREN}},
	{"TerminalAt", "@", []Token{ATSIGN}},
	{"TerminalLBrace", "{", []Token{LBRACE}},
	{"TerminalRBrace", "}", []Token{RBRACE}},
	{"InvalidBang", "!", []Token{INVALID_TOKEN}},
	{"InvalidHash", "#", []Token{INVALID_TOKEN}},
	{"InvalidDollar", "$", []Token{INVALID_TOKEN}},
	{"InvalidPercent", "%", []Token{INVALID_TOKEN}},
	{"InvalidHat", "^", []Token{INVALID_TOKEN}},
	{"InvalidPound", "&", []Token{INVALID_TOKEN}},
	{"InvalidUnderscore", "_", []Token{INVALID_TOKEN}},
	{"InvalidGT", ">", []Token{INVALID_TOKEN}},
	{"InvalidQuestion", "?", []Token{INVALID_TOKEN}},
	{"InvalidBacktick", "`", []Token{INVALID_TOKEN}},
	{"InvalidLSub", "[", []Token{INVALID_TOKEN}},
	{"InvalidRSub", "]", []Token{INVALID_TOKEN}},
	{"InvalidBackslash", "\\", []Token{INVALID_TOKEN}},
	{"InvalidPipe", "|", []Token{INVALID_TOKEN}},
	{"TokenAndInvalidBang", "a ! a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidHash", "a # a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidDollar", "a $ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidPercent", "a % a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidHat", "a ^ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidPound", "a & a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidUnderscore", "a _ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidGT", "a > a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidQuestion", "a ? a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidBacktick", "a ` a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidLSub", "a [ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidRSub", "a ] a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidBackslash", "a \\ a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
	{"TokenAndInvalidPipe", "a | a", []Token{OBJECTID, INVALID_TOKEN, OBJECTID}},
}

var testFiles = []string{
	"testdata/all_else_true.cl.cool",
	"testdata/backslash.cool",
	"testdata/backslash2.cool",
	"testdata/badidentifiers.cool",
	"testdata/badkeywords.cool",
	"testdata/bothcomments.cool",
	"testdata/comment_in_string.cl.cool",
	"testdata/endcomment.cool",
	"testdata/eofstring.cool",
	"testdata/escaped_chars_in_comment.cl.cool",
	"testdata/escapedeof.cool",
	"testdata/escapedquote.cool",
	"testdata/escapedunprintables.cool",
	"testdata/integers2.cool",
	"testdata/invalidcharacters.cool",
	"testdata/invalidinvisible.cool",
	"testdata/keywords.cool",
	"testdata/lineno2.cool",
	"testdata/lineno3.cool",
	"testdata/longcomment.cool",
	"testdata/longstring_escapedbackslashes.cool",
	"testdata/multilinecomment.cool",
	"testdata/nestedcomment.cool",
	"testdata/null_in_code.cl.cool",
	"testdata/null_in_string_unescaped_newline.cl.cool",
	"testdata/opencomment.cool",
	"testdata/operators.cool",
	"testdata/pathologicalstrings.cool",
	"testdata/s03.test.cool",
	"testdata/s04.test.cool",
	"testdata/s05.test.cool",
	"testdata/s14.test.cool",
	"testdata/s16.test.cool",
	"testdata/s19.test.cool",
	"testdata/s25.test.cool",
	"testdata/s26.test.cool",
	"testdata/s31.test.cool",
	"testdata/s32.test.cool",
	"testdata/s33.test.cool",
	"testdata/s34.test.cool",
	"testdata/simplestrings.cool",
	"testdata/stringcomment.cool",
	"testdata/stringwithescapes.cool",
	"testdata/twice_512_nested_comments.cl.cool",
	"testdata/validcharacters.cool",
	"testdata/weirdcharcomment.cool",
	"testdata/wq0607-c1.cool",
	"testdata/wq0607-c2.cool",
	"testdata/wq0607-c3.cool",
	"testdata/wq0607-c4.cool",
	"testdata/arith.cool",
	"testdata/atoi.cool",
	"testdata/book_list.cl.cool",
	"testdata/hairyscary.cool",
	"testdata/io.cool",
	"testdata/life.cool",
	"testdata/new_complex.cool",
	"testdata/objectid.test.cool",
	"testdata/palindrome.cool",
	"testdata/sort_list.cl.cool",
}

// TODO
// FAIL {"testdata/escapednull.cool", []Token{INVALID_TOKEN}},
// FAIL {"testdata/null_in_string.cl.cool", []Token{INVALID_TOKEN}},
// FAIL {"testdata/null_in_string_followed_by_tokens.cl.cool", []Token{INVALID_TOKEN, OBJECTID, PLUS}},
