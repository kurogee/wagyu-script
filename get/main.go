package get

import (
	"fmt"
	"os"
	"io"
	"net/http"
)

func variables_replacer(variables map[string]string, target string) string {
	val, ok := variables[target]
	if ok {
		return val
	}

	return target
}

func Run(name string, value []string, variables *map[string]string, functions *map[string][]string) {
	if name == "local" {
		// ローカルからvalue[2]の名前のファイルを取得し、value[1]の引数を持ち、value[0]の名前で関数として登録する
		value[2] = variables_replacer(*variables, value[2])
		file, err := os.ReadFile(value[2])
		if err != nil {
			fmt.Println("The error occurred in local function in get package. [1]")
			return
		}

		(*functions)[value[0]] = []string{value[1], string(file)}
	} else if name == "github" {
		// GitHubからvalue[2]の名前のリポジトリを取得し、value[3]の名前のファイルを取得し、value[1]の引数を持ち、value[0]の名前で関数として登録する
		value[2] = variables_replacer(*variables, value[2])
		value[3] = variables_replacer(*variables, value[3])
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
		value[2] = variables_replacer(*variables, value[2])
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
