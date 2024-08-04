package regex

import (
	"fmt"
	"strings"
	"regexp"
)

func variables_replacer(variables map[string]string, target string) string {
	val, ok := variables[target]
	if ok {
		return val
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

func Run(name string, value []string, variables *map[string]string) {
	if name == "replace" {
		// value[0] = 結果を入れる変数名 value[1] = 対象の文字列 value[2] = 置換前のパターン value[3] = 置換後の文字列
		// 置き換え回数が無い場合はすべて置換
		if len(value) == 4 {
			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in replace function in regex package. [1]")
			}

			target := variables_replacer(*variables, value[1])
			before := regexp.MustCompile(variables_replacer(*variables, value[2]))
			after := variables_replacer(*variables, value[3])

			result := before.ReplaceAllString(target, after)
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

			target := variables_replacer(*variables, value[1])
			pattern := variables_replacer(*variables, value[2])

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

			target := variables_replacer(*variables, value[1])
			pattern := variables_replacer(*variables, value[2])

			result := regexp.MustCompile(pattern).FindAllString(target, -1)
			// 結果を空白でjoinして代入
			(*variables)[value[0]] = strings.Join(result, " ")
		} else {
			fmt.Println("The error occurred in findAll function in regex package. [2]")
		}
	}
}

func Sharp(func_name string, args []string, variables *map[string]string) string {
	if func_name == "match" {
		// args[0] = パターン args[1] = 対象の文字列 -> true or false
		if len(args) == 2 {
			pattern := variables_replacer(*variables, args[0])
			target := variables_replacer(*variables, args[1])

			matched, err := regexp.MatchString(pattern, target)
			if err != nil {
				fmt.Println("The error occurred in match function in regex package. [1]")
				fmt.Println(err)
			}

			if matched {
				return "true"
			} else {
				return "false"
			}
		} else {
			fmt.Println("The error occurred in match function in regex package. [2]")
		}
	}

	return ""
}
