サンプルコード
---

### 九九表をゼロ埋めできれいに表示するプログラム
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

### 乱数が素数かどうかを判定するプログラム
```
sharp is_prime ($num) {
    match $num
    case (2 3 5 7) {
        return true;
    };

    if #($num == 1) {
        return false;
    } else {
        match $num
        case #(_ % 2 == 0) {
            return false;
        }
        case #(_ % 3 == 0) {
            return false;
        }
        case #(_ % 5 == 0) {
            return false;
        }
        case #(_ % 7 == 0) {
            return false;
        }
        default {
            return true;
        };
    };
};

printf ":1:が素数か: :2:\n" #var(x #rand(1 100)) #is_prime(x);
```

### 任意個のフィボナッチ数列を出力するプログラム
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
// 引数が大きすぎるとエラーが出るので注意;
```