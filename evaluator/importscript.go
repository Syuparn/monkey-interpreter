package evaluator

import (
	"../object"
	"fmt"
)

func importScript(env *object.Environment, fileName string) object.Object {
	// NOTE: use init() to avoid initialization loop (see below)
	// same as        "EvalScriptFile(fileName)"
	importEnv, err := _evalScriptFile(fileName)
	if err != nil {
		return newError(fmt.Sprintf("%s", err))
	}

	return &object.NameSpace{Env: importEnv}
}

// NOTE: use init() to avoid initialization loop
// builtin -> importScript -> EvalScriptFile
// -> Eval -> evalIdentifier -> builtin ->...
var _evalScriptFile func(string) (*object.Environment, error)

// (initは特殊な関数で、実行の直前に自動的に呼び出される(コード内での呼び出しは不能)
// 宣言を実行の直前(importの初期化の直後)まで遅延させるためループしなくなる)
func init() {
	_evalScriptFile = EvalScriptFile
}
