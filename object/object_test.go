package object

import (
	"testing"
)

// ハッシュキーによるハッシュの要素探索
// (直接Objectどうしを比較しても、.Valueが同じでも別のポインタであるため参照不可能)
// &String{Value: "name"} != &String{Value: "name"}
// 一方、全てのkeyの.Valueを使って照合する場合O(n)かかる…
// => ハッシュキーによりObjectの同値性を確かめる

func TestStringHashkey(t *testing.T) {
	hello1 := &String{Value: "hello world"}
	hello2 := &String{Value: "hello world"}
	diff1 := &String{Value: "My name is Jonney"}
	diff2 := &String{Value: "My name is Jonney"}

	if hello1.HashKey() != hello2.HashKey() {
		t.Errorf("strings with same contants has different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("strings with same contants has different hash keys")
	}

	if hello1.HashKey() == diff1.HashKey() {
		t.Errorf("strings with different contants has same hash keys")
	}
}

func TestBooleanHashKeys(t *testing.T) {
	true1 := &Boolean{Value: true}
	true2 := &Boolean{Value: true}
	false1 := &Boolean{Value: false}
	false2 := &Boolean{Value: false}

	if true1.HashKey() != true2.HashKey() {
		t.Errorf("booleans with same contants has different hash keys")
	}

	if false1.HashKey() != false2.HashKey() {
		t.Errorf("booleans with same contants has different hash keys")
	}

	if true1.HashKey() == false2.HashKey() {
		t.Errorf("booleans with different contants has same hash keys")
	}
}

func TestIntegerHashKeys(t *testing.T) {
	int1 := &Integer{Value: 1}
	int2 := &Integer{Value: 1}
	diff1 := &Integer{Value: 2}
	diff2 := &Integer{Value: 2}

	if int1.HashKey() != int2.HashKey() {
		t.Errorf("integers with same contants has different hash keys")
	}

	if diff1.HashKey() != diff2.HashKey() {
		t.Errorf("integers with same contants has different hash keys")
	}

	if int1.HashKey() == diff2.HashKey() {
		t.Errorf("integers with different contants has same hash keys")
	}
}
