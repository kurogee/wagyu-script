Wagyu-Script
---
Wagyu-scriptは、@kurogeが開発した半角スペース区切りが特徴のシンプルなプログラミング言語です。

始め方など
---
詳しくは、[Get Started](./docs/get-started.md)ページをご覧ください。

Examples
---

### Hello, World!
```
println "Hello, World!";
```

### FizzBuzz
```
start = value 1;
end = value 100;

while #(start <= end) {
    if #(start % 15 == 0) {
        println "FizzBuzz";
    } elif #(start % 5 == 0) {
        println "Buzz";
    } elif #(start % 3 == 0) {
        println "Fizz";
    } else {
        println start;
    };

    start = add 1;
};
```