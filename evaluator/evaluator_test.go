package evaluator

import (
	"../lexer"
	"../object"
	"../parser"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		{"1 && 2", 2},
		{"0 && 10", 10},
		{"1 || 2", 1},
		{"true && 2 || 3", 2},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	// 各テストごとに新しい(=独立した)環境
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{`"Hello" == "Hello"`, true},
		{`"Hello" == "bye"`, false},
		{`"Hello" == "hello"`, false},
		{`"hello" != "bye"`, true},
		{`"hello" != "hello"`, false},
		{`"hello" == 5`, false},
		{`"hello" != 5`, true},
		{`"hello" == true`, false},
		{`"hello" == false`, false},
		{`"hello" != true`, true},
		{`"hello" != false`, true},
		{"5 == true", false},
		{"5 == false", false},
		{"5 != true", true},
		{"5 != false", true},
		{"5 <= 10", true},
		{"10 <= 5", false},
		{"5 <= 5", true},
		{"10 >= 5", true},
		{"5 >= 10", false},
		{"10 >= 10", true},
		{"true && true", true},
		{"true && false", false},
		{"false && false", false},
		{"true || true", true},
		{"true || false", true},
		{"false || false", false},
		// 短絡評価 (未定義変数は評価されないのでエラー吐かない)
		{"false && unknownVar", false},
		{"true || unknownVar", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	// NOTE: monkeyでは、falseでもnullでもないものは「全て」truthy!
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10 }", nil},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (0) { 10 } else { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{"if (10 > 1) { if (10 > 1) { return 10; }; return 1; }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				} 
				return 1;
			}
			`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identifier not found: foobar",
		},
		{
			`"Hello" - "world"`,
			"unknown operator: STRING - STRING",
		},
		{
			`{"name": "Monkey"}[fn(x) { x }];`,
			"unusable as hash key: FUNCTION",
		},
		{
			"true <= 5",
			"type mismatch: BOOLEAN <= INTEGER",
		},
		{
			"false <= false",
			"unknown operator: BOOLEAN <= BOOLEAN",
		},
		{
			"true >= 5",
			"type mismatch: BOOLEAN >= INTEGER",
		},
		{
			"false >= false",
			"unknown operator: BOOLEAN >= BOOLEAN",
		},
		{
			"true.5",
			"type mismatch: BOOLEAN . INTEGER",
		},
		{
			`"string".false`,
			"type mismatch: STRING . BOOLEAN",
		},
		{
			"5.true",
			"type mismatch: INTEGER . BOOLEAN",
		},
		{
			"true.true",
			"unknown operator: BOOLEAN . BOOLEAN",
		},
		{
			`"string"."string"`,
			"unknown operator: STRING . STRING",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error objects returned. got=%T (%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	// NOTE: このテストは変数に束縛された値もテストする
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	tests := []struct {
		input             string
		expectedNumParams int
		expectedBody      string
	}{
		{
			"fn(x) { x + 2; };",
			1,
			"(x + 2)",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		fn, ok := evaluated.(*object.Function)
		if !ok {
			t.Fatalf("object is not function. got=%T (%+v)",
				evaluated, evaluated)
		}

		if len(fn.Parameters) != tt.expectedNumParams {
			t.Fatalf("function does not contain %d parameters. Parameters=%+v",
				tt.expectedNumParams, fn.Parameters)
		}

		if fn.Body.String() != tt.expectedBody {
			t.Fatalf("body is not %q. got=%q", tt.expectedBody, fn.Body.String())
		}

	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(3);
	`
	testIntegerObject(t, testEval(input), 5)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello, world!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not string. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello, world!" {
		t.Errorf("String has wrong value. got=%q, expected=%q",
			str.Value, "Hello, world!")
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "world!"`
	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not string got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello world!" {
		t.Errorf("String has wrong value. got=%q, expected=%q",
			str.Value, "Hello world!")
	}
}

func TestBuildinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, "argument to `len` not supported, got INTEGER"},
		{`len("one", "two")`, "wrong number of arguments. got=2, want=1"},
		{`len([])`, 0},
		{`len([1])`, 1},
		{`len([1, 2, 3, 4])`, 4},
		{`len([[1, 2], [3, 4]])`, 2},
		{`len([1, 2], [3])`, "wrong number of arguments. got=2, want=1"},
		{`first([1])`, 1},
		{`first([2, 3])`, 2},
		{`first(1)`, "argument to `first` must be ARRAY, got INTEGER"},
		{`first([1, 2], [3])`, "wrong number of arguments. got=2, want=1"},
		{`first([])`, nil},
		{`first([[1, 2], [3]])`, []int64{1, 2}},
		{`last([1])`, 1},
		{`last([1, 2, 3])`, 3},
		{`last([[1, 2], [3, 4]])`, []int64{3, 4}},
		{`last([1, 2], [3])`, "wrong number of arguments. got=2, want=1"},
		{`last(1)`, "argument to `last` must be ARRAY, got INTEGER"},
		{`last([])`, nil},
		{`rest([1, 2])`, []int64{2}},
		{`rest([1])`, []int64{}},
		{`rest([1, 2, 3])`, []int64{2, 3}},
		{`rest(rest([1, 2, 3]))`, []int64{3}},
		{`rest([1, 2], [3])`, "wrong number of arguments. got=2, want=1"},
		{`rest(1)`, "argument to `rest` must be ARRAY, got INTEGER"},
		{`rest([])`, nil},
		{`let a = [1, 2]; rest(a); a;`, []int64{1, 2}},
		{`push([1, 2, 3], 4)`, []int64{1, 2, 3, 4}},
		{`push([], 1)`, []int64{1}},
		{`let a = [1]; let b = push(a, 2); b;`, []int64{1, 2}},
		{`let a = [1]; let b = push(a, 2); a;`, []int64{1}}, // no side-effect
		{`push([1, 2])`, "wrong number of arguments. got=1, want=2"},
		{`push(1, 1)`, "argument to `push` must be ARRAY, got INTEGER"},
		{`puts(1)`, nil},
		{`puts("foo")`, nil},
		{`puts(true)`, nil},
		{`puts([1])`, nil},
		{`puts({"foo": "bar"})`, nil},
		{`puts()`, nil},
		{`puts("one", "two)"`, nil},
		{`puts(1, "two", ["three"], {"four": "five"}, true)`, nil},
		// NOTE: 戻り値の型がNameSpaceのテストはTestBuildinNameSpaceFunctionsで行う
		{`fn() { outer(); }() == self()`, true},
		{`(namespace { let o = outer(); }).o == self()`, true},
		// NOTE: import成功例のテストはTestImportで行う　(ファイル探索システムがmain.goを
		// 基準に作られているので、パスを絶対参照する必要がある)
		{`import()`, "wrong number of arguments. got=0, want=1"},
		{`import(1)`, "argument to `import` must be STRING, got INTEGER"},
		{`import("_")`, "file could not open: _.monkey"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			testBooleanObject(t, evaluated, expected)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}

			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		case nil:
			testNullObject(t, evaluated)
		case []int64:
			testIntegerArray(t, evaluated, expected)
		}
	}
}

func TestBuildinNameSpaceFunctions(t *testing.T) {
	tests := []struct {
		input        string
		expectedType object.ObjectType
		expected     string
	}{
		{

			`
			let x = 3;
			self();
			`,
			object.NAMESPACE_OBJ,
			`namespace {x: 3}`,
		},
		{
			`
			fn() {
				let x = 3;
				self();
			}();
			`,
			object.NAMESPACE_OBJ,
			`namespace {x: 3}`,
		},
		{
			`fn(x) { self(); }(3);`,
			object.NAMESPACE_OBJ,
			`namespace {x: 3}`,
		},
		{
			`self(1);`,
			object.ERROR_OBJ,
			`wrong number of arguments. got=1, want=0`,
		},
		{
			`
			let x = 3;
			fn() { outer() }();
			`,
			object.NAMESPACE_OBJ,
			`namespace {x: 3}`,
		},
		{
			`outer(1);`,
			object.ERROR_OBJ,
			`wrong number of arguments. got=1, want=0`,
		},
		{
			`outer();`,
			object.NULL_OBJ,
			`null`,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch tt.expectedType {
		case object.NAMESPACE_OBJ:
			nameSpace, ok := evaluated.(*object.NameSpace)
			if !ok {
				t.Fatalf("Eval didn't return NameSpace. got=%T (%+v)",
					evaluated, evaluated)
			}

			if nameSpace.Inspect() != tt.expected {
				t.Fatalf("nameSpace has wrong Env. want=%q, got=%q",
					tt.expected, nameSpace.Inspect())
			}
		case object.ERROR_OBJ:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Fatalf("Eval didn't return Error. got=%T (%+v)",
					evaluated, evaluated)
			}

			if errObj.Message != tt.expected {
				t.Fatalf("errObj returned wrong message. want=%q, got=%q",
					tt.expected, errObj.Message)
			}
		case object.NULL_OBJ:
			testNullObject(t, evaluated)
		}
	}
}

func TestImport(t *testing.T) {
	// NOTE: importのファイル探索システムはmain.goを基準に作られているので、
	// ファイルは絶対参照する
	curDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("fail to get current dir: %s", err)
	}
	path := filepath.Join(filepath.Dir(curDir), "scripts")

	tests := []struct {
		input    string
		expected []int64
	}{
		{
			`let std = import("%s/std"); std.map([1, 2, 3], fn(x) { x * x; });`,
			[]int64{1, 4, 9},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(fmt.Sprintf(tt.input, path))
		testIntegerArray(t, evaluated, tt.expected)
	}
}

func testIntegerArray(t *testing.T, evaluated object.Object, expected []int64) bool {
	array, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not array. got=%T (%+v)", evaluated, evaluated)
		return false
	}

	if len(array.Elements) != len(expected) {
		t.Fatalf("array has wrong num of elements. got=%d (%+v), want=%d",
			len(array.Elements), array.Elements, len(expected))
		return false
	}

	for i, exp := range expected {
		if !testIntegerObject(t, array.Elements[i], exp) {
			return false
		}
	}

	return true
}

func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"
	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. expected=3, got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"let i = 0; [1][i]",
			1,
		},
		{
			"[1, 2, 3][1 + 1]",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i];",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `
	let two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		true: 5,
		false: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():   1,
		(&object.String{Value: "two"}).HashKey():   2,
		(&object.String{Value: "three"}).HashKey(): 3,
		(&object.Integer{Value: 4}).HashKey():      4,
		TRUE.HashKey():                             5,
		FALSE.HashKey():                            6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong length. got=%d, want=%d",
			len(result.Pairs), len(expected))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`let key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{true: 5}[true]`,
			5,
		},
		{
			`{false: 5}[false]`,
			5,
		},
		{
			`{true: 5}[5 == 5]`,
			5,
		},
		{
			`{"age": 5}["a" + "ge"]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestNameSpaceLiteral(t *testing.T) {
	tests := []struct {
		input          string
		expectedIdents []string
		expectedVals   []interface{}
	}{
		{
			`
			namespace {
				let x = 1;
				let cond = true;
				let y = x + 2;
			}
			`,
			[]string{"x", "cond", "y"},
			[]interface{}{1, true, 3},
		},
		{
			`
			let mySpace = namespace { let x = 1; };
			mySpace;
			`,
			[]string{"x"},
			[]interface{}{1},
		},
		{
			`
			let mySpace = namespace { let x = 1; };
			let alias = mySpace;
			alias;
			`,
			[]string{"x"},
			[]interface{}{1},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		nameSpace, ok := evaluated.(*object.NameSpace)
		if !ok {
			t.Fatalf("Eval didn't return NameSpace. return=%T (%+v)",
				evaluated, evaluated)
		}

		for i, ident := range tt.expectedIdents {
			val, ok := nameSpace.Env.Get(ident)
			if !ok {
				t.Errorf("identifier %s wasn't bound to NameSpace", ident)
			}

			switch expectedVal := tt.expectedVals[i].(type) {
			case int:
				testIntegerObject(t, val, int64(expectedVal))
			case bool:
				testBooleanObject(t, val, expectedVal)
			}
		}
	}
}

func TestNameSpaceAsOOP(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`
			let Person = namespace {
				let new = fn(age) { self(); };
			};
			let person = Person.new(20);
			person.age;
			`,
			20,
		},
		{
			`
			let Person = namespace {
				let new = fn(age) { self(); };
				let canDrink = fn() { age >= 20; };
			};
			let person = Person.new(30);
			person.canDrink();
			`,
			true,
		},
		{
			// namespaceはimmutableなので、代わりに値の違う新たなnamespaceを返す
			`
			let Person = namespace {
				let new = fn(age) { self(); };
				let reachBirthDay = fn() { outer().new(age + 1); };
			};
			let person = Person.new(14);
			let person = person.reachBirthDay();
			person.age;
			`,
			15,
		},
		{
			`
			let Person = namespace {
				let new = fn(age) { self(); };
				let isOlder = fn(other) { age > other.age };
			};
			let mike = Person.new(14);
			let judy = Person.new(18);
			judy.isOlder(mike);
			`,
			true,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			testBooleanObject(t, evaluated, expected)
		}
	}
}

func TestNameSpaceLiteralScopes(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			`
			let mySpace = namespace { let x = 10; };
			let x = 1;
			x;
			`,
			1,
		},
		{
			`
			let x = 1;
			let mySpace = namespace { let x = 10; };
			x;
			`,
			1,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		integer, ok := evaluated.(*object.Integer)
		if !ok {
			t.Fatalf("Eval didn't return Integer. got=%T (%+v)",
				evaluated, evaluated)
		}

		testIntegerObject(t, integer, tt.expected)
	}
}

func TestEvalDotExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			`
			let mySpace = namespace {
				let five = 5;
			};
			mySpace.five;
			`,
			5,
		},
		{
			`
			let mySpace = namespace {
				let five = fn() { 5; };
			};
			mySpace.five();
			`,
			5,
		},
		{
			`
			let mySpace = namespace {
				let childSpace = namespace {
					let five = 5;
				};
			};
			mySpace.childSpace.five;
			`,
			5,
		},
		{
			`let mySpace = namespace {
				let five = 5;
			};
			let f = fn() { mySpace; };
			f().five;
			`,
			5,
		},
		{
			`namespace { let x = 5; }.x;`,
			5,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
