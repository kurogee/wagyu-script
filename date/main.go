package date

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	system_split "github.com/kurogee/wagyu-script/system_split"
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

/*

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

func Sharp(func_name string, args []string, args_in_quote []bool, variables *map[string]string) (string, bool) {
	if func_name == "nowYear" {
		return strconv.Itoa(time.Now().Year()), false
	} else if func_name == "nowMonth" {
		return strconv.Itoa(int(time.Now().Month())), false
	} else if func_name == "nowDay" {
		return strconv.Itoa(time.Now().Day()), false
	} else if func_name == "nowHour" {
		return strconv.Itoa(time.Now().Hour()), false
	} else if func_name == "nowMinute" {
		return strconv.Itoa(time.Now().Minute()), false
	} else if func_name == "nowSecond" {
		return strconv.Itoa(time.Now().Second()), false
	} else if func_name == "nowDow" || func_name == "nowDayOfWeek" {
		return strconv.Itoa(int(time.Now().Weekday())), false
	} else if func_name == "nowFull" {
		now := time.Now()
		// YYYY-MM-DD HH:MM:SSにフォーマットして返す
		return now.Format("2006-01-02 15:04:05"), true
	} else if func_name == "nowDate" {
		now := time.Now()
		// YYYY-MM-DDにフォーマットして返す
		return now.Format("2006-01-02"), true
	} else if func_name == "nowTime" {
		now := time.Now()
		// HH:MM:SSにフォーマットして返す
		return now.Format("15:04:05"), true
	} else if func_name == "nowUnix" {
		now := time.Now()
		// Unix時間にフォーマットして返す
		return strconv.FormatInt(now.Unix(), 10), false
	}

	return "", false
}

func Run(name string, value []string, value_in_quotes []bool, variables *map[string]string) {
	if name == "now" {
		if value[0] == "format" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in date function in date package. [1]")
				return
			}

			// value[2]には日付のフォーマットが入っている
			format := variables_replacer(variables, value[2], value_in_quotes[2], false)
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
		} else if value[0] == "dow" || value[0] == "dayOfWeek" {
			// value[1]に変数が存在するか確認
			_, ok := (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in dow (dayOfWeek) function in date package. [8]")
				return
			}

			(*variables)[value[1]] = strconv.Itoa(int(time.Now().Weekday()))
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
			base_time := variables_replacer(variables, value[2], value_in_quotes[2], false)

			// value[3]には配列で(年, 月, 日)が入っている
			value_array, value_array_in_quote := divide_split(split(variables_replacer(variables, value[3], value_in_quotes[3], false)))
			for i, v := range(value_array) {
				value_array[i] = variables_replacer(variables, v, value_array_in_quote[i], false)
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
			base_time := variables_replacer(variables, value[2], value_in_quotes[2], false)

			// value[3]には配列で(年, 月, 日)が入っている
			value_array, value_array_in_quote := divide_split(split(variables_replacer(variables, value[3], value_in_quotes[3], false)))
			for i, v := range(value_array) {
				value_array[i] = variables_replacer(variables, v, value_array_in_quote[i], false)
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
		
		value_array := strings.Split(variables_replacer(variables, value[1], value_in_quotes[1], false), " ")[1:]
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
