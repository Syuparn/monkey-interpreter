# Monkey language interpreter
code for "Writing An Interpreter In Go"

「Go言語でつくるインタプリタ」(Thorsten Ball著)を読んで実装したものです。

This is Monkey interpreter from "Writing An Interpreter In Go" by Thorsten Ball.

[O'Reilly Japan - Go言語でつくるインタプリタ](https://www.oreilly.co.jp/books/9784873118222/)

[Writing An Interpreter In Go \| Thorsten Ball](https://interpreterbook.com/)

# Differences from original
## String comparison

Strings are same means their values (not addresses) are same.

```
>> "monkey" == "monkey"
true
```

## `>=`, `<=`

```
>> 4 >= 4
true
>> 10 <= 9
false
```

## `&&`, `||` with shortcut

```
>> 1 == 2 && 2 == 2
false
>> true && 2
2
>> 5 || puts("hi")
5
>> false || puts("hi")
hi
null
```

## Running script files

You can run monkey script file by `-f` command.

```
./monkey -f myscript.monkey
```
or
```
go run main.go -f myscript.monkey
```

Of course you can also use REPL mode, which is familiar with you.

```
> ./monkey
Hello (your name)! This is the Monkey programming language!
Feel free to type in commands
```

## `type()`

```
>> type(1)
INTEGER
>> type("Monkey")
STRING
>> type("Monkey") == "STRING"
true
```

## `namespace {}`

`namespace` is an object of enclosed environment.

```
>> let a = 10;
>> namespace { let a = 5; };
namespace {a: 5}
>> puts(a);
10
null
```

`namespace` is a first-class object.

```
>> let mySpace = namespace { let a = 5; };
>> puts(mySpace);
namespace {a: 5}
```

You can access variables in `namespace` by `.`.

```
>> let mySpace = namespace { let a = 5;  let add = fn(x, y) {x + y}; };
>> mySpace.a
5
>> mySpace.add(1, 3)
4
```

`namespace` can be nested.

```
>> namespace { let ns = namespace { let one = 1; }; }.ns.one
1
```

Notice: `.` operator works weirdly in some cases.
`namespace` is just an object of an environment.
This means outer variables and literals can be referred by `.` operator.

```
>> let name = "John";
>> let ns = namespace { let a = 5; };
>> ns.name
John
>> ns.5
5
```

### `import()`

Notice: This function works only if you build exective monkey file (due to file reference system).

```
# ok
.../monkey$ go build
.../monkey$ ./monkey
#NG
.../monkey$ go run main.go
```

`import()` reads a script file and returns it as `namespace`.

```monkey:sample.monkey
let add = fn(x, y) { x + y; };
```

```
>> let sample = import("(...)/sample");
>> sample.add(1, 2)
3
```

If you want to import script in the same directory, const `THIS_DIR` may help you.

```
import(THIS_DIR + "otherscript")
```

#### Standard functions

Besides your scripts, you can import scripts in monkey/scripts without full-path.

```
>> let std = import("std");
>> std.abs(-3)
3
>> std.filter([1, 2, 4, 8], fn(x) { x < 4 })
[1, 2]
```

### `self()`, `outer()` for OOP

`self()` returns `namespace` with the current environment.

Likewise, `outer()` returns `namespace` with the outer environment of current one.

These built-in functions realize class-like system like this.

```monkey:person.monkey
let Person = namespace {
    let new = fn(age, name) {
        self();
    };
    
    let sayHi = fn() {
        puts("hi, I'm " + name);
    };
    
    let isElder = fn(other) {
        age > other.age
    };
    
    let grow = fn() {
        outer().new(age + 1, name)
    };
};
```

With namespace `Person`, you can write codes like below.

```
>> let Person = import("(...)\person").Person
>> let tom = Person.new(10, "Tom")
>> let judy = Person.new(14, "Judy")
>> tom
namespace {age: 10, name: Tom}
>> judy
namespace {age: 14, name: Judy}
>> tom.sayHi()
hi, I'm Tom
null
>> judy.sayHi()
hi, I'm Judy
null
>> judy.isElder(tom)
true
>> let tom = tom.grow()
>> tom
namespace {age: 11, name: Tom}
```

#### How do `self()` and `outer()` work?

```monkey:person.monkey
let Person = namespace {
    let new = fn(age, name) {
        self();
    };
    
    let sayHi = fn() {
        puts("hi, I'm " + name);
    };
    ...
};
```

```
>> let tom = Person.new(10, "Tom")
>> tom.sayHi()
```

`self()` can be used for constructor, returning namespace which contains its own variables and can access to other methods.

In `Person.new`, `self()` returns namespace of "current" enviroment, the environment in function `new` containing arguments `age` and `name`.
Also, since the namespace is inner of namespace `Person`, all functions in `Person` can called by `.`.

For example. `Person.new(10, "Tom")` returns environment of the function.
`age` is bound to `10` and `name` is bound to `"Tom"`.

`tom.sayHi` is evaluated to `sayHi` in `Person` because outer environment of `tom` is `Person`.
Then, `sayHi()` is called in the environment of namespace `tom`, where `age` is bound to `10`, like method call.


```monkey:person.monkey
let Person = namespace {
  ...
  let grow = fn() {
        outer().new(age + 1, name)
  };
};
``` 

```
>> let tom = tom.grow()
>> tom
namespace {age: 11, name: Tom}
```

`outer()` can be used for setter-like functions.

Since all values in Monkey are immutable, the "setter" returns new "instance" with different values instead.

Since outer environment of `grow` is `Person`, `outer()` returns environment of `Person`.
`outer().new(age + 1, name)` works as `Person.new(age + 1, name)`.

Tom turned to be eleven!

Notice: you can use `self().new(age + 1, name)` instead, but not recommended by  performance reason.

If you use `self()`, `new` is evaluated inside the environment of namespace `tom`, then returns a namespace inner of `tom` (not inner of `Person`).
When you reassign `tom.grow()` to `tom` 3 times, `tom` is a namespace in namespace in namespace in namespace in `Person`. These redundant nests make access to functions in `Person` slower (and look strange...).

## Avoiding panic

### Empty block returns `*object.Null` instead of `nil`

```
>> let x = fn() {};
>> x()
null
```

### Check ality of function call 

Error messages are same as built-in functions'.

```
>> let add = fn(x, y) { x + y };
>> add(1)
ERROR: wrong number of arguments. got=1, want=2
```
