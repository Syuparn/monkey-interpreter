package evaluator

import (
	"../object"
	"fmt"
)

// NOTE: envを引数に追加
// (import, self, outer等は、評価の際今のスコープを知る必要があるため)
var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
	},
	"first": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}

			return NULL
		},
	},
	"last": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `last` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}

			return NULL
		},
	},
	"rest": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `rest` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElements := make([]object.Object, length-1, length-1)
				copy(newElements, arr.Elements[1:length])
				return &object.Array{Elements: newElements}
			}

			return NULL
		},
	},
	"push": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2",
					len(args))
			}
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s",
					args[0].Type())
			}

			arr := args[0].(*object.Array)
			length := len(arr.Elements)

			newElements := make([]object.Object, length+1, length+1)
			copy(newElements, arr.Elements)
			newElements[length] = args[1]

			return &object.Array{Elements: newElements}
		},
	},
	"puts": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return NULL
		},
	},
	"self": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got=%d, want=0",
					len(args))
			}

			return &object.NameSpace{Env: env}
		},
	},
	"outer": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 0 {
				return newError("wrong number of arguments. got=%d, want=0",
					len(args))
			}

			if env.Outer() == nil {
				return NULL
			}

			return &object.NameSpace{Env: env.Outer()}
		},
	},
	"import": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			if args[0].Type() != object.STRING_OBJ {
				return newError("argument to `import` must be STRING, got %s",
					args[0].Type())
			}

			fileStem := args[0].(*object.String).Value

			// NOTE: namespace中THIS_DIR, THIS_FILEにファイル名をSTRINGで束縛
			return importScript(env, fileStem+".monkey")
		},
	},
	"type": &object.Builtin{
		Fn: func(env *object.Environment, args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}

			return &object.String{Value: string(args[0].Type())}
		},
	},
}
