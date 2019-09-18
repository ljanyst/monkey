
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
