package string

import (
	"fmt"
	"strings"
	"strconv"

	system_split "github.com/kurogee/wagyu-script/system_split"
)

var take_off_quotation = system_split.Take_off_quotation

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

func Run(name string, value []string, value_in_quotes []bool, variables *map[string]string) {
	if name == "replace" {
		// value[0] = 変数名 value[1] = 対象の文字 value[2] = 置換前の文字列 value[3] = 置換後の文字列 value[4] = 置き換え回数
		if len(value) == 5 {
			replace_count, err := strconv.Atoi(variables_replacer(variables, value[4], value_in_quotes[4], false))
			if err != nil {
				fmt.Println("The error occurred in replace function in string package. [1]")
				fmt.Println(err)
			}

			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in replace function in string package. [2]")
				fmt.Println("The variable is not found.")
			}

			for i, v := range(value[1:3]) {
				value[i + 1] = variables_replacer(variables, v, value_in_quotes[i + 1], false)
			}

			(*variables)[value[0]] = strings.Replace(value[1], value[2], value[3], replace_count)
		} else {
			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in replace function in string package. [3]")
			}

			for i, v := range(value[1:]) {
				value[i + 1] = variables_replacer(variables, v, value_in_quotes[i + 1], false)
			}

			(*variables)[value[0]] = strings.Replace(value[1], value[2], value[3], -1)
		}
	} else if name == "addend" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addend function in string package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(variables, value[1], value_in_quotes[1], true)

		(*variables)[value[0]] = (*variables)[value[0]] + val
	} else if name == "addbeg" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addbeg function in string package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(variables, value[1], value_in_quotes[1], true)

		(*variables)[value[0]] = val + (*variables)[value[0]]
	} else if name == "include" {
		// value[0] = 結果を入れる変数名 value[1] = 対象の文字 value[2] = 検索する文字列
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in include function in string package. [1]")
			fmt.Println("The variable is not found.")
		}

		for i, v := range(value[1:]) {
			value[i + 1] = variables_replacer(variables, v, value_in_quotes[i + 1], false)
		}

		if strings.Contains(value[1], value[2]) {
			(*variables)[value[0]] = "true"
		} else {
			(*variables)[value[0]] = "false"
		}
	} else if name == "substr" {
		// value[0] = 変数 [1] = 切り出したい文字 [2] = 切り出すはじめのインデックス [3] = 切り出す文字数(長さ)
		start, err := strconv.Atoi(variables_replacer(variables, value[2], value_in_quotes[2], false))
		if err != nil {
			fmt.Println("The error occurred in substr function in string package. [1]")
			fmt.Println(err)
		}

		length, err := strconv.Atoi(variables_replacer(variables, value[3], value_in_quotes[3], false))
		if err != nil {
			fmt.Println("The error occurred in substr function in string package. [2]")
			fmt.Println(err)
		}

		str := variables_replacer(variables, value[1], value_in_quotes[1], false)
		if start < 0 || start >= len(str) || length < 0 || start+length > len(str) {
			fmt.Println("The error occurred in substr function in string package. [3]")
			fmt.Println("Invalid start or length value.")
		}

		(*variables)[value[0]] = str[start : start+length]
	}
}

func Sharp(func_name string, args []string, args_in_quote []bool, variables *map[string]string) (string, bool) {
	if func_name == "replace" {
		if len(args) == 4 {
			// args[0] = 対象の文字 args[1] = 置換前の文字列 args[2] = 置換後の文字列 (args[3] = 置き換え回数)
			replace_count, err := strconv.Atoi(variables_replacer(variables, args[3], args_in_quote[3], false))
			if err != nil {
				fmt.Println("The error occurred in replace sharp function in string package. [1]")
				fmt.Println(err)
			}

			for i, v := range(args) {
				args[i] = variables_replacer(variables, v, args_in_quote[i], false)
			}

			return strings.Replace(args[0], args[1], args[2], replace_count), true
		} else {
			for i, v := range(args) {
				args[i] = variables_replacer(variables, v, args_in_quote[i], false)
			}

			return strings.Replace(args[0], args[1], args[2], -1), true
		}
	} else if func_name == "include" {
		// args[0] = 対象の文字 args[1] = 検索する文字列
		for i, v := range(args) {
			args[i] = variables_replacer(variables, v, args_in_quote[i], false)
		}

		if strings.Contains(args[0], args[1]) {
			return "true", false
		} else {
			return "false", false
		}
	} else if func_name == "substr" {
		// args[0] = 切り出したい文字 args[1] = 切り出すはじめのインデックス args[2] = 切り出す文字数(長さ)
		start, err := strconv.Atoi(variables_replacer(variables, args[1], args_in_quote[1], false))
		if err != nil {
			fmt.Println("The error occurred in substr sharp function in string package. [1]")
			fmt.Println(err)
		}

		length, err := strconv.Atoi(variables_replacer(variables, args[2], args_in_quote[2], false))
		if err != nil {
			fmt.Println("The error occurred in substr sharp function in string package. [2]")
			fmt.Println(err)
		}

		str := variables_replacer(variables, args[0], args_in_quote[0], false)
		if start < 0 || start >= len(str) || length < 0 || start+length > len(str) {
			fmt.Println("The error occurred in substr sharp function in string package. [3]")
			fmt.Println("Invalid start or length value.")
		}

		return str[start : start+length], true
	}
	
	return "", false
}