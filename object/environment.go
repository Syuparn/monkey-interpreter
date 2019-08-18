package object

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

func NewEnclosedEnvironment(outer *Environment) *Environment {
	// 関数内等のローカルな名前空間を返す
	env := NewEnvironment()
	env.outer = outer
	return env
}
