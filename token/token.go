package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// 識別子＋リテラル
	IDENT = "IDENT" // 変数x,y...
	INT   = "INT"   // 数1,2...

	// 演算子
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	DOT      = "."

	LT = "<"
	GT = ">"

	EQ     = "=="
	NOT_EQ = "!="
	GEQ    = ">="
	LEQ    = "<="
	AND    = "&&"
	OR     = "||"

	// デリミタ
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION  = "FUNCTION"
	LET       = "LET"
	TRUE      = "TRUE"
	FALSE     = "FALSE"
	IF        = "IF"
	ELSE      = "ELSE"
	RETURN    = "RETURN"
	NAMESPACE = "NAMESPACE"

	STRING = "STRING"
)

var keywords = map[string]TokenType{
	"fn":        "FUNCTION",
	"let":       "LET",
	"true":      "TRUE",
	"false":     "FALSE",
	"if":        "IF",
	"else":      "ELSE",
	"return":    "RETURN",
	"namespace": "NAMESPACE",
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
