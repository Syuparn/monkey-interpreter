package object

import (
	"../ast"
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"
)

// NOTE: HashKey()によるハッシュの要素探索
// {"name": "Taro"}["name"]
// 直接Objectどうしを比較しても、.Valueが同じでも別のポインタであるため参照不可能
// &String{Value: "name"} != &String{Value: "name"}
// 一方、全てのkeyの.Valueを使って照合する場合O(n)かかる…
// => ハッシュキーによりObjectの同値性を確かめる
type HashKey struct {
	Type  ObjectType // Valueが偶然同じでも、型が違えばハッシュキーは等しくない
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

// monkeyでは、評価された値の内部表現は全てObjectで表される
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILDIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	NAMESPACE_OBJ    = "NAMESPACE"
)

// NOTE: 内部表現によってフィールドが違う(boolとint等)のでstructではなくinterface
type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: b.Type(), Value: value}
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }

// NOTE: ErrorはReturnValueのように使う
// (エラーもreturnも「その先の評価を中断し脱出」という点で同じ)
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

// EnvをFunctionのフィールドにしたのは、クロージャを実現するため
// (Envに入るのは関数が作られたときの、この関数のすぐ外側の名前空間。
// そのため、関数outer内で関数innerを生成すると、inner内ではouterの束縛は
// 「outerの外側で呼び出されたときも」参照可能！)
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// NOTE: 過去に計算したhashkeyをキャッシュすることでString.HashKey()を
// 1.2~8倍高速化 (詳細はcashhashkey_test.go)
var stringHashKeyCashes = make(map[string]HashKey)

func (s *String) HashKey() HashKey {
	// NOTE: ごくまれに別の文字列に同じ整数が与えられてしまう(ハッシュの衝突)
	// 実用上問題ないが、絶対に一意にするには「チェイン法」、「オープンアドレス法」用いる

	// hashkey既に求めていたら使い回し
	if hashKey, ok := stringHashKeyCashes[s.Value]; ok {
		return hashKey
	} else {
		h := fnv.New64a()
		h.Write([]byte(s.Value))

		return HashKey{Type: s.Type(), Value: h.Sum64()}
	}
}

// NOTE: envを引数に追加
// (import, self, outer等は、評価の際今のスコープを知る必要があるため)
type BuiltinFunction func(env *Environment, args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILDIN_OBJ }
func (b *Builtin) Inspect() string  { return "buildin function" }

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

// NOTE: map[HashKey]Objectとしてもハッシュにはなるが、keyを参照しづらくなる
// イテレートするときやkeyとのペアを取得するにはHashPairがあるほうが都合がいい
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+": "+pair.Value.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type NameSpace struct {
	Env *Environment
}

func (ns *NameSpace) Type() ObjectType { return NAMESPACE_OBJ }
func (ns *NameSpace) Inspect() string {
	var out bytes.Buffer
	out.WriteString("namespace ")
	out.WriteString(ns.Env.Inspect())

	return out.String()
}
