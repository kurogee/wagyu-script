package file

import (
	"fmt"
	"strings"
	"strconv"
	"os"
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

func replaceSymbols(input string) string {
	input = strings.ReplaceAll(input, "\n", "\\n")
	input = strings.ReplaceAll(input, "\t", "\\t")
	input = strings.ReplaceAll(input, "\\", "\\\\")
	input = strings.ReplaceAll(input, "\"", "\\\"")
	input = strings.ReplaceAll(input, "'", "\\'")

	return input
}

func Run(name string, value []string, variables *map[string]string) {
	if name == "read" {
		// value[0] ... ファイルのパス value[1] ... 値を入れる変数名
		data, err := os.ReadFile(variables_replacer(*variables, value[0]))
		if err != nil {
			fmt.Println("The error occurred in read function in file package. [1]")
			fmt.Println(err)
			return
		}

		_, ok := (*variables)[value[1]]
		if !ok {
			fmt.Println("The error occurred in read function in file package. [2]")
			fmt.Println("The variable is not found.")
			return
		}

		(*variables)[value[1]] = string(strings.ReplaceAll(string(data), "\n", "\\n"))
	} else if name == "write" {
		// value[0] ... 書き込むファイルのパス value[1] ... 書き込む値（変数もあり）
		write_data := []byte(strings.ReplaceAll(variables_replacer(*variables, value[1]), "\\n", "\n"))
		err := os.WriteFile(variables_replacer(*variables, value[0]), write_data, os.ModePerm)
		if err != nil {
			fmt.Println("The error occurred in write function in file package. [1]")
			fmt.Println(err)
			return
		}
	} else if name == "append" {
		// value[0] ... 追記するファイルのパス value[1] ... 追記する値（変数もあり）
		append_data := strings.ReplaceAll(variables_replacer(*variables, value[1]), "\\n", "\n")
		file_data, err := os.ReadFile(variables_replacer(*variables, value[0]))
		if err != nil {
			fmt.Println("The error occurred in append function in file package. [1]")
			fmt.Println(err)
			return
		}

		all_data := []byte(string(file_data) + append_data)

		err2 := os.WriteFile(variables_replacer(*variables, value[0]), all_data, os.ModePerm)
		if err2 != nil {
			fmt.Println("The error occurred in append function in file package. [2]")
			fmt.Println(err)
			return
		}
	} else if name == "remove" {
		// value[0] ... 削除するファイルのパス
		err := os.Remove(variables_replacer(*variables, value[0]))
		if err != nil {
			fmt.Println("The error occurred in remove function in file package. [1]")
			fmt.Println(err)
			return
		}
	} else if name == "rename" {
		// value[0] ... リネームするファイルのパス value[1] ... 新しいファイル名
		err := os.Rename(variables_replacer(*variables, value[0]), variables_replacer(*variables, value[1]))
		if err != nil {
			fmt.Println("The error occurred in rename function in file package. [1]")
			fmt.Println(err)
			return
		}
	} else if name == "readline" {
		// value[0] ... ファイルのパス value[1] ... 値を入れる変数名 (value[2] ... 読み込む行数 デフォルトは全部)
		data, err := os.ReadFile(variables_replacer(*variables, value[0]))
		if err != nil {
			fmt.Println("The error occurred in readline function in file package. [1]")
			fmt.Println(err)
			return
		}

		lines := strings.Split(string(data), "\n")
		if len(value) == 2 {
			// クオーテーションマークを付ける
			for i, line := range lines[:] {
				lines[i] = replaceSymbols(line)
				lines[i] = "\"" + line + "\""
			}

			fmt.Println(lines)

			(*variables)[value[1]] = strings.Join(lines[:], " ")
		} else {
			how_many, err := strconv.Atoi(variables_replacer(*variables, value[2]))
			if err != nil {
				fmt.Println("The error occurred in readline function in file package. [2]")
				fmt.Println(err)
				return
			}

			if how_many < 0 {
				fmt.Println("The error occurred in readline function in file package. [2]")
				fmt.Println("The value is not a number.")
				return
			}

			// how_manyがlinesの長さより大きい場合は全部読み込む
			if how_many > len(lines) {
				how_many = len(lines)
			}

			// クオーテーションマークを付ける
			for i, line := range lines[:how_many] {
				lines[i] = replaceSymbols(line)
				lines[i] = "\"" + line + "\""
			}

			(*variables)[value[1]] = strings.Join(lines[:how_many], " ")
		}
	}
}
