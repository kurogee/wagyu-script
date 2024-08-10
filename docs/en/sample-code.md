Sample Code
---

### A program that neatly displays multiplication tables by filling in zeros
```
x = array #from(1 9);
y = array #from(1 9);

each y > i {
    each x > j {
        res = value #(i * j);
        printf ":1::2:  " #repeat("0" #(2 - #len(res))) res;
    };
    println "";
};
```

### A program to determine whether an arbitrary number is a prime number
```
sharp is_prime ($n) {
    if #($n < 2) { return false; };
    if #($n == 2) { return true };
    if #($n % 2 == 0) { return false };

    mem = #math.sqrt($n);
    i = 3;
    while #(i <= mem) {
        if #($n % i == 0) { return false };
        i = add 2;
    };

    return true;
};

println #is_prime(89); // true;
```

### A program to print any number of Fibonacci numbers
```
fnc fib ($n) {
    a = 0;
    b = 1;

    count = #from(0 $n);
    each count : c {
        println a;

        tmp = a;
        a = b;
        b = #(tmp + b);
    };
};

fib (20);
// Note that if the argument is too large, an error will occur.;
```