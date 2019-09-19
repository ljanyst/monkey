
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
