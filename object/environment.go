package object

import (
	"bytes"
	"fmt"
	"strings"
)

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

type Environment struct {
	store map[string]Object
	// 内部のenvironmentに存在しない束縛は、外側のenvironmentを参照する
	outer *Environment
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	// 内部のenvironmentに存在しない束縛は、外側のenvironmentを参照する
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

func (e *Environment) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for k, v := range e.store {
		// NOTE: 下記の問題を回避するため、namespace内のnamespaceは略記する
		// namespaceをそのまま表示すると自己参照で無限ループ
		// `namespace { let ns = self(); };`
		switch v.Type() {
		case NAMESPACE_OBJ:
			// 内部は省略
			pairs = append(pairs, fmt.Sprintf("%s: namespace {...}", k))
		default:
			// 束縛されたobjectをinspect
			pairs = append(pairs, fmt.Sprintf("%s: %s", k, v.Inspect()))
		}
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// 組み込み関数outer()の評価に必要なためgetterとして追加
func (e *Environment) Outer() *Environment {
	return e.outer
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	// 関数内等のローカルな名前空間を返す
	env := NewEnvironment()
	env.outer = outer
	return env
}
