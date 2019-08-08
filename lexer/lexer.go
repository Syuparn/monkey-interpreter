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

// Lexerのコンストラクタ
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // 一文字目を読みこんでおく
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF //newTokenで生成しないのは、null文字をstringで変換できないため？
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
