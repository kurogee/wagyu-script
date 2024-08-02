package date

import (
	"fmt"
	"strconv"
	"strings"

	"time"
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

func _set_format(format string) string {
	format = strings.ReplaceAll(format, "YYYY", "2006")
	format = strings.ReplaceAll(format, "MM", "01")
	format = strings.ReplaceAll(format, "DD", "02")
	format = strings.ReplaceAll(format, "HH", "15")
	format = strings.ReplaceAll(format, "mm", "04")
	format = strings.ReplaceAll(format, "SS", "05")

	return format
}

func Run(name string, value []string, variables *map[string]string) {
	if name == "now" {
		if value[0] == "format" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in date function in date package. [1]")
				return
			}

			// value[2]には日付のフォーマットが入っている
			format := variables_replacer(*variables, value[2])
			format = strings.ReplaceAll(format, ":auto:", "YYYY-MM-DD HH:mm:SS")
			format = _set_format(format)

			// フォーマットの形式: YYYY-MM-DD HH:MM:SS
			(*variables)[value[1]] = time.Now().Format(format)
		} else if value[0] == "year" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in year function in date package. [2]")
				return
			}

			(*variables)[value[1]] = strconv.Itoa(time.Now().Year())
		} else if value[0] == "month" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in month function in date package. [3]")
				return
			}

			(*variables)[value[1]] = strconv.Itoa(int(time.Now().Month()))
		} else if value[0] == "day" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in day function in date package. [4]")
				return
			}

			(*variables)[value[1]] = strconv.Itoa(time.Now().Day())
		} else if value[0] == "hour" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in hour function in date package. [5]")
				return
			}

			(*variables)[value[1]] = strconv.Itoa(time.Now().Hour())
		} else if value[0] == "minute" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in minute function in date package. [6]")
				return
			}

			(*variables)[value[1]] = strconv.Itoa(time.Now().Minute())
		} else if value[0] == "second" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in second function in date package. [7]")
				return
			}

			(*variables)[value[1]] = strconv.Itoa(time.Now().Second())
		} else {
			fmt.Println("The error occurred in now in date package. [8]")
			fmt.Println("The function name is invalid.")
			return
		}
	} else if name == "calc" {
		if value[0] == "add" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in add function in date package. [1]")
				return
			}

			// value[2]には日付がフォーマットに則って入っている
			base_time := variables_replacer(*variables, value[2])

			// value[3]には配列で(年, 月, 日)が入っている
			value_array := strings.Split(variables_replacer(*variables, value[3]), " ")
			for i, v := range(value_array) {
				_, err := strconv.Atoi(v)
				if v == "" || err != nil {
					value_array[i] = "0"
				}
			}

			// フォーマットの形式: YYYY-MM-DD HH:MM:SS
			t, err := time.Parse(_set_format("YYYY-MM-DD"), base_time)
			if err != nil {
				fmt.Println("The error occurred in add function in date package. [2]")
				return
			}

			// 日付の加算
			add_year, err := strconv.Atoi(value_array[0])
			if err != nil {
				fmt.Println("The error occurred in add function in date package. [3]")
				return
			}

			add_month, err := strconv.Atoi(value_array[1])
			if err != nil {
				fmt.Println("The error occurred in add function in date package. [4]")
				return
			}

			add_day, err := strconv.Atoi(value_array[2])
			if err != nil {
				fmt.Println("The error occurred in add function in date package. [5]")
				return
			}

			t = t.AddDate(add_year, add_month, add_day)

			(*variables)[value[1]] = t.Format(_set_format("YYYY-MM-DD"))
		} else if value[0] == "sub" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in add function in date package. [1]")
				return
			}

			// value[2]には日付がフォーマットに則って入っている
			base_time := variables_replacer(*variables, value[2])

			// value[3]には配列で(年, 月, 日)が入っている
			value_array := strings.Split(variables_replacer(*variables, value[3]), " ")[1:]
			for i, v := range(value_array) {
				if v == "" {
					value_array[i] = "0"
				}
			}

			// フォーマットの形式: YYYY-MM-DD HH:MM:SS
			t, err := time.Parse(_set_format("YYYY-MM-DD"), base_time)
			if err != nil {
				fmt.Println("The error occurred in add function in date package. [2]")
				fmt.Println(err)
				return
			}

			// 日付の加算
			add_year, err := strconv.Atoi(value_array[0])
			if err != nil {
				fmt.Println("The error occurred in add function in date package. [3]")
				return
			}

			add_month, err := strconv.Atoi(value_array[1])
			if err != nil {
				fmt.Println("The error occurred in add function in date package. [4]")
				return
			}

			add_day, err := strconv.Atoi(value_array[2])
			if err != nil {
				fmt.Println("The error occurred in add function in date package. [5]")
				return
			}

			t = t.AddDate(-add_year, -add_month, -add_day)

			(*variables)[value[1]] = t.Format(_set_format("YYYY-MM-DD"))
		} else {
			fmt.Println("The error occurred in calc in date package. [1]")
			fmt.Println("The function name is invalid.")
			return
		}
	} else if name == "format" {
		// value[0]...代入する変数名
		// value[1]...日付時間の配列 長さが3の時は(年, 月, 日) 長さが6の時は(年, 月, 日, 時, 分, 秒)
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in format function in date package. [1]")
			return
		}
		
		value_array := strings.Split(variables_replacer(*variables, value[1]), " ")
		// value_arrayを2桁に変換 1桁の場合は前に0を付ける
		for i, v := range(value_array) {
			if len(v) == 1 {
				value_array[i] = "0" + v
			}
		}

		if len(value_array) == 3 {
			// フォーマットの形式: YYYY-MM-DD
			t, err := time.Parse(_set_format("YYYY-MM-DD"), value_array[0] + "-" + value_array[1] + "-" + value_array[2])
			if err != nil {
				fmt.Println("The error occurred in format function in date package. [2]")
				return
			}

			(*variables)[value[0]] = t.Format(_set_format("YYYY-MM-DD"))
		} else if len(value_array) == 6 {
			// フォーマットの形式: YYYY-MM-DD HH:MM:SS
			t, err := time.Parse(_set_format("YYYY-MM-DD HH:mm:SS"), value_array[0] + "-" + value_array[1] + "-" + value_array[2] + " " + value_array[3] + ":" + value_array[4] + ":" + value_array[5])
			if err != nil {
				fmt.Println("The error occurred in format function in date package. [3]")
				return
			}

			(*variables)[value[0]] = t.Format(_set_format("YYYY-MM-DD HH:mm:SS"))
		} else {
			fmt.Println("The error occurred in format function in date package. [4]")
			fmt.Println("The length of the date array is invalid.")
			return
		}
	}
}
