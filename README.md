MyLisp
=======

A _Scheme_ dialect.

I built this project to learn the core concepts of programming languages.
_MyLisp_ is a general purposed, strongly typed programming language, which
performs dynamic type checking just like other _Scheme_ dialects.

Currently support:
- a REPL user interface
- datatype: 64-bit integers, booleans, symbols, pairs, procedures & closures
- syntax: `define`, `lambda`, `cond`, `quote`

## Syntax

### Literals

``` scheme
123     ; integer
true    ; boolean
'abc    ; symbol
```

### Identifiers

keywords:
```
define  set!   lambda  cond   quote 
mod     and    or      not    eq?
equal?  cons   car     cdr    list
```

acceptable identifiers: `[_a-zA-z][_a-zA-z0-9]*[!\?]?`

``` scheme
a     ; variable a
b     ; variable b
foo   ; variable foo
```

### Procedure calls

basic procedure calls: 
`+` `-` `*` `/` `mod` `=` `<` `<=` `>` `>=` `and` `or` `not` `eq?` `equal?` `cons` `car` `cdr` `list`

``` scheme
(+ a b c)          ; add
(- a b c)          ; subtract
(* x y z)          ; multiply
(/ a b c)          ; divide by
(mod a b)          ; modulo operation
(= a b c)          ; equal
(< a b c)          ; less than
(<= a b c)         ; less than or equal
(> a b c)          ; greater than
(>= a b c)         ; greater than or equal
(and b1 b2 b3)     ; and
(or b1 b2 b3)      ; or
(not b)            ; not
(eq? a b)          ; equal in terms of memory
(equal? a b)       ; equal in terms of inherent value
(cons a b)         ; construct a pair
(car a b)          ; get the first item of a pair
(cdr a b)          ; get the second item of a pair
(list a b c d)     ; construct a list
```

regular procedure calls

``` scheme
(f x y z)
(foo)
(gcd 123 456)
(display (list 1 2 3))
```

### Definitions

variable definitions: start with keyword `define`

``` scheme
(define x 123)
(define a x)

(define I (lambda (x) x))
(define K (lambda (x) (lambda (y) x)))
(define S (lambda (x) (lambda (y) (lambda (z) ((x z) (y z))))))
```

### Assignments

variable assignment: start with keyword `set!`

``` scheme
(define x 123)
(set! x 456)

x    ; 456

(define p (cons 1 2))
(set! p (cons (car p) 3))

p    ; (1 . 3)
```

### Lambda expressions

``` scheme
(lambda () 123)           ; constant procedure
(lambda (x) x)            ; with 1 argument
(lambda (x y) (+ x y))    ; with 2 arguments
```

### Conditional expressions

``` scheme
(define a 4)
(define b 5)
(cond ((= a b) 9)
      ((> a b) 8)
      ((< a b) 7))     ; 7


(define gcd
 (lambda (a b)
  (cond ((= b 0) a)
        (else (gcd b (mod a b))))))

(gcd 65 13)      ; 13
(gcd 64 48)      ; 16
```

### Quoting and `'`

The quoting syntax is like `(quote expr)`, where `expr` is either atom or a list expression.
`'expr` is a shorthand for `(quote expr)`, and the former will be expanded into the latter during parsing.

When:
- `expr` is a primitive value, it returns the given value.
- `expr` is an identifier, it returns the symbol value to the identifier.
- `expr` is a list expression, it is like applying `quote` to each item of the list.

``` scheme
'1234            ; 1234
'a               ; a
'(1 2 3)         ; (1 2 3)
'(a b (c d))     ; (a b (c d))
'(a b 'c)        ; (a b (quote c))
```

## Implementation

An expression will go through the following processes after typed into the interpreter:

`lexer` => `parser` => `compiletime` => `runtime` => `printer`

### Lexer

The `lexer` converts input character stream into token streams.

### Parser

The `parser` converts token stream into list expressions.

### Compile time

The `compiletime` performs transformations on list expressions. It outputs AST nodes for the `runtime`.

### Runtime

The `runtime` evaluate AST nodes and outputs runtime values.

### Printer

The `printer` print runtime values to the console.

## Examples

Church numerals

```
(define n0 (lambda (f) (lambda (x) x)))
(define n1 (lambda (f) (lambda (x) (f x))))
(define show (lambda (n) ((n (lambda (x) (+ x 1))) 0)))
(define add (lambda (a b) (lambda (f) (lambda (x) ((a f) ((b f) x))))))
(define mul (lambda (a b) (lambda (f) (lambda (x) ((a (b f)) x)))))

(define n2 (add n1 n1))
(define n3 (add n1 n2))
(define n4 (add n2 n2))
(define n5 (add n2 n3))
(define n8 (add n3 n5))
(define n13 (add n5 n8))
(define n65 (mul n5 n13))
(define n32 (mul n4 n8))
(define n64 (mul n8 n8))
(define n1024 (mul n32 n32))

(show n4)
$ 4
(show n5)
$ 5
(show n13)
$ 13
(show n1024)
$ 1024
(show n64)
$ 64

(define n65536 (mul n64 n1024))
(show n65536)
$ 65536
```

Message-passing style OOP

```
(define NewProfile
 (lambda ()
  (define id 0)
  (define name 'name)
  (define setId (lambda (x) (set! id x)))
  (define setName (lambda (x) (set! name x)))
  (lambda (msg)
   (cond ((eq? msg 'Id) id)
         ((eq? msg 'SetId) setId)
         ((eq? msg 'Name) name)
         ((eq? msg 'SetName) setName)))))

(NewProfile)
$ <procedure>:proc

(define p (NewProfile))
p
$ <procedure p>

(p 'Id)
$ 0

((p 'SetId) 1)
(p 'Id)
$ 1

(p 'Name)
$ name

((p 'SetName) 'dy)
(p 'Name)
$ dy
```