
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
