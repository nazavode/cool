# Grammar of the Cool programming language
# as specified here:
# http://theory.stanford.edu/~aiken/software/cool/cool.html

language cool(go);

lang = "cool"
package = "github.com/nazavode/cool"
eventBased = true

:: lexer

%s initial;
%x inComment;

# Accept end-of-input in all states.
<*> eoi: /{eoi}/

whitespace: /[\n\r\t ]+/ (space)

invalid_token:
error:

# Multiline, nested comment blocks
<initial, inComment>
EnterBlockComment:  /\(\*/ (space) { l.enterBlockComment() }

<initial>
invalid_token: /\*\)/ (space)

<inComment> {
ExitBlockComment:  /\*\)/ (space) { l.exitBlockComment() }

BlockComment: /([^\*\(]|\*+[^\)]|\([^\*])*/ (space)
}

SingleLineComment: /\-\-.*/ (space)

# Identifiers
# All identifier rules conflict with keywords.
# We cannot use the (class) qualifier here
# b/c keywords in Cool aren't constant words
# (case insensitive make them parametric).
# Just use explicit priority.
ObjectId: /[a-z][_\w]*/ -1
TypeId  : /[A-Z][_\w]*/ -1

# Literals
IntegerLiteral: /\d+/
BoolLiteral   : /(t{R}{U}{E}|f{A}{L}{S}{E})/
StringLiteral : /"[^"]*"/

# Keywords (case insensitive)
class   : /{C}{L}{A}{S}{S}/
else    : /{E}{L}{S}{E}/
if      : /{I}{F}/
fi      : /{F}{I}/
in      : /{I}{N}/
inherits: /{I}{N}{H}{E}{R}{I}{T}{S}/
isvoid  : /{I}{S}{V}{O}{I}{D}/
let     : /{L}{E}{T}/
loop    : /{L}{O}{O}{P}/
pool    : /{P}{O}{O}{L}/
then    : /{T}{H}{E}{N}/
while   : /{W}{H}{I}{L}{E}/
case    : /{C}{A}{S}{E}/
esac    : /{E}{S}{A}{C}/
new     : /{N}{E}{W}/
of      : /{O}{F}/
not     : /{N}{O}{T}/

# Punctuation
'{': /\{/
'}': /\}/
'(': /\(/
')': /\)/
'.': /\./
';': /;/
',': /,/
'<': /</
'<=': /<=/
'>=': />=/
'@': /@/
'+': /\+/
'-': /-/
'*': /\*/
'/': /\//
'<-': /<-/
'~': /~/
':': /:/
'=': /=/

# Case insensitive chars
A = /(a|A)/
B = /(b|B)/
C = /(c|C)/
D = /(d|D)/
E = /(e|E)/
F = /(f|F)/
G = /(g|G)/
H = /(h|H)/
I = /(i|I)/
J = /(j|J)/
K = /(k|K)/
L = /(l|L)/
M = /(m|M)/
N = /(n|N)/
O = /(o|O)/
P = /(p|P)/
Q = /(q|Q)/
R = /(r|R)/
S = /(s|S)/
T = /(t|T)/
U = /(u|U)/
V = /(v|V)/
W = /(w|W)/
X = /(x|X)/
Y = /(y|Y)/
Z = /(z|Z)/

%%

${template go_lexer.stateVars-}
	commentLevel int // number of open nested block comments
${end}

${template go_lexer.initStateVars-}
	l.commentLevel = 0
${end}

${template go_lexer.lexer-}
${call base-}

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

${end}
