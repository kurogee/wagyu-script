package array

import (
	"fmt"
	"strings"
	"strconv"

	system_split "github.com/kurogee/wagyu-script/system_split"
)

var take_off_quotation = system_split.Take_off_quotation
var split = system_split.Split
var divide_split = system_split.Divide_split

func variables_replacer(variables *map[string]string, target string, target_in_quote, add_quotes bool) string {
	if target_in_quote {
		if add_quotes {
			return "\"" + target + "\""
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
	// 基本的にvalue[0]は変数名
	if name == "reset" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in reset function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		(*variables)[value[0]] = ""
	} else if name == "split" {
		// value[0] = 変数 value[1] = 対象の文字列 value[2] = 区切り文字
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in split function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		value[1] = variables_replacer(variables, value[1], value_in_quotes[1], false)
		value[2] = variables_replacer(variables, value[2], value_in_quotes[2], false)

		splited_args := strings.Split(value[1], value[2])

		for i, val := range(splited_args) {
			splited_args[i] = "\"" + val + "\""
		}

		(*variables)[value[0]] = strings.Join(splited_args, " ")
	} else if name == "join" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in join function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		(*variables)[value[0]] = strings.Join(strings.Split(variables_replacer(variables, value[1], value_in_quotes[1], false), " "), value[2])
	} else if name == "addbeg" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addbeg function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(variables, value[1], value_in_quotes[1], true)
		// val = add_quotation(val)

		// もしvalue[0]の変数が空文字列なら、スペースを追加しない
		if (*variables)[value[0]] == "" {
			(*variables)[value[0]] = val
		} else {
			(*variables)[value[0]] = val + " " + (*variables)[value[0]]
		}
	} else if name == "addend" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addend function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(variables, value[1], value_in_quotes[1], true)

		// もしvalue[0]の変数が空文字列なら、スペースを追加しない
		if (*variables)[value[0]] == "" {
			(*variables)[value[0]] = val
		} else {
			(*variables)[value[0]] = (*variables)[value[0]] + " " + val
		}
	} else if name == "addnth" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in add function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		index, err := strconv.Atoi(variables_replacer(variables, value[1], value_in_quotes[1], false))
		if err != nil {
			fmt.Println("The error occurred in add function in array package. [2]")
			fmt.Println("The index is not integer.")
		}

		val := variables_replacer(variables, value[2], value_in_quotes[2], true)
		// val = add_quotation(val)

		(*variables)[value[0]] = (*variables)[value[0]][:index + 1] + " " + val + (*variables)[value[0]][index + 1:]
	} else if name == "replace" {
		// value[0] = target, value[1] = old index, value[2] = new value
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in replace function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		index, err := strconv.Atoi(variables_replacer(variables, value[1], value_in_quotes[1], false))
		if err != nil {
			fmt.Println("The error occurred in replace function in array package. [2]")
			fmt.Println("The index is not integer.")
		}

		// 一回配列に変換してから置換する
		slice, _ := divide_split(split((*variables)[value[0]]))

		val := variables_replacer(variables, value[2], value_in_quotes[2], false)
		// val = add_quotation(val)

		slice[index] = val
		(*variables)[value[0]] = strings.Join(slice, " ")
	} else if name == "delnth" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in delnth function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		index, err := strconv.Atoi(variables_replacer(variables, value[1], value_in_quotes[1], false))
		if err != nil {
			fmt.Println("The error occurred in delnth function in array package. [2]")
			fmt.Println("The index is not integer.")
		}

		// 一回配列に変換してから削除する
		slice, _ := divide_split(split((*variables)[value[0]]))
		slice = append(slice[:index], slice[index + 1:]...)

		(*variables)[value[0]] = strings.Join(slice, " ")
	} else if name == "sort" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in sort function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		// 一回配列に変換してからソートする
		slice, _ := divide_split(split((*variables)[value[0]]))
		for i := 0; i < len(slice); i++ {
			for j := i + 1; j < len(slice); j++ {
				if slice[i] > slice[j] {
					slice[i], slice[j] = slice[j], slice[i]
				}
			}
		}

		(*variables)[value[0]] = strings.Join(slice, " ")
	} else if name == "reverse" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in reverse function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		// 一回配列に変換してから逆順にする
		slice, _ := divide_split(split((*variables)[value[0]]))
		for i := 0; i < len(slice) / 2; i++ {
			slice[i], slice[len(slice) - i - 1] = slice[len(slice) - i - 1], slice[i]
		}

		(*variables)[value[0]] = strings.Join(slice, " ")
	} else if name == "search" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in search function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		_, ok2 := (*variables)[value[1]]
		if !ok2 {
			fmt.Println("The error occurred in search function in array package. [2]")
			fmt.Println("The variable is not found.")
		}

		// 一回配列に戻す
		slice, _ := divide_split(split((*variables)[value[1]]))

		// 配列から値を検索し、そのインデックスをvalue[0]の変数に格納 なければ-1
		index := -1
		for i, val := range slice {
			if val == variables_replacer(variables, value[2], value_in_quotes[2], false) {
				index = i
				break
			}
		}

		(*variables)[value[0]] = strconv.Itoa(index)
	}
}

func Sharp(func_name string, args []string, args_in_quote []bool, variables *map[string]string) (string, bool) {
	if func_name == "search" {
		// args[0] = 対象の配列 args[1] = 検索する文字列 -> インデックスを返す(なければ-1)

		args[0] = variables_replacer(variables, args[0], args_in_quote[0], false)
		args[1] = variables_replacer(variables, args[1], args_in_quote[1], false)

		// 配列に変換
		slice, _ := divide_split(split(args[0]))

		// 配列から値を検索し、そのインデックスを返す なければ-1
		index := -1
		for i, val := range(slice) {
			if val == args[1] {
				index = i
				break
			}
		}

		return strconv.Itoa(index), false
	} else if func_name == "split" {
		// args[0] = 変数 args[1] = 区切り文字
		args[0] = variables_replacer(variables, args[0], args_in_quote[0], false)
		args[1] = variables_replacer(variables, args[1], args_in_quote[1], false)

		splited_args := strings.Split(args[0], args[1])

		// 変数が入っている場合は、変数を展開する
		for i, val := range(splited_args) {
			splited_args[i] = "\"" + val + "\""
		}

		return "(" + strings.Join(splited_args, " ") + ")", false
	} else if func_name == "join" {
		// args[0] = 値or変数 args[1] = 区切り文字
		args[0] = variables_replacer(variables, args[0], args_in_quote[0], false)
		args[1] = variables_replacer(variables, args[1], args_in_quote[1], false)

		// args[0]に直接配列が入っている場合は、両端の括弧を取り除く
		if strings.HasPrefix(args[0], "(") && strings.HasSuffix(args[0], ")") {
			args[0] = args[0][1:len(args[0]) - 1]
		}

		divided_arg, divided_arg_in_quote := divide_split(split(args[0]))

		// 変数が入っている場合は、変数を展開する
		for i, val := range(divided_arg) {
			divided_arg[i] = variables_replacer(variables, val, divided_arg_in_quote[i], false)
		}

		// join関数を実行
		result := strings.Join(divided_arg, args[1])

		// divided_tokensを返す
		return result, true
	} else if func_name == "same" {
		// args[0] = 変数 args[1] = 変数
		// 配列の中身が同じかどうかを判定する
		args[0] = variables_replacer(variables, args[0], args_in_quote[0], false)
		args[1] = variables_replacer(variables, args[1], args_in_quote[1], false)

		// 一回配列に変換してから比較する
		slice1, _ := divide_split(split(args[0]))
		slice2, _ := divide_split(split(args[1]))

		if len(slice1) != len(slice2) {
			return "false", false
		}

		for i := 0; i < len(slice1); i++ {
			if slice1[i] != slice2[i] {
				return "false", false
			}
		}

		return "true", false
	}

	return "", false
}