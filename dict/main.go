package array

import (
	"strings"
	"strconv"
	"regexp"

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

func ParseDict(dict string, variables *map[string]string) map[string]string {
	// dictは、"key1" : "value1" "key2" : "value2" "key3" : "value3"のような形式
	// まずは、dictを:がないスペースで分割する
	// ""はあってもなくてもいい
	
	// まずはクオート内のスペースを置き換える
	re := regexp.MustCompile(`"[^"]+"`)
	matches := re.FindAllString(dict, -1)
	for _, match := range matches {
		dict = strings.ReplaceAll(dict, match, strings.ReplaceAll(match, " ", "__SPACE__"))
	}
	
	dict = dict + " "

	re = regexp.MustCompile(`[^\s]+\s*:\s*[^\s]+\s`)
	matches2 := re.FindAllStringSubmatch(dict, -1)

	// keyとvalueを取り出す
	dict_map := make(map[string]string)
	for _, match := range matches2 {
		// matchを:で分割する
		re = regexp.MustCompile(`([^\s]+)\s*:\s*([^\s]+)\s`)
		match = re.FindStringSubmatch(match[0])

		key := match[1]
		value := match[2]

		// keyとvalueをsplitする
		splited_key, key_in_value := divide_split(split(key))
		splited_value, value_in_quote := divide_split(split(value))

		// keyとvalueを変数で置き換える
		key = variables_replacer(variables, splited_key[0], key_in_value[0], false)
		value = variables_replacer(variables, splited_value[0], value_in_quote[0], false)

		dict_map[key] = value
	}

	// クオート内のスペースを元に戻す
	for key, value := range dict_map {
		if strings.Contains(value, "__SPACE__") {
			dict_map[key] = strings.ReplaceAll(value, "__SPACE__", " ")
		}

		if strings.Contains(key, "__SPACE__") {
			dict_map[strings.ReplaceAll(key, "__SPACE__", " ")] = value
			delete(dict_map, key)
		}
	}

	return dict_map
}

func Run(name string, value []string, value_in_quotes []bool, variables *map[string]string) {
	if name == "new" {
		// value[0] = 変数名 value[1] = keyが入った配列 value[2] = valueが入った配列
		splited_value1, value1_in_quote := divide_split(split(value[1]))
		splited_value2, value2_in_quote := divide_split(split(value[2]))

		// keyとvalueを変数で置き換える
		for i := 0; i < len(splited_value1); i++ {
			splited_value1[i] = variables_replacer(variables, splited_value1[i], value1_in_quote[i], true)
			splited_value2[i] = variables_replacer(variables, splited_value2[i], value2_in_quote[i], true)
		}

		// forでkeyとvalueを取り出す
		dict := make(map[string]string)
		for i := 0; i < len(splited_value1); i++ {
			dict[splited_value1[i]] = splited_value2[i]
		}

		// dictをkey1 : value1 key2 : value2の形式に変換する
		dict_str := ""
		for key, value := range dict {
			dict_str += key + " : " + value + " "
		}

		// 最後の空白を取り除く
		dict_str = dict_str[:len(dict_str) - 1]

		(*variables)[value[0]] = dict_str
	} else if name == "set" {
		// value[0] = 変数名 value[1] = key value[2] = value
		// value[0]が変数名かどうかを判定する
		dict := variables_replacer(variables, value[0], value_in_quotes[0], false)
		
		// dictをパースする
		dict_map := ParseDict(dict, variables)

		// keyとvalueを変数で置き換える
		splited_value1, _ := divide_split(split(value[1]))
		splited_value2, value2_in_quote := divide_split(split(value[2]))

		splited_value1[0] = variables_replacer(variables, splited_value1[0], true, false)
		splited_value2[0] = variables_replacer(variables, splited_value2[0], value2_in_quote[0], true)

		// keyとvalueを追加する
		dict_map[splited_value1[0]] = splited_value2[0]

		// dictをkey1 : value1 key2 : value2の形式に変換する
		dict_str := ""
		for key, value := range dict_map {
			dict_str += key + " : " + value + " "
		}

		// 最後の空白を取り除く
		dict_str = dict_str[:len(dict_str) - 1]

		(*variables)[value[0]] = dict_str
	} else if name == "remove" {
		// value[0] = 変数名 value[1] = key
		// value[0]が変数名かどうかを判定する
		dict := variables_replacer(variables, value[0], value_in_quotes[0], false)

		// dictをパースする
		dict_map := ParseDict(dict, variables)

		// keyを変数で置き換える
		splited_value1, _ := divide_split(split(value[1]))
		splited_value1[0] = variables_replacer(variables, splited_value1[0], true, false)

		// keyを削除する
		delete(dict_map, splited_value1[0])

		// dictをkey1 : value1 key2 : value2の形式に変換する
		dict_str := ""
		for key, value := range dict_map {
			dict_str += key + " : " + value + " "
		}

		// 最後の空白を取り除く
		dict_str = dict_str[:len(dict_str) - 1]

		(*variables)[value[0]] = dict_str
	}
}

func Sharp(func_name string, args []string, args_in_quote []bool, variables *map[string]string) (string, bool) {
	if func_name == "get" {
		// args[0] = 辞書or変数名 args[1] = key -> valueを返す なければ-1を返す
		// args[0]が変数名かどうかを判定する
		dict := variables_replacer(variables, args[0], args_in_quote[0], false)

		// dictをパースする
		dict_map := ParseDict(dict, variables)

		// keyを入れる
		key := variables_replacer(variables, args[1], args_in_quote[1], false)

		// keyが存在するか確認する
		value, ok := dict_map[key]
		if !ok {
			return "-1", false
		}

		// valueが数値に変換できそうだったらfalseを返す
		if _, err := strconv.Atoi(value); err == nil {
			return value, false
		} else {
			return value, true
		}
	} else if func_name == "check" {
		// args[0] = 辞書or変数名 args[1] = key
		// args[0]が変数名かどうかを判定する
		dict := variables_replacer(variables, args[0], args_in_quote[0], false)

		// dictをパースする
		dict_map := ParseDict(dict, variables)

		// keyを入れる
		key := variables_replacer(variables, args[1], args_in_quote[1], false)

		// keyが存在するか確認する
		_, ok := dict_map[key]
		return strconv.FormatBool(ok), false
	} else if func_name == "search" {
		// args[0] = 辞書or変数名 args[1] = value -> keyを返す なければ-1を返す
		// args[0]が変数名かどうかを判定する
		dict := variables_replacer(variables, args[0], args_in_quote[0], false)

		// dictをパースする
		dict_map := ParseDict(dict, variables)

		// valueを入れる
		value := variables_replacer(variables, args[1], args_in_quote[1], false)

		// キーを探す
		for key, val := range dict_map {
			if val == value {
				return key, false
			}
		}

		return "-1", false
	}

	return "", false
}