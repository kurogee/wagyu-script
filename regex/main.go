package regex

import (
	"fmt"
	"strings"
	"strconv"
	"regexp"
)

func take_off_quotation(target string) string {
	if strings.HasPrefix(target, "'") && strings.HasSuffix(target, "'") {
		return strings.Trim(target, "'")
	} else if strings.HasPrefix(target, "\"") && strings.HasSuffix(target, "\"") {
		return strings.Trim(target, "\"")
	} else if strings.HasPrefix(target, "`") && strings.HasSuffix(target, "`") {
		return strings.Trim(target, "`")
	}

	return target
}

func variables_replacer(variables *map[string]string, target string, target_in_quote, add_quotes bool) string {
	if target_in_quote {
		if add_quotes {
			return target
		}

		// 両端のクオートを取り除く
		target = take_off_quotation(target)
		return target
	}

	val, ok := (*variables)[target]
	if ok {
		// もし数値や配列ではなかったらクオートで囲む
		if add_quotes {
			if _, err := strconv.Atoi(val); err != nil {
				if _, err := strconv.ParseFloat(val, 64); err != nil {
					if len(strings.Split(val, " ")) == 1 {
						return "\"" + val + "\""
					}
				}
			}

			return val
		} else {
			return val
		}
	}

	return target
}

/*
func take_off_quotation(target string) string {
	if strings.HasPrefix(target, "'") && strings.HasSuffix(target, "'") {
		return strings.Trim(target, "'")
	} else if strings.HasPrefix(target, "\"") && strings.HasSuffix(target, "\"") {
		return strings.Trim(target, "\"")
	}

	return target
}
*/

func Run(name string, value []string, value_in_quotes []bool, variables *map[string]string) {
	if name == "replace" {
		// value[0] = 結果を入れる変数名 value[1] = 対象の文字列 value[2] = 置換前のパターン value[3] = 置換後の文字列
		if len(value) == 4 {
			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in replace function in regex package. [1]")
			}

			target := variables_replacer(variables, value[1], value_in_quotes[1], false)
			pattern := variables_replacer(variables, value[2], value_in_quotes[2], false)
			replacement := variables_replacer(variables, value[3], value_in_quotes[3], false)

			result := regexp.MustCompile(pattern).ReplaceAllString(target, replacement)
			(*variables)[value[0]] = result
		} else {
			fmt.Println("The error occurred in replace function in regex package. [2]")
		}
	} else if name == "find" {
		// value[0] = 結果を入れる変数名 value[1] = 対象の文字列 value[2] = パターン
		if len(value) == 3 {
			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in find function in regex package. [1]")
			}

			target := variables_replacer(variables, value[1], value_in_quotes[1], false)
			pattern := variables_replacer(variables, value[2], value_in_quotes[2], false)

			result := regexp.MustCompile(pattern).FindString(target)
			(*variables)[value[0]] = result
		} else {
			fmt.Println("The error occurred in find function in regex package. [2]")
		}
	} else if name == "findAll" {
		// value[0] = 結果を入れる変数名 value[1] = 対象の文字列 value[2] = パターン
		if len(value) == 3 {
			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in findAll function in regex package. [1]")
			}

			target := variables_replacer(variables, value[1], value_in_quotes[1], false)
			pattern := variables_replacer(variables, value[2], value_in_quotes[2], false)

			result := regexp.MustCompile(pattern).FindAllString(target, -1)
			// 結果を空白でjoinして代入
			(*variables)[value[0]] = strings.Join(result, " ")
		} else {
			fmt.Println("The error occurred in findAll function in regex package. [2]")
		}
	}
}

func Sharp(func_name string, args []string, args_in_quote []bool, variables *map[string]string) (string, bool) {
	if func_name == "match" {
		// args[0] = パターン args[1] = 対象の文字列 -> true or false
		if len(args) == 2 {
			pattern := variables_replacer(variables, args[0], args_in_quote[0], false)
			target := variables_replacer(variables, args[1], args_in_quote[1], false)

			matched, err := regexp.MatchString(pattern, target)
			if err != nil {
				fmt.Println("The error occurred in match function in regex package. [1]")
				fmt.Println(err)
			}

			if matched {
				return "true", false
			} else {
				return "false", false
			}
		} else {
			fmt.Println("The error occurred in match function in regex package. [2]")
		}
	} else if func_name == "find" {
		// args[0] = パターン args[1] = 対象の文字列 -> マッチした文字列 or 空文字列
		if len(args) == 2 {
			pattern := variables_replacer(variables, args[0], args_in_quote[0], false)
			target := variables_replacer(variables, args[1], args_in_quote[1], false)

			result := regexp.MustCompile(pattern).FindString(target)
			return result, true
		} else {
			fmt.Println("The error occurred in find function in regex package. [1]")
		}
	} else {
		fmt.Println("The error occurred in find function in regex package. [2]")
	}

	return "", false
}
