package array

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

func replaceSymbols(input string) string {
	input = strings.ReplaceAll(input, "\n", "\\n")
	input = strings.ReplaceAll(input, "\t", "\\t")
	input = strings.ReplaceAll(input, "\\", "\\\\")
	input = strings.ReplaceAll(input, "\"", "\\\"")
	input = strings.ReplaceAll(input, "'", "\\'")

	return input
}

func add_quotation(val string) string {
	if strings.Contains(val, " ") {
		val = replaceSymbols(val)
		val = "\"" + val + "\""
	}

	return val
}

func Run(name string, value []string, variables *map[string]string) {
	// 基本的にvalue[0]は変数名
	if name == "reset" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in reset function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		(*variables)[value[0]] = ""
	} else if name == "split" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in split function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		(*variables)[value[0]] = strings.Join(strings.Split(variables_replacer(*variables, value[1]), value[2]), " ")
	} else if name == "addbeg" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addbeg function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(*variables, value[1])
		val = add_quotation(val)

		(*variables)[value[0]] = val + " " + (*variables)[value[0]]
	} else if name == "addend" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addend function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(*variables, value[1])
		val = add_quotation(val)

		(*variables)[value[0]] = (*variables)[value[0]] + " " + val
	} else if name == "addnth" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in add function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		index, err := strconv.Atoi(variables_replacer(*variables, value[1]))
		if err != nil {
			fmt.Println("The error occurred in add function in array package. [2]")
			fmt.Println("The index is not integer.")
		}

		val := variables_replacer(*variables, value[2])
		val = add_quotation(val)

		(*variables)[value[0]] = (*variables)[value[0]][:index + 1] + " " + val + (*variables)[value[0]][index + 1:]
	} else if name == "replace" {
		// value[0] = target, value[1] = old index, value[2] = new value
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in replace function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		index, err := strconv.Atoi(variables_replacer(*variables, value[1]))
		if err != nil {
			fmt.Println("The error occurred in replace function in array package. [2]")
			fmt.Println("The index is not integer.")
		}

		// 一回配列に変換してから置換する
		slice := strings.Split((*variables)[value[0]], " ")

		val := variables_replacer(*variables, value[2])
		val = add_quotation(val)

		slice[index] = val
		(*variables)[value[0]] = strings.Join(slice, " ")
	} else if name == "delnth" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in delnth function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		index, err := strconv.Atoi(variables_replacer(*variables, value[1]))
		if err != nil {
			fmt.Println("The error occurred in delnth function in array package. [2]")
			fmt.Println("The index is not integer.")
		}

		// 一回配列に変換してから削除する
		slice := strings.Split((*variables)[value[0]], " ")
		slice = append(slice[:index], slice[index + 1:]...)
		(*variables)[value[0]] = strings.Join(slice, " ")
	} else if name == "sort" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in sort function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		// 一回配列に変換してからソートする
		slice := strings.Split((*variables)[value[0]], " ")
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
		slice := strings.Split((*variables)[value[0]], " ")
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
		slice := strings.Split((*variables)[value[1]], " ")

		// 配列から値を検索し、そのインデックスをvalue[0]の変数に格納 なければ-1
		index := -1
		for i, val := range slice {
			if val == variables_replacer(*variables, value[2]) {
				index = i
				break
			}
		}

		(*variables)[value[0]] = strconv.Itoa(index)
	}
}
