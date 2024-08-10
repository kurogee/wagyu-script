package http

import (
	"fmt"
	"strings"
	"strconv"
	"net/http"
	"encoding/json"
	"io"
	"os"

	system_split "github.com/kurogee/wagyu-script/system-split"
)

var take_off_quotation = system_split.Take_off_quotation
var split = system_split.Split
var divide_split = system_split.Divide_split

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
	if name == "get" {
		// value[0] = URL value[1] = 変数名
		if len(value) != 2 {
			fmt.Println("The error occurred in the get function in http package. [1]")
			fmt.Println("The number of arguments is invalid.")
			return
		}

		url := variables_replacer(variables, value[0], value_in_quotes[0], false)
		_, ok := (*variables)[value[1]]
		if !ok {
			fmt.Println("The error occurred in the get function in http package. [2]")
			fmt.Println("The variable name is invalid.")
			return
		}

		response, err := http.Get(url)
		if err != nil {
			fmt.Println("The error occurred in the get function in http package. [3]")
			fmt.Println(err)
			return
		}

		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("The error occurred in the get function in http package. [4]")
			fmt.Println(err)
			return
		}

		(*variables)[value[1]] = string(body)
	} else if name == "post" {
		// value[0] = URL value[1] = 変数名 value[2] = データ(配列でキー) value[3] = データ(配列で値)
		if len(value) != 4 {
			fmt.Println("The error occurred in the post function in http package. [1]")
			fmt.Println("The number of arguments is invalid.")
			return
		}

		url := variables_replacer(variables, value[0], value_in_quotes[0], false)
		_, ok := (*variables)[value[1]]
		if !ok {
			fmt.Println("The error occurred in the post function in http package. [2]")
			fmt.Println("The variable name is invalid.")
			return
		}

		d := split(value[2])
		data, data_in_quotes := divide_split(d)

		d_v := split(value[3])
		data_value, data_value_in_quotes := divide_split(d_v)

		if len(data) != len(data_value) {
			fmt.Println("The error occurred in the post function in http package. [3]")
			fmt.Println("The number of data is invalid.")
			return
		}

		// dataに変数が含まれている場合は変数を置き換える
		for i := 0; i < len(data); i++ {
			data[i] = variables_replacer(variables, data[i], data_in_quotes[i], true)
		}

		// data_valueに変数が含まれている場合は変数を置き換える
		for i := 0; i < len(data_value); i++ {
			data_value[i] = variables_replacer(variables, data_value[i], data_value_in_quotes[i], true)
		}

		// dataをURLエンコードする
		data_str := ""
		for i := 0; i < len(data); i++ {
			data_str += data[i] + "=" + data_value[i] + "&"
		}

		data_str = data_str[:len(data_str) - 1]

		response, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data_str))

		if err != nil {
			fmt.Println("The error occurred in the post function in http package. [4]")
			fmt.Println(err)
			return
		}

		defer response.Body.Close()

		// byteを動的に確保
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("The error occurred in the post function in http package. [5]")
			fmt.Println(err)
			return
		}

		(*variables)[value[1]] = string(body)
	} else if name == "toJSON" {
		// value[0] = 変数名 value[1] = JSONファイル名
		if len(value) != 2 {
			fmt.Println("The error occurred in the toJSON function in http package. [1]")
			fmt.Println("The number of arguments is invalid.")
			return
		}

		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in the toJSON function in http package. [2]")
			fmt.Println("The variable name is invalid.")
			return
		}

		json_file := variables_replacer(variables, value[1], value_in_quotes[1], false)
		var json_data interface{}

		err := json.Unmarshal([]byte((*variables)[value[0]]), &json_data)
		if err != nil {
			fmt.Println("The error occurred in the toJSON function in http package. [3]")
			fmt.Println(err)
			return
		}

		file, err := json.MarshalIndent(json_data, "", "  ")
		if err != nil {
			fmt.Println("The error occurred in the toJSON function in http package. [4]")
			fmt.Println(err)
			return
		}

		err = os.WriteFile(json_file, file, 0644)
		if err != nil {
			fmt.Println("The error occurred in the toJSON function in http package. [5]")
			fmt.Println(err)
			return
		}
	}
}

// func Sharp(func_name string, args []string, args_in_quote []bool, variables *map[string]string) (string, bool) {
// 	
// }