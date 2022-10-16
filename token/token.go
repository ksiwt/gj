package token

// Token identifies the type of Lex items.
type Token int

const (
	Unknown      Token = iota // unknown
	LeftBrace                 // {
	RightBrace                // }
	LeftBracket               // [
	RightBracket              // ]
	String                    // "foo"
	Number                    // 1
	True                      // true
	False                     // false
	Null                      // null
	Comma                     // ,
	Colon                     // :
	EOF                       // eof
	Error                     // error
)
