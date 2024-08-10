package get

import (
	"fmt"
	"strings"
	"strconv"
	"os"
	"io"
	"net/http"

	system_split "github.com/kurogee/wagyu-script/system-split"
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

func Run(name string, value []string, value_in_quotes []bool, variables *map[string]string, functions *map[string][]string) {
	if name == "local" {
		// ローカルからvalue[2]の名前のファイルを取得し、value[1]の引数を持ち、value[0]の名前で関数として登録する
		value[2] = variables_replacer(variables, value[2], value_in_quotes[2], false)
		file, err := os.ReadFile(value[2])
		if err != nil {
			fmt.Println("The error occurred in local function in get package. [1]")
			return
		}

		(*functions)[value[0]] = []string{value[1], string(file)}
	} else if name == "github" {
		// GitHubからvalue[2]の名前のリポジトリを取得し、value[3]の名前のファイルを取得し、value[1]の引数を持ち、value[0]の名前で関数として登録する
		value[2] = variables_replacer(variables, value[2], value_in_quotes[2], false)
		value[3] = variables_replacer(variables, value[3], value_in_quotes[3], false)
		url := "https://raw.githubusercontent.com/" + value[2] + "/main/" + value[3]
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("The error occurred in github function in get package. [1]")
			return
		}
		defer resp.Body.Close()
		text, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("The error occurred in github function in get package. [2]")
			return
		}

		(*functions)[value[0]] = []string{value[1], string(text)}
	} else if name == "url" {
		// URLからvalue[2]の名前のファイルを取得し、value[1]の引数を持ち、value[0]の名前で関数として登録する
		value[2] = variables_replacer(variables, value[2], value_in_quotes[2], false)
		resp, err := http.Get(value[2])
		if err != nil {
			fmt.Println("The error occurred in url function in get package. [1]")
			return
		}
		defer resp.Body.Close()
		text, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("The error occurred in url function in get package. [2]")
			return
		}

		(*functions)[value[0]] = []string{value[1], string(text)}
	} else {
		fmt.Println("The error occurred in get package. [1]")
	}
}
