package ast

import (
	"../token"
)

type Node interface { // 範疇(category)を扱うインターフェース
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode() // Expressionとの混同をコンパイルエラーにするためのダミーメソッド
}

type Expression interface {
	Node
	expressionNode() // Statementとの混同をコンパイルエラーにするためのダミーメソッド
}

type Program struct { // プログラム全体(=S)
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // token.Let
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token // 'ident' token
	Value string
}

// NOTE: let文は値を返さないが、後で値を返すidentifierも作るのでexpressionにする
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token // 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}
