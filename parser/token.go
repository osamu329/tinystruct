package parser

//go:generate stringer -type Token
type Token int

const (
	INVALID Token = iota
	EOF
	IDENT
	STRUCT
	TYPEDEF
	COMMENT
	LBRACE
	RBRACE
	LBRACK
	RBRACK
	RUNE
	INT
	LF
	SEMICOLON
)

var keywords = map[string]Token{
	"struct":  STRUCT,
	"typedef": TYPEDEF,
}
