package main

import (
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"math/rand/v2"
	"os"

	system_split "github.com/kurogee/wagyu-script/system-split"

	array_pack "github.com/kurogee/wagyu-script/array"
	date_pack "github.com/kurogee/wagyu-script/date"
	file_pack "github.com/kurogee/wagyu-script/file"
	string_pack "github.com/kurogee/wagyu-script/string"
	get_pack "github.com/kurogee/wagyu-script/get"
	regex_pack "github.com/kurogee/wagyu-script/regex"
	dict_pack "github.com/kurogee/wagyu-script/dict"
	// http_pack "github.com/kurogee/wagyu-script/http"

	math_sharp_functions "github.com/kurogee/wagyu-script/maths"

	"github.com/Knetic/govaluate"
)

var split = system_split.Split
var divide_split = system_split.Divide_split
var take_off_quotation = system_split.Take_off_quotation

type Parse struct {
	parsed []string
	parsed_in_quotes []bool
	parsed_meaning CodeType
}

type CodeType struct {
	name string
	conjunction string
	value []string
	value_in_quotes []bool
}

type Package map[string]func(string, []string, []bool, *map[string]string)
type Package_with_functions map[string]func(string, []string, []bool, *map[string]string, *map[string][]string)
type Package_sharp_functions map[string]func(string, []string, []bool, *map[string]string) (string, bool)

// packageの一覧を定義
var packages = Package{
	"date": date_pack.Run,
	"file": file_pack.Run,
	"array": array_pack.Run,
	"string": string_pack.Run,
	"regex": regex_pack.Run,
	"dict": dict_pack.Run,
	// "http": http_pack.Run,
}

var packages_with_functions = Package_with_functions{
	"get": get_pack.Run,
}

var packages_sharp_functions = Package_sharp_functions{
	"math": math_sharp_functions.Sharp,
	"date": date_pack.Sharp,
	"regex": regex_pack.Sharp,
	"string": string_pack.Sharp,
	"array": array_pack.Sharp,
	"file": file_pack.Sharp,
	"dict": dict_pack.Sharp,
}

func contains(s []string, e string) bool {
	for _, a := range(s) {
		if a == e {
			return true
		}
	}

	return false
}

func giveSymbols(input string) string {
	input = strings.ReplaceAll(input, "\\n", "\n")
	input = strings.ReplaceAll(input, "\\t", "\t")
	input = strings.ReplaceAll(input, "\\\\", "\\")
	input = strings.ReplaceAll(input, "\\\"", "\"")
	input = strings.ReplaceAll(input, "\\'", "'")
	
	return input
}


func parser(code string) (p Parse) {
	// codeから不要な改行や文の始めのインデントを削除
	code = strings.ReplaceAll(code, "\n", "")
	code = strings.ReplaceAll(code, "\r", "")
	code = strings.ReplaceAll(code, "\t", "")

	// 正規表現で文初めの2つもしくは4つのスペースを削除
	re := regexp.MustCompile(`^ {2}|^ {4}`)
	code = re.ReplaceAllString(code, "")

	var mem string = strings.ReplaceAll(code, " ", "")

	// #関数名(引数)で、#関数名 (引数) と書かれていたら、間のスペースを削除
	re = regexp.MustCompile(`(#[^\(]*) (\()`)
	code = re.ReplaceAllString(code, "$1$2")

	// (sharp|fnc) 関数名(引数) と書かれていたら、間にスペースを入れる
	re = regexp.MustCompile(`(sharp|fnc) ([^\(]+)\(`)
	code = re.ReplaceAllString(code, "$1 $2 (")
	
	if mem == "" {
		p.parsed = []string{"None", ""}
		p.parsed_meaning = CodeType{"None", "", []string{}, []bool{}}
		return
	}

	// 先頭に//があったらコメントアウトを削除
	if strings.HasPrefix(code, "//") {
		p.parsed = []string{"None", ""}
		p.parsed_meaning = CodeType{"None", "", []string{}, []bool{}}
		return
	}

	divided_code := split(code)

	// divided_codeをstringとboolに分ける
	var divided_code_str []string
	var divided_code_bool []bool
	for _, val := range(divided_code) {
		for key := range(val) {
			divided_code_str = append(divided_code_str, key)
		}
		for _, val2 := range(val) {
			divided_code_bool = append(divided_code_bool, val2)
		}
	}

	p.parsed = divided_code_str
	p.parsed_in_quotes = divided_code_bool

	// もしdivided_code[1]が記号ではなかったら、divided_code[1:]をvalueに格納しconjunctionには<を格納
	re = regexp.MustCompile(`^[><=]$`)
	if re.MatchString(divided_code_str[1]) {
		p.parsed_meaning.name = strings.ReplaceAll(divided_code_str[0], " ", "")
		p.parsed_meaning.conjunction = divided_code_str[1]
		p.parsed_meaning.value = divided_code_str[2:]
		p.parsed_meaning.value_in_quotes = divided_code_bool[2:]
	} else {
		p.parsed_meaning.name = strings.ReplaceAll(divided_code_str[0], " ", "")
		p.parsed_meaning.conjunction = "<"
		p.parsed_meaning.value = divided_code_str[1:]
		p.parsed_meaning.value_in_quotes = divided_code_bool[1:]
	}

	return
}

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

	// if add_quotes {
	// 	if _, err := strconv.Atoi(target); err != nil {
	// 		if _, err := strconv.ParseFloat(target, 64); err != nil {
	// 			if len(strings.Split(target, " ")) == 1 {
	// 				return "\"" + target + "\""
	// 			}
	// 		}
	// 	}
	// 
	// 	return target
	// }

	return target
}

func variables_replacers(variables map[string]string, sentence string, targets []string, targets_in_quote []bool) string {
	var result string = sentence
	var count int = 1
	for i, target := range(targets) {
		if targets_in_quote[i] {
			// 両端のクオートを取り除く
			mem := take_off_quotation(target)
			result = strings.ReplaceAll(result, ":" + strconv.Itoa(count) + ":", mem)
		} else {
			val, ok := variables[target]
			if ok {
				val = take_off_quotation(val)
				result = strings.ReplaceAll(result, ":" + strconv.Itoa(count) + ":", val)
			} else {
				mem := take_off_quotation(target)
				result = strings.ReplaceAll(result, ":" + strconv.Itoa(count) + ":", mem)
			}
		}

		count++
	}

	return result
}

func calc_expression(expression string, parameters map[string]interface{}) string {
	govaluate_result, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		fmt.Println("The error occurred in calc. [1] / 演算の準備をしている際にエラーが発生しました。")
		fmt.Println(err)
		return ""
	}

	result, err2 := govaluate_result.Evaluate(parameters)
	if err2 != nil {
		fmt.Println("The error occurred in calc. [2] / 演算をする際にエラーが発生しました。")
		fmt.Println(err2)
		return ""
	}

	return fmt.Sprintf("%v", result)
}

func sharp_functions(func_name string, args []string, args_in_quote []bool, variables *map[string]string, functions *map[string][]string, sharps *map[string][]string) (string, bool) {
	if variables == nil {
		variables = &map[string]string{}
	}

	args_str := args
	args_in_quote2 := args_in_quote

	// argsの中にさらにシャープ関数がある場合は、それを実行
	for i, arg2 := range(args) {
		if strings.HasPrefix(arg2, "#") {
			func_name2 := strings.Split(arg2, "(")[0][1:]
			re := regexp.MustCompile(`^#[^\(]*\(|\)$`)

			args2 := split(re.ReplaceAllString(arg2, ""))
			args2_str, args2_in_quote := divide_split(args2)
			
			args_str[i], args_in_quote2[i] = sharp_functions(func_name2, args2_str, args2_in_quote, variables, functions, sharps)
		}
	}

	// もしfunc_nameに.が含まれていたら、そのパッケージを実行
	if strings.Contains(func_name, ".") {
		package_name := strings.Split(func_name, ".")[0]
		val, ok := packages_sharp_functions[package_name]
		if !ok {
			fmt.Println("The error occurred in sharp_functions. [1] / パッケージが見つからないためエラーが発生しました。")
			fmt.Println("The package is not found.")
			return "", false
		}

		result, result_bool := val(strings.Split(func_name, ".")[1], args, args_in_quote, variables)
		return result, result_bool
	}

	if func_name == "arrayAt" {
		// args[0] は配列名、args[1] はインデックス
		// インデックスは数値である必要がある
		index, err := strconv.Atoi(variables_replacer(variables, args_str[1], args_in_quote2[1], false))
		if err != nil {
			fmt.Println("The error occurred in arrayAt. [1] / インデックスが数値でないためエラーが発生しました。")
			fmt.Println(err)
		}

		// 配列名が存在する場合は、その配列を取得
		array, _ := divide_split(split(variables_replacer(variables, args_str[0], args_in_quote2[0], false)))

		// インデックスが配列の範囲内にあるか確認
		if index < 0 || index >= len(array) {
			fmt.Println("The error occurred in arrayAt. [3] / インデックスが範囲外のためエラーが発生しました。")
			fmt.Println("The index is out of range.")
			return "", false
		}

		return array[index], args_in_quote2[0]
	} else if func_name == "at" {
		// args[0] は文字列or変数名、args[1] はインデックス
		// インデックスは数値である必要がある
		index, err := strconv.Atoi(variables_replacer(variables, args_str[1], args_in_quote2[1], false))
		if err != nil {
			fmt.Println("The error occurred in at. [1] / インデックスが数値でないためエラーが発生しました。")
			fmt.Println(err)
		}

		// 文字列が変数名である場合は、その変数名の値を取得
		mem := strings.Split(variables_replacer(variables, args_str[0], args_in_quote2[0], false), "")[index]

		// 両端にクオートをつける
		mem = "\"" + mem + "\""

		return mem, true
	} else if func_name == "arrayLen" {
		// args[0] は配列名
		// 配列名が存在するか確認
		_, ok := (*variables)[args[0]]
		if !ok {
			fmt.Println("The error occurred in arrayLen. [1] / 配列名が見つからないためエラーが発生しました。")
			fmt.Println("The array name is not found.")
			return "", false
		}

		// 配列名が存在する場合は、その配列を取得
		array, _ := divide_split(split((*variables)[args[0]]))

		return strconv.Itoa(len(array)), args_in_quote2[0]
	} else if func_name == "len" {
		// args[0] は文字列
		// 文字列の長さを返す
		return strconv.Itoa(len(variables_replacer(variables, args[0], args_in_quote2[0], false))), args_in_quote2[0]
	} else if func_name == "calc" {
		// args[0] は計算式
		// args[1:] は変数名

		if len(args) == 1 {
			return calc_expression(variables_replacer(variables, args[0], args_in_quote2[0], true), map[string]interface{}{}), false
		} else if len(args) > 1 {
			return calc_expression(variables_replacers((*variables), args[0], args[1:], args_in_quote2[1:]), map[string]interface{}{}), false
		}
	} else if func_name == "from" {
		// 引数の1番目の数字から2番目の数字までの連番の配列を作成
		// args[0] は開始番号、args[1] は終了番号

		// もし引数が変数名である場合は、その変数名の値を取得
		start_arg := variables_replacer(variables, args_str[0], args_in_quote2[0], false)
		end_arg := variables_replacer(variables, args_str[1], args_in_quote2[1], false)

		// fmt.Println(start_arg, end_arg)

		start, err := strconv.Atoi(start_arg)
		if err != nil {
			fmt.Println("The error occurred in from. [1] / 開始番号が(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
			return "", false
		}

		end, err := strconv.Atoi(end_arg)
		if err != nil {
			fmt.Println("The error occurred in from. [2] / 終了番号が(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
			return "", false
		}

		var result string
		var builder strings.Builder

		for i := start; i <= end; i++ {
			builder.WriteString(strconv.Itoa(i) + " ")
		}

		result = builder.String()

		// 最後のスペースを削除
		result = result[:len(result) - 1]

		return result, false
	} else if func_name == "format" {
		// args[0]は接続する文字列フォーマット、args[1:]は変数名
		// printfと同じように変数を埋め込む
		return variables_replacers((*variables), args[0], args[1:], args_in_quote2[1:]), true
	} else if func_name == "var" {
		// 変数を宣言し、args[0]にargs[1]を代入
		// 変数名を返す
		(*variables)[args[0]] = variables_replacer(variables, args_str[1], args_in_quote2[1], false)
		
		return args[0], args_in_quote2[0]
	} else if func_name == "remove" {
		// args[0] は配列名、args[1] は削除する値のインデックス
		_, ok := (*variables)[args[0]]
		if !ok {
			fmt.Println("The error occurred in remove. [1] / 配列名が見つからないためエラーが発生しました。")
			fmt.Println("The array name is not found.")
		}

		// 配列名が存在する場合は、その配列を取得
		array, _ := divide_split(split((*variables)[args[0]]))

		// インデックスが数値であるか確認
		index, err := strconv.Atoi(variables_replacer(variables, args_str[1], args_in_quote2[1], false))
		if err != nil {
			fmt.Println("The error occurred in remove. [2] / インデックスが(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
		}

		// インデックスが配列の範囲内にあるか確認
		if index < 0 || index >= len(array) {
			fmt.Println("The error occurred in remove. [3] / インデックスが範囲外のためエラーが発生しました。")
			fmt.Println("The index is out of range.")
		}

		// インデックスの要素を削除
		array = append(array[:index], array[index + 1:]...)

		// 配列を再度文字列に変換して格納
		return strings.Join(array, " "), false
	} else if func_name == "rand" {
		// args[0] は開始値、args[1] は終了値
		start, err := strconv.Atoi(variables_replacer(variables, args_str[0], args_in_quote2[0], false))
		if err != nil {
			fmt.Println("The error occurred in rand. [1] / 開始値が(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
		}

		end, err := strconv.Atoi(variables_replacer(variables, args_str[1], args_in_quote2[1], false))
		if err != nil {
			fmt.Println("The error occurred in rand. [2] / 終了値が(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
		}

		// start〜endの範囲（endを含む）の乱数を返す
		return strconv.Itoa(rand.IntN(end - start + 1) + start), false
	} else if func_name == "convert" {
		// args[0] は変換する値、args[1] は変換後の型
		switch args[1] {
		case "int":
			// 小数点があったら、それを取り除く
			val, err := strconv.Atoi(strings.Split(variables_replacer(variables, args_str[0], args_in_quote2[0], false), ".")[0])
			if err != nil {
				fmt.Println("The error occurred in convert. [1] / 変換する値が(変数の中の値が)数値でないためエラーが発生しました。")
				fmt.Println(err)
			}

			return strconv.Itoa(val), args_in_quote2[0]
		case "float":
			val, err := strconv.ParseFloat(variables_replacer(variables, args_str[0], args_in_quote2[0], false), 64)
			if err != nil {
				fmt.Println("The error occurred in convert. [2] / 変換する値が(変数の中の値が)数値でないためエラーが発生しました。")
				fmt.Println(err)
			}

			return strconv.FormatFloat(val, 'f', 10, 64), args_in_quote2[0]
		}
	} else if func_name == "repeat" {
		// args[0] は繰り返す文字列、args[1] は繰り返す回数
		repeat_num, err := strconv.Atoi(variables_replacer(variables, args_str[1], args_in_quote2[1], false))
		if err != nil {
			fmt.Println("The error occurred in repeat. [1] / 繰り返す回数が(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
		}

		return strings.Repeat(variables_replacer(variables, args_str[0], args_in_quote2[0], false), repeat_num), args_in_quote2[0]
	} else if func_name == "all" {
		// argsの長さが1より大きければ、args[0:]が値or変数名
		if len(args) > 1 {
			for i, arg := range(args) {
				mem := variables_replacer(variables, arg, args_in_quote2[i], false)
				if mem != "1" && mem != "true" {
					return "false", false
				}
			}

			return "true", false
		}

		// args[0] は変数名の配列
		// すべての変数が1もしくはtrueであるか確認
		mem := split(variables_replacer(variables, args_str[0], args_in_quote2[0], false))
		vals, vals_in_quote := divide_split(mem)

		for i, val := range(vals) {
			if variables_replacer(variables, val, vals_in_quote[i], false) != "1" && variables_replacer(variables, val, vals_in_quote[i], false) != "true" {
				return "false", false
			}
		}

		return "true", false
	} else if func_name == "any" {
		// argsの長さが1より大きければ、args[0:]が値or変数名
		if len(args) > 1 {
			for i, arg := range(args) {
				mem := variables_replacer(variables, arg, args_in_quote2[i], false)
				if mem == "1" || mem == "true" {
					return "true", false
				}
			}

			return "false", false
		}

		// args[0] は変数名の配列
		// どれかの変数が1もしくはtrueであるか確認
		mem := split(variables_replacer(variables, args_str[0], args_in_quote2[0], false))
		vals, vals_in_quote := divide_split(mem)

		for i, val := range(vals) {
			if variables_replacer(variables, val, vals_in_quote[i], false) == "1" || variables_replacer(variables, val, vals_in_quote[i], false) == "true" {
				return "true", false
			}
		}

		return "false", false
	} else if func_name == "" {
		for i, arg := range(args) {
			args[i] = variables_replacer(variables, arg, args_in_quote2[i], true)
		}

		var formula string = strings.Join(args, " ")

		return calc_expression(formula, map[string]interface{}{}), false
	// もしもfunc_nameがsharpsに含まれていたら、その関数を実行
	} else if _, ok := (*sharps)[func_name]; ok {
		// args[0:] は引数
		// sharps[func_name][0] は引数の名前（スペース区切り）、sharps[func_name][1] は関数の中身

		// 引数の数が一致しているか確認
		if len(args) != len(split((*sharps)[func_name][0])) {
			fmt.Println("The error occurred in sharp_functions. [2] / 引数の数が一致していないためエラーが発生しました。")
			fmt.Println("The number of arguments is not match.")
			return "", false
		}

		// 引数を変数に格納
		arg_str, arg_in_quote := divide_split(split((*sharps)[func_name][0]))
		for i, arg := range(args) {
			if arg_in_quote[i] {
				(*variables)[arg_str[i]] = arg
			} else {
				(*variables)[arg_str[i]] = variables_replacer(variables, arg, arg_in_quote[i], false)
			}
		}

		// 関数の中身を実行
		codes := splitOutsideSemicolons((*sharps)[func_name][1])

		for _, code := range(codes) {
			p := parser(code)
			status := p.runner(variables, functions, sharps, &func_name, false)
			if status == 1 {
				break
			}
		}

		// 関数の中身を実行した後、関数の中身で定義された変数を削除
		for _, arg := range(arg_str) {
			delete((*variables), arg)
		}

		// 返り値を返す
		// 文字列か文字列じゃないかのboolも返す
		return (*variables)["0__return__"], false
	}

	return "", false
}

func (p Parse) runner(variables *map[string]string, functions *map[string][]string, sharps *map[string][]string, before_func_name *string, top bool) int {
	/*
		returnの数字の意味
		0: 正常終了
		1: return文があった、関数を抜ける
		2: while文をbreakする
		3: while文をcontinueする
		-1: エラー、異常終了
	*/
	name := p.parsed_meaning.name
	conjunction := p.parsed_meaning.conjunction
	value := p.parsed_meaning.value
	value_in_quotes := p.parsed_meaning.value_in_quotes

	if name == "None" {
		return 0
	}

	if name == "$" {
		code := *before_func_name + " " + conjunction + " " + strings.Join(value, " ")
		parsedCode := parser(code)
		status := parsedCode.runner(variables, functions, sharps, before_func_name, false)
		if status == 1 {
			return 1
		} else if status == 2 {
			return 2
		} else if status == 3 {
			return 3
		}
		return 0
	}

	for i, val := range(value) {
		if name != "while" && name != "if" && name != "match" {
			if strings.HasPrefix(val, "#") {
				func_name := strings.Split(val, "(")[0][1:]
				re := regexp.MustCompile(`^#[^\(]*\(|\)$`)

				args := split(re.ReplaceAllString(val, ""))
				args_str, args_in_quote := divide_split(args)

				value[i], value_in_quotes[i] = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)
			}
		}
	}
	
	mem := &value
	if name != "if" {
		for i, val := range(*mem) {
			(*mem)[i] = take_off_quotation(val)
			if strings.HasPrefix(val, "(") && strings.HasSuffix(val, ")") {
				(*mem)[i] = strings.Trim(val, "()")
			} else if strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}") {
				(*mem)[i] = strings.Trim(val, "{}")
			}
		}
	}

	// valueにカッコがついている場合は取り除く
	for i, val := range(*mem) {
		if strings.HasPrefix(val, "(") && strings.HasSuffix(val, ")") {
			(*mem)[i] = strings.Trim(val, "()")
		} else if strings.HasPrefix(val, "{") && strings.HasSuffix(val, "}") {
			(*mem)[i] = strings.Trim(val, "{}")
		}
	}

	if conjunction == "<" {
		// もしもpackagesにnameが含まれていたら、そのパッケージを実行
		package_name := strings.Split(name, ".")[0]
		val, ok := packages[package_name]
		if ok {
			val(strings.Split(name, ".")[1], value, value_in_quotes, variables)
		}

		// もしもpackages_with_functionsにnameが含まれていたら、そのパッケージを実行
		val2, ok2 := packages_with_functions[package_name]
		if ok2 {
			val2(strings.Split(name, ".")[1], value, value_in_quotes, variables, functions)
		}

		if name == "print" {
			fmt.Print(giveSymbols(variables_replacer(variables, value[0], value_in_quotes[0], false)))
		}

		if name == "printf" {
			fmt.Print(giveSymbols(variables_replacers((*variables), value[0], value[1:], value_in_quotes[1:])))
		}

		if name == "println" {
			fmt.Println(giveSymbols(variables_replacer(variables, value[0], value_in_quotes[0], false)))
		}

		if name == "if" {
			for i, val := range(value) {
				if strings.HasPrefix(val, "#") {
					func_name := strings.Split(val, "(")[0][1:]
					re := regexp.MustCompile(`^#[^\(]*\(|\)$`)
					args := split(re.ReplaceAllString(val, ""))
					args_str, args_in_quote := divide_split(args)
					value[i], value_in_quotes[i] = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)
				}
			}

			if len(value) == 4 && value[2] == "else" {
				if value[0] == "true" || value[0] == "1" {
					// value[1]と[3] は {} で囲まれた部分。[2]はelse
					codes, _ := divide_split(split(value[1]))
					codes_str := strings.Join(codes, " ")
					parsed_codes := splitOutsideSemicolons(codes_str)

					for _, code := range(parsed_codes) {
						p := parser(code)
						status := p.runner(variables, functions, sharps, before_func_name, false)
						if status == 1 {
							return 1
						} else if status == 2 {
							return 2
						} else if status == 3 {
							return 3
						}
					}
				} else if value[0] == "false" || value[0] == "0" {
					// value[1]と[3] は {} で囲まれた部分。[2]はelse
					codes, _ := divide_split(split(value[3]))
					codes_str := strings.Join(codes, " ")
					parsed_codes := splitOutsideSemicolons(codes_str)
					
					for _, code := range(parsed_codes) {
						p := parser(code)
						status := p.runner(variables, functions, sharps, before_func_name, false)
						if status == 1 {
							return 1
						} else if status == 2 {
							return 2
						} else if status == 3 {
							return 3
						}
					}
				}
			} else if value[0] == "true" || value[0] == "1" {
				codes, _ := divide_split(split(value[1]))
				codes_str := strings.Join(codes, " ")
				parsed_codes := splitOutsideSemicolons(codes_str)

				for _, code := range(parsed_codes) {
					p := parser(code)
					status := p.runner(variables, functions, sharps, before_func_name, false)
					if status == 1 {
						return 1
					} else if status == 2 {
						return 2
					} else if status == 3 {
						return 3
					}
				}
			} else if len(value) >= 5 {
				for i := 1; i < len(value); i += 1 {
					if value[i] == "elif" {
						if value[i + 1] == "true" || value[i + 1] == "1" {
							codes, _ := divide_split(split(value[i + 2]))
							codes_str := strings.Join(codes, " ")
							parsed_codes := splitOutsideSemicolons(codes_str)

							for _, code := range(parsed_codes) {
								p := parser(code)
								status := p.runner(variables, functions, sharps, before_func_name, false)
								if status == 1 {
									return 1
								} else if status == 2 {
									return 2
								} else if status == 3 {
									return 3
								}
							}

							break
						}
					} else if value[i] == "else" {
						codes, _ := divide_split(split(value[i + 1]))
						codes_str := strings.Join(codes, " ")
						parsed_codes := splitOutsideSemicolons(codes_str)

						for _, code := range(parsed_codes) {
							p := parser(code)
							status := p.runner(variables, functions, sharps, before_func_name, false)
							if status == 1 {
								return 1
							} else if status == 2 {
								return 2
							} else if status == 3 {
								return 3
							}
						}

						break
					}
				}
			}
		}

		if name == "fnc" {
			// value[0] は関数名、value[1]は引数、value[2]は関数の中身
			(*functions)[value[0]] = []string{value[1], value[2]}
		}

		if name == "sharp" {
			// value[0] はシャープ関数名、value[1]は引数、value[2]は関数の中身
			(*sharps)[value[0]] = []string{value[1], value[2]}
		}

		if name == "while" {
			// value[0] は条件式、value[1] は中身
			// value[0]に#を含む関数があった場合は、その関数を実行してから条件式を評価する
			var sharp string = ""
			var value_mem string = value[0]

			if strings.HasPrefix(value[0], "#") {
				func_name := strings.Split(value[0], "(")[0][1:]
				args := split(strings.Trim(strings.Split(value[0], "(")[1], ")"))
				args_str, args_in_quote := divide_split(args)
				sharp, _ = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)
			}

			if sharp != "" {
				value[0] = sharp
			}

			eval := calc_expression(variables_replacer(variables, value[0], value_in_quotes[0], false), map[string]interface{}{})
			
			var control string = ""
			for eval == "true" || eval == "1" {
				codes := splitOutsideSemicolons(value[1])

				for _, code := range(codes) {
					p := parser(code)
					status := p.runner(variables, functions, sharps, before_func_name, false)
					if status == 2 {
						control = "break"
						break
					} else if status == 3 {
						control = "continue"
						break
					} else if status == 1 {
						return 1
					}
				}

				if control == "break" {
					break
				} else if control == "continue" {
					control = ""
					continue
				}

				// 条件式を再評価
				value[0] = value_mem
				sharp = ""
				if strings.HasPrefix(value[0], "#") {
					func_name := strings.Split(value[0], "(")[0][1:]
					args := split(strings.Trim(strings.Split(value[0], "(")[1], ")"))
					args_str, args_in_quote := divide_split(args)

					sharp, _ = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)
				}

				if sharp != "" {
					value[0] = sharp
				}

				eval = calc_expression(variables_replacer(variables, value[0], value_in_quotes[0], false), map[string]interface{}{})
			}
		}

		if name == "each" {
			// value[0] は配列名、value[2] は回している配列の値を入れる変数名、value[3] は実行する中身
			// value[1] は繰り返しの種類を識別するための記号
			// 配列名が存在するか確認
			var control string = ""

			if value[1] == ">" {
				val, ok := (*variables)[value[0]]
				if !ok {
					// 無くても、配列の形をしていれば良い
					if len(strings.Split(value[0], " ")) == 1 {
						fmt.Println("The error occurred in each(function). [1] / 配列名が見つからないためエラーが発生しました。")
						fmt.Println("The array name is not found.")
						return -1
					} else {
						val = value[0]
					}
				}

				// value[2]の変数名が既に存在する場合は、警告を出す
				_, ok = (*variables)[value[2]]
				if ok {
					fmt.Println("The warning issued in each(function). [1] / [警告] 変数名がすでに存在するので、上書きされます。")
					fmt.Println("The variable name is already exist.")
				}

				// 配列名が存在する場合は、その配列を取得
				array, _ := divide_split(split(val))

				for _, val := range(array) {
					(*variables)[value[2]] = take_off_quotation(val)

					codes := splitOutsideSemicolons(value[3])

					for _, code := range(codes) {
						p := parser(code)
						status := p.runner(variables, functions, sharps, before_func_name, false)
						if status == 2 {
							control = "break"
							break
						} else if status == 3 {
							control = "continue"
							break
						} else if status == 1 {
							return 1
						}
					}

					// もしbreakがあったら、eachを抜ける
					if control == "break" {
						break
					} else if control == "continue" {
						control = ""
						continue
					}
				}

				// eachが終わったら、value[2]の変数を削除
				delete((*variables), value[2])
			} else if value[1] == ":" {
				// value[0] は配列名、value[2] は回している配列のインデックスを入れる変数名、value[3] は実行する中身
				// 配列名が存在するか確認
				val, ok := (*variables)[value[0]]
				if !ok {
					// 無くても、配列の形をしていれば良い
					if len(strings.Split(value[0], " ")) == 1 {
						fmt.Println("The error occurred in each(function). [2] / 配列名が見つからないためエラーが発生しました。")
						fmt.Println("The array name is not found.")
						return -1
					} else {
						val = value[0]
					}
				}

				// value[2]の変数名が既に存在する場合は、警告を出す
				_, ok = (*variables)[value[2]]
				if ok {
					fmt.Println("The warning issued in each(function). [2] / [警告] 変数名がすでに存在するので、上書きされます。")
					fmt.Println("The variable name is already exist.")
				}

				// 配列名が存在する場合は、その配列を取得
				array, _ := divide_split(split(val))

				for i := 0; i < len(array); i++ {
					(*variables)[value[2]] = strconv.Itoa(i)

					codes := splitOutsideSemicolons(value[3])

					for _, code := range(codes) {
						p := parser(code)
						status := p.runner(variables, functions, sharps, before_func_name, false)
						if status == 2 {
							control = "break"
							break
						} else if status == 3 {
							control = "continue"
							break
						} else if status == 1 {
							return 1
						}
					}

					// もしbreakがあったら、eachを抜ける
					if control == "break" {
						break
					} else if control == "continue" {
						control = ""
						continue
					}
				}
			} else {
				fmt.Println("The error occurred in each(function). [3] / 繰り返しの種類が不明なためエラーが発生しました。")
				fmt.Println("The type of repetition is unknown.")
				return -1
			}
		}

		if name == "match" {
			// value[0] は条件式、value[1:] はcase文とそれが真だった場合の処理
			// value[0]に#を含む関数があった場合は、その関数を実行してから条件式を評価する
			var sharp string = ""

			if strings.HasPrefix(value[0], "#") {
				func_name := strings.Split(value[0], "(")[0][1:]
				args := split(strings.Trim(strings.Split(value[0], "(")[1], ")"))
				args_str, args_in_quote := divide_split(args)

				sharp, _ = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)

				if sharp != "" {
					value[0] = sharp
				}
			} else {
				value[0] = variables_replacer(variables, value[0], value_in_quotes[0], false)
			}

			// 計算できそうだったら計算する
			var eval string = ""
			mem, err := govaluate.NewEvaluableExpression(variables_replacer(variables, value[0], value_in_quotes[0], false))
			if err == nil {
				_, err2 := mem.Evaluate(map[string]interface{}{})
				if err2 == nil {
					eval = calc_expression(variables_replacer(variables, value[0], value_in_quotes[0], false), map[string]interface{}{})
				}
			}

			if eval != "" {
				value[0] = eval
			}

			var all_false bool = true
			for i := 1; i < len(value); i += 1 {
				if value[i] == "case" {
					// value[i + 1] がシャープ関数だった場合は、その関数を実行してから評価する。その中の引数に「_」があったらそこに評価対象の値を入れる
					// それ以外は変数の値と同じかどうかを評価する
					var sharped bool = false

					if strings.HasPrefix(value[i + 1], "#") {
						func_name := strings.Split(value[i + 1], "(")[0][1:]
						args := split(strings.Trim(strings.Split(value[i + 1], "(")[1], ")"))
						args_str, args_in_quote := divide_split(args)

						for j, arg := range(args_str) {
							if arg == "_" {
								args_str[j] = value[0]
							}
						}

						value[i + 1], _ = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)
						sharped = true
					} else {
						value[i + 1] = variables_replacer(variables, value[i + 1], value_in_quotes[i + 1], false)
					}

					// もしvalue[i + 1]が配列の変数だった場合は、value[i + 1]を展開する
					if !sharped {
						value[i + 1] = variables_replacer(variables, value[i + 1], value_in_quotes[i + 1], false)
					}
					
					if len(value) > i + 2 {
						// もしvalue[i + 1]が配列だった場合は、その配列の中にvalue[0]が含まれているかどうかを評価する
						if sharped && (value[i + 1] == "true" || value[i + 1] == "1") {
							// sharped = false
							all_false = false

							codes := splitOutsideSemicolons(value[i + 2])

							for _, code := range(codes) {
								p := parser(code)
								status := p.runner(variables, functions, sharps, before_func_name, false)
								if status == 1 {
									return 1
								} else if status == 2 {
									return 2
								} else if status == 3 {
									return 3
								}
							}
						} else if l, _ := divide_split(split(value[i + 1])); len(l) > 1 {
							array, _ := divide_split(split(value[i + 1]))
							if contains(array, value[0]) {
								all_false = false

								codes := splitOutsideSemicolons(value[i + 2])

								for _, code := range(codes) {
									p := parser(code)
									status := p.runner(variables, functions, sharps, before_func_name, false)
									if status == 1 {
										return 1
									} else if status == 2 {
										return 2
									} else if status == 3 {
										return 3
									}
								}
							}
						} else if value[0] == value[i + 1] {
							all_false = false

							codes := splitOutsideSemicolons(value[i + 2])

							for _, code := range(codes) {
								p := parser(code)
								status := p.runner(variables, functions, sharps, before_func_name, false)
								if status == 1 {
									return 1
								} else if status == 2 {
									return 2
								} else if status == 3 {
									return 3
								}
							}
						}
					}
				} else if value[i] == "default" {
					if all_false {
						codes := splitOutsideSemicolons(value[i + 1])

						for _, code := range(codes) {
							p := parser(code)
							status := p.runner(variables, functions, sharps, before_func_name, false)
							if status == 1 {
								return 1
							} else if status == 2 {
								return 2
							} else if status == 3 {
								return 3
							}
						}
					}
				}
			}
		}

		if name == "calc" {
			(*variables)[value[0]] = calc_expression(variables_replacers((*variables), value[1], value[2:], value_in_quotes[2:]), map[string]interface{}{})
		}

		if name == "input" {
			var input string
			fmt.Scanln(&input)
			(*variables)[value[0]] = input
		}

		if name == "vars" {
			// value[0] は変数名の配列、value[1] は変数の値の配列
			splited_names := split(value[0])
			splited_values := split(value[1])

			names, _ := divide_split(splited_names)
			values, _ := divide_split(splited_values)

			mem := split(strings.Join(values, " "))
			new_values, new_values_in_quote := divide_split(mem)

			// シャープがあった場合は、その関数を実行してから変数に格納する
			for i := 0; i < len(new_values); i += 1 {
				if strings.HasPrefix(new_values[i], "#") {
					func_name := strings.Split(new_values[i], "(")[0][1:]
					args := split(strings.Trim(strings.Split(new_values[i], "(")[1], ")"))
					args_str, args_in_quote := divide_split(args)
					
					new_values[i], new_values_in_quote[i] = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)
				}
			}

			for i, name := range(names) {
				if len(strings.Split(new_values[i], " ")) == 1 {
					(*variables)[name] = take_off_quotation(variables_replacer(variables, new_values[i], new_values_in_quote[i], false))
				} else {
					(*variables)[name] = new_values[i]
				}
			}
		}

		if name == "make" {
			if value[0] == "random" {
				// value[1]変数名、value[2]は型(int | float | array)、value[3]は開始値、value[4]は終了値
				var start int
				var end int
				var result string

				_, ok := (*variables)[value[1]]
				if !ok {
					fmt.Println("The error occurred in random in make. [1] / 変数名が見つからないためエラーが発生しました。")
					fmt.Println("The variable name is not found.")
					return -1
				}

				if len(value) == 3 {
					// それ以外はvalue[2]をarrayとみなして、その配列からランダムに1つ選ぶ
					array, _ := divide_split(split(variables_replacer(variables, value[2], value_in_quotes[2], false)))
					// arrayからランダムに1つ選ぶ
					result = array[rand.IntN(len(array))]
				} else {
					if value[2] == "int" {
						start, _ = strconv.Atoi(variables_replacer(variables, value[3], value_in_quotes[3], false))
						end, _ = strconv.Atoi(variables_replacer(variables, value[4], value_in_quotes[4], false))
						result = strconv.Itoa(rand.IntN(end - start) + start)
					} else if value[2] == "float" {
						start, _ = strconv.Atoi(variables_replacer(variables, value[3], value_in_quotes[3], false))
						end, _ = strconv.Atoi(variables_replacer(variables, value[4], value_in_quotes[4], false))
						result = strconv.FormatFloat(rand.Float64() * float64(end - start) + float64(start), 'f', 10, 64)
					} else if value[2] == "array" {
						array, _ := divide_split(split(variables_replacer(variables, value[3], value_in_quotes[3], false)))
						// arrayからランダムに1つ選ぶ
						result = array[rand.IntN(len(array))]
					}
				}

				(*variables)[value[1]] = result
			} else if value[0] == "var" {
				if len(value) >= 3 {
					// 値が複数ある場合は、それらをすべて変数名とみなし初期化する
					for _, val := range(value[1:]) {
						(*variables)[val] = ""
					}
				} else {
					// value[1]の変数名で空の変数を作成
					(*variables)[value[1]] = ""
				}
				
			}
		}

		if name == "return" {
			(*variables)["0__return__"] = variables_replacer(variables, value[0], value_in_quotes[0], false)
			return 1
		}

		if name == "runmyself" {
			// value[0]を再度パースして実行
			for _, code := range(splitOutsideSemicolons(variables_replacer(variables, value[0], value_in_quotes[0], false))) {
				p := parser(code)
				status := p.runner(variables, functions, sharps, before_func_name, false)
				if status == 1 {
					break
				} else if status == 2 {
					return 2
				} else if status == 3 {
					return 3
				}
			}
		}

		if name == "break" {
			return 2
		}

		if name == "continue" {
			return 3
		}

		if name == "delete" {
			// もし変数名が複数あったら、それらをすべて削除する
			if len(value) >= 2 {
				for _, val := range(value) {
					_, ok := (*variables)[val]
					if !ok {
						fmt.Println("The error occurred in delete. [1] / 変数名が見つからないためエラーが発生しました。")
						fmt.Println("The variable name is not found.")
						return -1
					}
					delete((*variables), val)
				}
			} else {
				// value[0] は変数名
				// 変数が存在するか確認
				_, ok := (*variables)[value[0]]
				if !ok {
					fmt.Println("The error occurred in delete. [2] / 変数名が見つからないためエラーが発生しました。")
					fmt.Println("The variable name is not found.")
					return -1
				}

				delete((*variables), value[0])
			}
		}

		if name == "swap" {
			// value[0] は変数名1、value[1] は変数名2
			// 変数の値を入れ替える
			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in swap. [1] / 変数1が見つからないためエラーが発生しました。")
				fmt.Println("The variable (1) is not found.")
				return -1
			}

			_, ok = (*variables)[value[1]]
			if !ok {
				fmt.Println("The error occurred in swap. [2] / 変数2が見つからないためエラーが発生しました。")
				fmt.Println("The variable (2) is not found.")
				return -1
			}

			(*variables)[value[0]], (*variables)[value[1]] = (*variables)[value[1]], (*variables)[value[0]]
		}

		// もしオリジナル関数名がnameに存在する場合は、その関数を実行
		// value[0]は引数の値のリスト、value[1]は「to」など、value[2]は返り値を格納する変数名
		// toがあったら、その変数に返り値を格納する
		if _, ok := (*functions)[name]; ok {
			args, _ := divide_split(split((*functions)[name][0]))
			splited_value, _ := divide_split(split(value[0]))
			
			if len(args) != 0 && len(splited_value) != 0 && len(args) == len(splited_value) {
				for i, arg := range(args) {
					args[i] = take_off_quotation(arg)
					splited_value[i] = take_off_quotation(splited_value[i])

					// もし#があったら、その関数を実行してから格納する
					if strings.HasPrefix(splited_value[i], "#") {
						func_name := strings.Split(splited_value[i], "(")[0][1:]
						re := regexp.MustCompile(`^#[^\(]*\(|\)$`)

						args := split(re.ReplaceAllString(splited_value[i], ""))
						args_str, args_in_quote := divide_split(args)
						
						splited_value[i], _ = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)
					}

					// もしsplited_value[i]に波括弧があったら、その中身を取り出す
					if strings.HasPrefix(splited_value[i], "{") && strings.HasSuffix(splited_value[i], "}") {
						splited_value[i] = strings.Trim(splited_value[i], "{}")
						splited_value[i] = (*variables)[splited_value[i]]
					}
					(*variables)[arg] = splited_value[i]
				}
			// もし引数名があるが値が設定されていなかったらその引数を空白にするという処理を追加する
			} else if len(args) != 0 && len(splited_value) == 0 {
				for _, arg := range(args) {
					(*variables)[arg] = ""
				}
			}

			// 関数の中身を実行
			codes := splitOutsideSemicolons((*functions)[name][1])
			
			for _, code := range(codes) {
				p := parser(code)
				status := p.runner(variables, functions, sharps, before_func_name, false)
				if status == 1 {
					break
				} else if status == 2 {
					return 2
				} else if status == 3 {
					return 3
				}
			}

			// 返り値を格納する変数がある場合は、その変数に格納する
			// そもそもvalue[1]が存在するか確認
			if len(value) > 2 {
				if value[1] == "to" {
					(*variables)[value[2]] = variables_replacer(variables, (*variables)["0__return__"], false, false)
				}
			}

			// 引数を削除
			for _, arg := range(args) {
				delete((*variables), arg)
			}
			delete((*variables), "0__return__")
		}

	} else if conjunction == "=" {
		// もし変数名の最初が数字だったら、エラーを返す
		if contains([]string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}, name[0:1]) {
			fmt.Println("The error occurred in variables. [1] / 変数名が不正(数字から始まっているなど)のためエラーが発生しました。")
			fmt.Println("The variable name is invalid.")
			return -1
		}
		
		if len(value) == 1 {
			if len(strings.Split(value[0], " ")) > 1 {
				array_value, array_value_in_quotes := divide_split(split(value[0]))

				for i := 0; i < len(array_value); i++ {
					// もし#があったら、その関数を実行してから格納する
					if strings.HasPrefix(array_value[i], "#") {
						func_name := strings.Split(array_value[i], "(")[0][1:]
						args := split(strings.Trim(strings.Split(array_value[i], "(")[1], ")"))
						args_str, args_in_quote := divide_split(args)

						array_value[i], array_value_in_quotes[i] = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)
					}

					array_value[i] = variables_replacer(variables, array_value[i], array_value_in_quotes[i], true)
				}

				array_value_edited := strings.Join(array_value, " ")

				(*variables)[name] = array_value_edited
			} else {
				(*variables)[name] = variables_replacer(variables, value[0], value_in_quotes[0], false)
			}
		} else if value[0] == "value" {
			(*variables)[name] = variables_replacer(variables, value[1], value_in_quotes[1], false)
		} else if value[0] == "array" {
			array_value, array_value_in_quotes := divide_split(split(value[1]))

			for i := 0; i < len(array_value); i++ {
				// もし#があったら、その関数を実行してから格納する
				if strings.HasPrefix(array_value[i], "#") {
					func_name := strings.Split(array_value[i], "(")[0][1:]
					args := split(strings.Trim(strings.Split(array_value[i], "(")[1], ")"))
					args_str, args_in_quote := divide_split(args)

					array_value[i], array_value_in_quotes[i] = sharp_functions(func_name, args_str, args_in_quote, variables, functions, sharps)
				}

				array_value[i] = variables_replacer(variables, array_value[i], array_value_in_quotes[i], true)
			}

			array_value_edited := strings.Join(array_value, " ")

			(*variables)[name] = array_value_edited
		} else if value[0] == "add" {
			var variables_value_float, value1_float float64
			var variables_value, value1 int
			var err1, err2, err3, err4 error

			// 数値に変換して足し算してから文字列に変換して格納
			variables_value, err1 = strconv.Atoi((*variables)[name])
			if err1 != nil {
				// 少数として格納されている場合
				variables_value_float, err2 = strconv.ParseFloat((*variables)[name], 64)
				if err2 != nil {
					fmt.Println("The error occurred in variables. [1] / 代入先の変数が数値でないためエラーが発生しました。")
					fmt.Println(err2)
					return -1
				}
			}

			value1, err3 = strconv.Atoi(variables_replacer(variables, value[1], value_in_quotes[1], false))
			if err3 != nil {
				// 小数として格納されている場合
				value1_float, err4 = strconv.ParseFloat(variables_replacer(variables, value[1], value_in_quotes[1], false), 64)
				if err4 != nil {
					fmt.Println("The error occurred in variables. [2] / 変数が数値でないためエラーが発生しました。")
					fmt.Println(err4)
					return -1
				}
			}

			if err1 == nil && err3 == nil {
				result := variables_value + value1
				(*variables)[name] = strconv.Itoa(result)
			} else if err2 == nil && err4 == nil {
				result := value1_float + variables_value_float
				(*variables)[name] = strconv.FormatFloat(result, 'f', 10, 64)
			}
		}
	}

	if top {
		*before_func_name = name
	}

	return 0
}

func splitOutsideSemicolons(input string) []string {
	// 一番外側にあるセミコロンで分割する
	// ただし、波括弧内にあるセミコロンは無視し、波括弧内をひとかたまりとして扱う

	var result []string
	var count int = 0
	var mem string = ""

	// \;を__SEMICOLON__に変換
	re := regexp.MustCompile(`\\;`)
	input = re.ReplaceAllString(input, "__SEMICOLON__")

	for _, val := range(strings.Split(input, ";")) {
		count += strings.Count(val, "{")
		count -= strings.Count(val, "}")

		if count == 0 {
			result = append(result, mem + val)
			mem = ""
		} else {
			mem += val + ";"
		}
	}

	for i := 0; i < len(result); i++ {
		// 空白でなく改行で引数が区切られている場合は、空白に変換する
		re := regexp.MustCompile(`\n`)
		result[i] = re.ReplaceAllString(result[i], " ")

		// __SEMICOLON__を;に変換
		re = regexp.MustCompile(`__SEMICOLON__`)
		result[i] = re.ReplaceAllString(result[i], ";")
	}

	return result
}

func main() {
	// コードが書かれたファイルを読み込む
	// コマンドラインでファイル名を指定する
	// 例: このファイル run --path ファイル名.wg
	var code string = ""
	var path string

	args := os.Args
	if len(args) == 1 {
		fmt.Println("The error occurred in main. [2] / コマンドが指定されていないためエラーが発生しました。")
		fmt.Println("The command is not found.")
		return
	}

	// もし run なら
	if args[1] == "run" {
		path = args[2]
		// ファイルを読み込む
		file, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("The error occurred in main. [3] / ファイルが読み込めないためエラーが発生しました。")
			fmt.Println(err)
			return
		}

		code = string(file)
	}

	variables := make(map[string]string)
	functions := make(map[string][]string)
	sharps := make(map[string][]string)

	var before_func_name string = ""

	for _, val := range(splitOutsideSemicolons(code)) {
		parsed := parser(val)
		status := parsed.runner(&variables, &functions, &sharps, &before_func_name, true)
		if status == -1 {
			break
		}
	}
}