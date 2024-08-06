Get Started
---
In "Get Started," you'll learn about this language, set it up, and get to know it easily.

What is this language?
---
This language is a programming language called "Wagyu-script," which is characterized by half-width space delimiters. The name comes from my online username, "@kuroge."

How to get started with this language
---
There is an item called "Releases" on the left side of the repository page, so please download the latest version for each OS from there. No installation is required.

After downloading, unzip the file and run it as follows. (The file extension is ``.wg``, but if it's a text file, ``.txt`` is fine.)

GNU/Linux
```sh
./wagyu-script run filename.wg
```
Windows
```cmd
./wagyu-script.exe run filename.wg
```

Hello, World!
---
This language is easier to write than C or Java. Below is the most basic program that outputs ``Hello, World!`` to the prompt.
```
println "Hello, World!";
```
General output (adding a newline to the end of the line) can be done with the ``println`` function. Basically, arguments are separated by half-width spaces and follow the function name. The ``"~"`` argument is called a **string type**. It is surrounded by quotation marks. It can also be expressed as ``'~'`` using single quotation marks.

When using regular expressions, you can also enclose a string in `` `〜` `` (backquotes).

This language also requires a ``;`` (semicolon) at the end of each line. The meaning of the program does not change whether or not there is a line break for each function, so it is possible to write it all on one line.
```
println "Hello!";println "World!";
```
※ Variables can also be entered as arguments to ``println``.

Comments
---
Comments are where you can freely write notes about the program. The contents of the comment do not affect the operation.

```
// Content of comment;
```
Please note that comments also require ``;`` at the end of the line.

Types
---
This language has functions to specify types, but they are rarely used because they are mostly automatically interpreted. Below is a list of types.

| Name | Type name | Contents to be put into the type |
|---|---|---|
| String type | string | String (ex: ``"Hello!"``) |
| Array type | array | Array (ex: ``("Hello" "World" "Wagyu")``) |
| Integer type | int | Integer value (ex: ``10``) |
| Decimal type | float | Decimal value (ex: ``3.14``)|

Other output
---
In addition to ``println``, the following output functions are available.

- ``printf`` Insert variables into a string according to the format and output (without a newline at the end of the line).
    ```
    printf "format" variable 1 variable 2 ...;
    ```
    In the format, for example, if you enter ``:1:``, the first variable name added to the argument will be inserted there. Below is an example program using this.
    
    ※ In this example, the variable ``x`` contains the string type value ``Tanaka``.
    ```
    printf "Hello, :1:!" x;
    // -> "Hello, Tanaka!" is output.
    ```
- ``print`` Outputs a string without a newline at the end of the sentence.
    ```
    print string or variable name;
    ```

Variable
---
There are two ways to declare and initialize variables.

1. Declare and initialize variables

- For normal values

    ``Variable name = value value;``

- For arrays

    ``Variable name = array (value1 value2 value3 ...);``

- Can be omitted

    ``Variable name = value;``

- When declaring multiple variables

    ``vars (variable name1 variable name2 ...) (value1 value2 ...);``

2. Declare only empty variables

- When there is only one

    ``make var variable name;``

- When there are two or more

    ``make var variable name1 variable name2 ...;``

※ All values ​​are assigned as strings and are reverted to numeric types each time a calculation is performed.

※ When referencing a variable, if you include the variable name in the argument, it will automatically be reverted to the referenced value.

If you simply want to increase/decrease the value of a variable, you can also use ``add``.
```
Variable name = add number (negative numbers are supported);
```
The following is an example of a program that uses add to add 1 to the value of the variable ``x``.
```
x = value 1;
x = add 1;

println x;
// -> "2" is output;
```

Calculation
---
If you want to put the results of calculations into variable assignments or function arguments, you can use the unique feature **Sharp Functions** called ``#(Calculation formula)``. For more information on this function, see [Sharp Functions](./sharp-function.md). (We recommend reading this after reading through these.)

For example, if you want to put the calculation result of ``1 + 1`` into the variable ``x``, do it like this.
```
x = value #(1 + 1);
```
You can also put variable names. The operator, value, and variable name must each be separated by a half-width space.
```
y = value #(x+1);
// An example that will cause an error. Be sure to separate them with a half-width space;
```

Conditional branching
---
Conditional branching (if statements) in this language is written as follows. There are three ways to write it.

1. When writing a program that is executed only when the value is true
```
if condition {
    Execution content;
    ...
};
```
2. When writing a program that is executed when the value is true and when the value is false
```
if condition {
    Execution content when the value is true;
    ...
} else {
    Execution content when the value is false;
    ...
};
```
3. When setting up conditional branches (else-if) for when the value is true, false, and also when the value is false
```
if condition {
    Execution content when the value is true;
} elif condition2 {
    Execution content when the value is true;
} else {
    Execution content when the value is false;
};
```

Conditional statements are judged by the values ​​of the calculation results, ``true`` and ``false``, or ``1`` (true) and ``0`` (false). Also, the conditional statements can be written using ``#(～)`` as mentioned earlier.

The following is a program that prints ``Good!`` if the variable ``x`` is 1, and ``Oh, no...`` if it is not.
```
if #(x == 1) {
    println "Good!";
} else {
    println "Oh, no...";
};
```

Match statement
---
This language provides a function called ``match`` because ``if`` is too troublesome to use when you just want to check if numbers or letters match.'' ``match`` allows you to write the content that will be executed when the condition is met.

``match`` is written as follows.
```
match variable name
case value 1 {
    content;
    ...
}
case value 2 {
    content;
    ...
}
default {
    content;
    ...
};
```
``case`` writes the value that matches the condition. If the value matches the ``variable name`` value, the contents of that ``case`` will be executed. Also, if there are multiple cases that match the condition, they will all be executed, so use ``break`` appropriately.

If you use a sharp function, it will be executed if the result is ``true`` or ``1``. In that case, you can substitute ``_`` for the ``variable name`` in the sharp function argument.
```
x = value 1;

match x
case #(_ == 1) { // In reality, it would be better to write 1 normally for this level of difficulty, but for the sake of example;
    println "Good!";
}
case #(_ == 2) {
    println "Nice!";
};
```
The ``_`` part will be replaced with the value of the ``x`` variable, that is, ``1``.

``default`` is the content to be executed if none of the ``case``s match. It must be written at the very end of the ``match`` statement.

The following is an example of a program that prints ``Good!`` when the variable ``x`` is 1, ``Nice!`` when it is 2, and ``Oh, no...`` otherwise.
```
match x
case 1 {
    println "Good!";
}
case 2 {
    println "Nice!";
}
default {
    println "Oh, no...";
};
```

Repeat
---
This language provides a repeat function called ``while``.

``while`` allows you to repeat processing while the condition is true or 1. The basic form is as follows.

```
while condition {
    execution content;
    ...
};
```

The following is an example of a program that repeats the contents of ``while`` five times until the variable ``x`` becomes 5.
```
x = value 1;
while #(x <= 5) {
    println "Hello!";
    x = add 1;
};
```

The ``each`` function is provided as an **array repetition** function. ``each`` allows you to retrieve elements of an array one by one and perform repeated processing.

``each`` is written as follows.
```
each array name > variable name {
    Execution content;
    ...
};
```
For ``variable name``, specify the name of the variable to which the value will be put when the elements of the array are retrieved one by one. The variable is created automatically and is only valid within that block.

Below is an example of a program that retrieves and outputs the elements of the array ``arr`` one by one.
```
arr = array ("Good" "Nice" "Great");
each arr > i {
    printf ":1:!\n" i;
};
```
When this is executed, ``Good!``, ``Nice!``, and ``Great!`` are output.

To get the index of an array, change ``>`` to ``:``. Here is an example.
```
arr = array ("Good" "Nice" "Great");
each arr : i {
    printf ":1:\n" i;
};
```
When you run this, the index of ``arr``, that is, ``0``, ``1``, and ``2``, are printed.

The ``break`` and ``continue`` functions are also provided to interrupt and restart the iteration.

The ``break`` function is written as follows, and when executed, it interrupts the iteration and exits the loop.
```
break "";
```

The ``continue`` function is written as follows, and when executed, it returns to the beginning of the iteration and restarts it.
```
continue "";
```

Functions
---
In this language, ``fnc`` is used to declare functions. By declaring a function, you can group various processes into one and call them any number of times.

Function declarations are written as follows:
```
fnc function name (argument 1 argument 2 ... (optional)) {
    function execution;
    ...
};
```

Function calls are written as follows:
```
function name (argument 1 argument 2 ... (write if there is));
```

If there is a return value (``return``), add the argument ``to variable name``.
```
function name (argument) to variable name to put return value into;
```

To return a return value from a function, use ``return value;``. When it is executed, the return value is returned and the function processing ends.

The following is an example of a program that uses a function to input a person's name and output a greeting.
```
fnc say_hello ($name) {
    printf "Hello, :1:!\n" $name;
};
say_hello ("Tanaka");
say_hello ("Yamada");
// -> "Hello, Tanaka!" and "Hello, Yamada!" are output. ;
```

Below is an example of a program that uses a function to return the sum of two numbers.
```
fnc add ($x $y) {
return #($x + $y);
};
make var result;
add (1 2) to result;
println result;
// -> "3" is output;
```

About other functions
---
We have explained basic functions and features.
For other functions, please refer to the following page
→ [List of standard functions](./common-function.md)

Sample code
---
Various sample codes are available. Please refer to them.
→ [Sample code](./sample-code.md)