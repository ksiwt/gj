package token

// Token identifies the type of Lex items.
type Token int

const (
	Unknown Token = iota
	LeftBrace
	RightBrace
	LeftBracket
	RightBracket
	String
	Number
	True
	False
	Null
	Comma
	Colon
	EOF
	Error
)
