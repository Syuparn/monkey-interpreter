package lexer

import (
	"../token"
)

type Lexer struct {
	input        string
	position     int  // 入力における現在読んでいる文字の位置
	readPosition int  // 次に読む文字の位置(token区切りを判断するために先読み)
	ch           byte // 現在検査中の文字
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // inputの終端に到達したらnull文字セット
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// 次の文字をのぞき見(peek) (readCharと違い読み進めない)
func (l *Lexer) peekChar() byte {
	if l.position >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// Lexerのコンストラクタ
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // 一文字目を読みこんでおく
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpace() // 空白読み飛ばさないと不要なILLIGAL tokenが生成されてしまう

	switch l.ch {
	case '=':
		if l.peekChar() == '=' { // '=='
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case ':':
		tok = newToken(token.COMMA, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF //newTokenで生成しないのは、null文字をstringで変換できないため？
	default:
		if isLetter(l.ch) {
			// 記号とは別処理なのでreadCharしない(ident終わるまで塊で１tokenとして読むため)
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal) // 変数名 or keyword
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch) // 使用不可能記号
		}
	}

	l.readChar() // 次の文字を読む
	return tok
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		// end of string
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
