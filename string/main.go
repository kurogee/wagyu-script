package string

import (
	"fmt"
	"strings"
	"strconv"
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

func variables_replacers(variables map[string]string, sentence string, targets []string) string {
	var result string = sentence
	var count int = 1
	for _, target := range(targets) {
		val, ok := variables[target]
		if ok {
			result = strings.ReplaceAll(result, ":" + strconv.Itoa(count) + ":", val)
			count++
		} else {
			result = strings.ReplaceAll(result, ":" + strconv.Itoa(count) + ":", target)
			count++
		}
	}

	return result
}*/

func Run(name string, value []string, variables *map[string]string) {
	if name == "replace" {
		// value[0] = 変数名 value[1] = 対象の文字 value[1] = 置換前の文字列 value[2] = 置換後の文字列 value[3] = 置き換え回数
		if len(value) == 4 {
			replace_count, err := strconv.Atoi(variables_replacer(*variables, value[3]))
			if err != nil {
				fmt.Println("The error occurred in replace function in string package. [1]")
				fmt.Println(err)
			}

			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in replace function in string package. [2]")
			}

			for i, v := range(value[1:]) {
				value[i + 1] = variables_replacer(*variables, v)
			}

			(*variables)[value[0]] = strings.Replace(value[1], value[2], value[3], replace_count)
		} else {
			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in replace function in string package. [3]")
			}

			for i, v := range(value[1:]) {
				value[i + 1] = variables_replacer(*variables, v)
			}

			(*variables)[value[0]] = strings.Replace(value[1], value[2], value[3], -1)
		}
	} else if name == "addend" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addend function in string package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(*variables, value[1])

		(*variables)[value[0]] = (*variables)[value[0]] + val
	} else if name == "addbeg" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addbeg function in string package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(*variables, value[1])

		(*variables)[value[0]] = val + (*variables)[value[0]]
	} else if name == "include" {
		// value[0] = 結果を入れる変数名 value[1] = 対象の文字 value[2] = 検索する文字列
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in include function in string package. [1]")
			fmt.Println("The variable is not found.")
		}

		for i, v := range(value[1:]) {
			value[i + 1] = variables_replacer(*variables, v)
		}

		if strings.Contains(value[1], value[2]) {
			(*variables)[value[0]] = "true"
		} else {
			(*variables)[value[0]] = "false"
		}
	}
}
