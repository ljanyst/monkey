Monkey
======

Monkey is a programming language and an execution environment that I like to
play with. The language itself is not meant to be useful in any way, shape, or
form. The primary source of amusement, and perhaps value, is the implementation
of the interpreter (and at some future point the compiler) and the execution
environment. It's loosely based on (and very significantly extended) the
language described in [Writing An Interpreter In Go](https://interpreterbook.com).
However, since the book is meant for complete beginners both in compilers and in
go, I lacked the patience to read most of it. Therefore, the implementation is
for the most part my own.

Building
--------

Use the standard go worflow:

    go install github.com/ljanyst/monkey/cmd/monkey

Examples
--------

All the code below can be found in the `examples` subdirectory. To run in, just
say something like:

    monkey quicksort.monkey

### Closures ###

```
let adder = fn(x) {
  return fn(y) { return x + y; };
};

let multiplier = fn(x) {
  return fn(y) { return x * y; };
};

let compositor = fn(f1, f2) {
  return fn(x) { return f1(f2(x)); };
};

let result = compositor(adder(5), multiplier(2))(3);

print("Result: #", result);
```

### Quicksort ###

```
let string = "zażółć gęślą jaźń";

let swap = fn(array, a, b) {
  let c = array[a];
  array[a] = array[b];
  array[b] = c;
};

let pivot = fn(array, start, end) {
  let i = end;
  for (let j = end; j > start; j = j - 1) {
    if (array[j] > array[start]) {
      swap(array, j, i);
      i = i - 1;
    };
  };
  swap(array, start, i);
  return i;
};

let qs = fn(array, start, end) {
  if (end - start < 1) {
    return nil;
  };
  let p = pivot(array, start, end);
  qs(array, start, p-1);
  qs(array, p+1, end);
};

let quicksort = fn(string) {
  qs(string, 0, len(string)-1);
};

print("Unsorted: #", string);
quicksort(string);
print("Sorted:   #", string);
```

### The towers of hanoi ###

```
let a = {3, 2, 1};
let b = {};
let c = {};

let printState = fn(a, b, c) {
  print("a = #", a);
  print("b = #", b);
  print("c = #", c);
  print("----------------");
};

let move = fn(n, source, target, auxiliary) {
  if (n <= 0) {
    return nil;
  };

  move(n - 1, source, auxiliary, target);
  append(target, pop(source));
  printState(a, b, c);
  move(n - 1, auxiliary, target, source);
};

printState(a, b, c);
move(3, a, c, b);
```
