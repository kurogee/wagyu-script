Sharp functions
---
Sharp functions are a feature unique to this language. Because of the nature of functions in this language, normal functions cannot be included in arguments, so I created this function because I wanted to create a function that simply returns a value and passes a value to an argument.

Basic form
---
The basic form of a sharp function is as follows. A sharp function can be called by adding ``#`` before the function name.
```
#Function name(Argument 1 Argument 2 ...)
```
※ Spaces cannot be left between the function name and the argument. (We are currently considering allowing spaces.)

For example, if we take the function ``#(2 + 3)`` in the example of get-started.md, it will be in the following form.
```
println #(2 + 3);
```
By doing this, the calculation result of ``2 + 3``, that is, ``5``, will be entered in the argument of the ``println`` function. This is the basic usage of sharp functions.

List of standard functions
---
The variables in the examples in this explanation have the following values:
```
arr = array ("apple" "banana" "orange" "grape");
val = value "Wagyu";
num = value 10;
```
* If the return value of the execution result is a decimal, the significant digits may differ depending on the function.

| Function name | Execution content | Example |
|---|---|---|
| arrayAt(array name index) | Returns the value at any index of the array | ``#arrayAt(arr 2)`` → ``"orange"``|
| arrayLen(array name) | Returns the length of the array | ``#arrayLen(arr)`` → ``4`` |
| at(string index) | Returns the character at any index of the string | ``#at("Hello, World!" 7)`` → ``"W"`` |
| len(value or variable name) | Returns the length of the value | ``#len(val)`` → ``5`` |
| from(integer value 1 integer value 2) | Creates an array containing consecutive numbers ranging from ``integer value 1`` to ``integer value 2`` | ``#from(1 10)`` → ``(1 2 3 ... 9 10)`` |
| format(format variable 1 variable 2 ...) | Returns a string formatted like ``printf`` | ``#format("Hello, :1:." val)`` → ``"Hello, Wagyu."`` |
| remove(arr_name index) | Returns an array with the value at index removed | ``#remove(arr 1)`` → ``("apple" "orange" "grape")`` |
| rand(integer_value1 integer_value2) | Returns a random number in the range ``integer_value1`` to ``integer_value2`` | ``#rand(1 10)`` → ``(random number between 1 and 10)`` |
| convert(value or variable_name type_name) | Converts the type of a value | ``#convert(num float)`` → ``10.00`` |
| var(variable_name value) | Declare and initialize a variable, and enter the variable name in the argument with the function | ``println #var(x 10);`` → ``10`` (``10`` is entered into the variable ``x``, and ``10`` is output.) |
| repeat(string Repeat count) | Returns a string with a specified number of repetitions | ``#repeat("Hello, " 3)`` → ``"Hello, Hello, Hello, "`` |
| all(array or any number of arguments) | Returns ``true`` if all arguments are ``true`` or ``1`` | ``#all(true false true)`` → ``false`` |
| any(array or any number of arguments) | Returns ``true`` if there is at least one ``true`` or ``1`` in the arguments | ``#any(true false true)`` → ``true`` |

Create your own function
---
Just like normal functions, you can create your own sharp functions. Basically, write them as follows:
```
sharp function name (argument 1 argument 2 ...) {
    process...;
    return return value;
};
```
※ Leave a space between the function name and the argument.

This function can be called in the form ``#function name (argument 1 argument 2 ...)``, and the return value of ``return`` will be returned to wherever this function is entered.

The following is an example of a sharp function that returns the sum of two numbers entered as arguments.
```
sharp add ($a $b) {
    return #($a + $b);
};

println #add(2 3);
// -> ``5'' is output. ;
```

List of functions from other packages
---

### math package
The maths package performs mathematical operations. Although it is omitted in the ``Function name`` column of the table, the function format is ``#math.function name``.

The variables that appear in the examples in this explanation are assumed to have the following values.
```
arr = array (2 3 1 5 4);
```

| Function name | Action | Example |
|---|---|---|
| pi() | Returns pi | ``3.1415926535`` |
| e() | Returns Napier's constant | ``2.7182818284`` |
| cos(value) | Returns the cosine of a value | ``#math.cos(0)`` → ``1.00`` |
| sin(value) | Returns the sine of a value | ``#math.sin(0)`` → ``0.00`` |
| tan(value) | Returns the tangent of a value | ``#math.tan(0)`` → ``0.00`` |
| acos(value) | Returns the arccosine of a value | ``#math.acos(1)`` → ``0.00`` |
| asin(value) | Returns the arcsine of a value | ``#math.asin(0)`` → ``0.00`` |
| atan(value) | Returns the arctangent of a value | ``#math.atan(0)`` → ``0.00`` |
| sqrt(value) | Returns the square root of a value | ``#math.sqrt(4)`` → ``2.00`` |
| pow(value1 value2) | Returns the square of value1 | ``#math.pow(2 3)`` → ``8.00`` |
| log(value) | Returns the natural logarithm of a value | ``#math.log(2.7182818284)`` → ``1.00`` |
| log10(value) | Returns the base 10 logarithm of a value | ``#math.log10(10)`` → ``1.00`` |
| abs(value) | Returns the absolute value of a value | ``#math.abs(-1)`` → ``1.00`` |
| median(array) | Returns the median value of an array | ``#math.median(arr)`` → ``3.00`` |
| mode(array) | Returns the most frequent value in an array. If all occur the same number of times, return the first value. | ``#math.mode(arr)`` → ``2.00`` |
| average(arr) | Returns the average value of an array | ``#math.average(arr)`` → ``3.00`` |
| floor(value) | Returns the value rounded down to the nearest integer | ``#math.floor(3.14)`` → ``3.00`` |
| ceil(value) | Returns the value rounded up to the nearest integer | ``#math.ceil(3.14)`` → ``4.00`` |
| round(value number of decimal places) | Returns the value rounded to the specified number of decimal places | ``#math.round(3.14 1)`` → ``3.10`` |

### regex package
The regex package allows you to manipulate strings using regular expressions. The function format is #regex.function name, which is omitted in the ``Function Name`` column of the table.

The variables in this explanation example have the following values:
```
str = value "Hello, World!";
```

| Function name | Execution content | Example |
|---|---|---|
| match(regex value) | Returns whether the value matches the regular expression | ``#regex.match("Hello" str)`` → ``true`` |
| find(regex value) | Returns the part of the value that matches the regular expression | ``#regex.find("Hello" str)`` → ``"Hello"`` |

### date package
The date package performs operations related to dates. The function format is #date.function name, which is omitted in the ``Function Name`` column of the table.

| Function name | Action | Example |
|---|---|---|
| nowFull() | Returns the current date and time | ``2024-08-05 12:34:56`` |
| nowDate() | Returns the current date | ``2024-08-05`` |
| nowTime() | Returns the current time | ``12:34:56`` |
| nowYear() | Returns the current year | ``2024`` |
| nowMonth() | Returns the current month | ``8`` |
| nowDay() | Returns the current day | ``5`` |
| nowHour() | Returns the current hour | ``12`` |
| nowMinute() | Returns the current minute | ``34`` |
| nowSecond() | Returns the current second | ``56`` |
| nowDow (nowDayOfWeek) | Returns the current day of the week | (0:Sunday 1:Monday 2:Tuesday 3:Wednesday 4: Thursday 5: Friday 6: Saturday) |
| nowUnix() | Returns the current UNIX time | ``1234567890`` |

Examples
---
### Multiplication tables using Sharp functions effectively
```
x = array #from(1 9);
y = array #from(1 9);

each y > i {
    each x > j {
        res = value #(i * j);
        printf ":1::2: " #repeat("0" #(2 - #len(res))) res;
    };
    println "";
};
```

### Prime number determination using Sharp functions (for advanced users)
```
sharp is_prime ($num) {
    match $num
    case (2 3 5 7) {
        return "true";
    };
    if #($num == 1) {
        return "false";
    } else {
        match $num
        case #(_ % 2 == 0) {
            return "false";
        }
        case #(_ % 3 == 0) {
            return "false";
        }
        case #(_ % 5 == 0) {
            return "false";
        }
        case #(_ % 7 == 0) {
            return "false";
        }
        default {
            return "true";
        };
    };
};

println #is_prime(89);
// -> "true" is output.;
```