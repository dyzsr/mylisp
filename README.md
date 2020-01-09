# mylisp
My implementation of the Scheme Language.

Currently support:
- datatypes: integers, booleans
- syntaxes: `define`, `lambda`, `cond`
- features: dynamic type checking, closures

## Demo

Basic operations

```
123
$ 123

true
$ true

(+ 1 (* 2 (- (/ 99 3) 4)))
$ 59

(define a 4)
(define b 5)
(cond ((= a b) 9)
      ((> a b) 8)
      ((< a b) 7))
$ 7
```

Procedures & closures

```
((lambda (x) (+ x x)) 32)
$ 64

(define f (lambda (x) (lambda (y) (+ x y))))
((f 10) 20)
$ 30

(define g (f 20))
(g 20)
$ 40
```

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