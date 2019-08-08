package lexer

type Lexer struct {
	input        string
	position     int  // 入力における現在読んでいる文字の位置
	readPosition int  // 次に読む文字の位置(token区切りを判断するために先読み)
	ch           byte // 現在検査中の文字
}

func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.ch = 0 // inputの終端に到達したらnull文字セット
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// Lexerのコンストラクタ
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // 一文字目を読みこんでおく
	return l
}
