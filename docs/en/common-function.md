Standard Functions List
---
This page introduces a list of standard functions. Note that functions introduced in get-started.md are not included here.

| Function Form | Execution Content |
| --- | --- |
| input variable_name; | Assigns input to the variable_name. |
| make random variable_name type[int, float, array] value1 (value2); | Assigns a random number from value1 to value2 to the variable_name. However, if the type is array, it randomly selects from value1, so value2 can be omitted. |
| runmyself code; | Executes the code. |
| delete variable_name (variable_name2 ...); | Deletes the variable_name. Multiple variable names can be specified. (Optional) |
| swap variable_name1 variable_name2; | Swaps the values of variable_name1 and variable_name2. |

Package Functions List
---
Packages refer to functions other than standard functions. Below is a list of functions for each package. Note that although the list omits this detail, functions are actually called in the form of `package_name.function_name` (e.g., for the array package, it would be `array.function_name`).

### array Package
| Function Form | Execution Content |
| --- | --- |
| reset variable_name; | Empties the array named variable_name. |
| split variable_name1 array_or_variable_name delimiter; | Splits `array_or_variable_name` by `delimiter` and assigns it to `variable_name1`. |
| join variable_name1 array_or_variable_name delimiter; | Joins `array_or_variable_name` by `delimiter` and assigns it to `variable_name1`. |
| addbeg variable_name value; | Adds `value` to the beginning of the array named `variable_name`. |
| addend variable_name value; | Adds `value` to the end of the array named `variable_name`. |
| addnth variable_name index value; | Adds `value` to the array named `variable_name` at the `index` position, shifting the original value backward. |
| replace variable_name index value; | Replaces the value at the `index` position of the array named `variable_name` with `value`. |
| delnth variable_name index; | Deletes the value at the `index` position of the array named `variable_name`, shifting the original value forward. |
| sort variable_name; | Sorts the array named `variable_name` in ascending order. |
| reverse variable_name; | Reverses the order of the array named `variable_name`. |
| search variable_name1 variable_name2 value; | Searches for `value` in the array named `variable_name2` and assigns the index to `variable_name1` if found; otherwise, assigns -1. |

### string Package
| Function Form | Execution Content |
| --- | --- |
| replace variable_name1 value_or_variable_name before after count; | Replaces `before` with `after` in the string `value_or_variable_name` `count` times and assigns it to `variable_name1`. |
| addbeg variable_name1 value_or_variable_name; | Adds the string `value_or_variable_name` to the beginning of the string `variable_name1`. |
| addend variable_name1 value_or_variable_name; | Adds the string `value_or_variable_name` to the end of the string `variable_name1`. |
| include variable_name1 value_or_variable_name value; | Checks if `value` is included in the string `value_or_variable_name`, and assigns true to `variable_name1` if it is, otherwise assigns false. |
| substr variable_name1 value_or_variable_name start length; | Extracts a substring of `length` starting from `start` from the string `value_or_variable_name` and assigns it to `variable_name1`. |

### date Package
The date format for the date package is as follows: YYYY-MM-DD HH:mm:SS

`YYYY` represents the year, `MM` the month, `DD` the day, `HH` the hour, `mm` the minute, and `SS` the second.

| Function Form | Execution Content |
| --- | --- |
| now format variable_name format; | Formats the current time with `format` and assigns it to `variable_name`. |
| now year variable_name; | Assigns the current year to `variable_name`. |
| now month variable_name; | Assigns the current month to `variable_name`. |
| now day variable_name; | Assigns the current day to `variable_name`. |
| now hour variable_name; | Assigns the current hour to `variable_name`. |
| now minute variable_name; | Assigns the current minute to `variable_name`. |
| now second variable_name; | Assigns the current second to `variable_name`. |
| now dow variable_name; (now dayOfWeek variable_name;) | Assigns the current day of the week to `variable_name`. (0: Sunday, 1: Monday, 2: Tuesday, 3: Wednesday, 4: Thursday, 5: Friday, 6: Saturday) |
| calc add variable_name1 value_or_variable_name (year month day); | Adds `year` `month` `day` to the date formatted `value_or_variable_name` and assigns it to `variable_name1`. `(year month day)` is specified as an array and defaults to 0 if omitted. |
| calc sub variable_name1 value_or_variable_name (year month day); | Subtracts `year` `month` `day` from the date formatted `value_or_variable_name` and assigns it to `variable_name1`. `(year month day)` is specified as an array and defaults to 0 if omitted. |

### file Package
| Function Form | Execution Content |
| --- | --- |
| read filename value_or_variable_name; | Reads the `filename` file and assigns it to `value_or_variable_name`. |
| readline filename variable_name (line_number); | Reads the file named ``filename`` up to line_number and assigns it to ``variable_name`` as an array. ``(line_number)`` is optional, and if omitted, all lines will be read. |
| write filename value_or_variable_name; | Writes the value of `value_or_variable_name` to the `filename` file. |
| addend filename value_or_variable_name; | Adds the value of `value_or_variable_name` to the end of the `filename` file. |
| remove filename; | Deletes the `filename` file. |
| rename filename1 filename2; | Renames the file `filename1` to `filename2`. |

### get Package
| Function Form | Execution Content |
| --- | --- |
| local function_name (arg1 arg2 ...) filename; | Reads the `filename` file and defines the program within it as `function_name`. Arguments are optional, but `()` cannot be omitted. |
| url function_name (arg1 arg2 ...) URL; | Reads the program from the `URL` and defines the program within it as `function_name`. Arguments are optional, but `()` cannot be omitted. |
| github function_name (arg1 arg2 ...) "username/repository" filename; | Reads the `filename` within the repository `username/repository` on GitHub and defines the program within it as `function_name`. Arguments are optional, but `()` cannot be omitted. |

### regex Package
| Function Form | Execution Content |
| --- | --- |
| replace variable_name1 value_or_variable_name regex replacement_string; | Replaces parts of the string `value_or_variable_name` that match `regex` with `replacement_string` and assigns it to `variable_name1`. |
| find variable_name1 value_or_variable_name regex; | Searches for the string that matches `regex` in `value_or_variable_name` and assigns it to `variable_name1`. |
| findAll variable_name1 value_or_variable_name regex; | Searches for all strings that match `regex` in `value_or_variable_name` and assigns them to `variable_name1` as an array. |
