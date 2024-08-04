package main

import (
	"fmt"
	"strings"
	"regexp"
	"strconv"
	"math/rand/v2"
	"os"

	date_pack "github.com/kurogee/wagyu-script/date"
	get_pack "github.com/kurogee/wagyu-script/get"
	file_pack "github.com/kurogee/wagyu-script/file"
	array_pack "github.com/kurogee/wagyu-script/array"
	string_pack "github.com/kurogee/wagyu-script/string"

	math_sharp_functions "github.com/kurogee/wagyu-script/maths"

	"github.com/Knetic/govaluate"
)

type Parse struct {
	parsed []string
	parsed_meaning CodeType
}

type CodeType struct {
	name string
	conjunction string
	value []string
}

type Package map[string]func(string, []string, *map[string]string)
type Package_with_functions map[string]func(string, []string, *map[string]string, *map[string][]string)
type Package_sharp_functions map[string]func(string, []string, *map[string]string) string

// packageの一覧を定義
var packages = Package{
	"date": date_pack.Run,
	"file": file_pack.Run,
	"array": array_pack.Run,
	"string": string_pack.Run,
}

var packages_with_functions = Package_with_functions{
	"get": get_pack.Run,
}

var packages_sharp_functions = Package_sharp_functions{
	"math": math_sharp_functions.Sharp,
}

func contains(s []string, e string) bool {
	for _, a := range(s) {
		if a == e {
			return true
		}
	}

	return false
}

func split(input string) []string {
	var tokens []string
	var buffer strings.Builder
	var quoteChar rune
	var parenStack []rune

	var haveSpace bool = false

	// 違う種類の括弧が重なっているところがあったら、それを離す
	// 例: ({ → ( { にする など
	re := regexp.MustCompile(`([\(\{\[\)\}\]])\s+([\(\{\[\)\}\]])`)
	if re.MatchString(input) {
		haveSpace = true
	}

	input = re.ReplaceAllString(input, "$1 $2")

	// \\;という文字列を__SEMICOLON__に置き換える
	re = regexp.MustCompile(`\\;`)
	input = re.ReplaceAllString(input, "__SEMICOLON__")

	inQuote := false
	inParen := false
	inFunctionCall := false

	// Regular expressions for matching parentheses and braces
	parenRegex := regexp.MustCompile(`^[\(\)\{\}]$`)
	functionCallRegex := regexp.MustCompile(`^#\w*\(`)

	for _, char := range input {
		switch {
		case inQuote:
			buffer.WriteRune(char)
			if char == quoteChar {
				inQuote = false
				tokens = append(tokens, buffer.String())
				buffer.Reset()
			}

		case inParen:
			buffer.WriteRune(char)
			if char == parenStack[len(parenStack)-1] && len(parenStack) > 0 {
				parenStack = parenStack[:len(parenStack)-1]
				if len(parenStack) == 0 {
					inParen = false
					if inFunctionCall {
						tokens = append(tokens, buffer.String())
						buffer.Reset()
						inFunctionCall = false
					}
				}
			} else if parenRegex.MatchString(string(char)) {
				parenStack = append(parenStack, matchingParen(char))
			}

		case inFunctionCall:
			buffer.WriteRune(char)
			if char == '(' {
				parenStack = append(parenStack, ')')
				inParen = true
			} else if char == ')' && len(parenStack) == 0 {
				inFunctionCall = false
				tokens = append(tokens, buffer.String())
				buffer.Reset()
			}

		default:
			switch char {
			case ' ', '\t', '\n':
				if buffer.Len() > 0 {
					tokens = append(tokens, buffer.String())
					buffer.Reset()
				}

			case '\'', '"':
				inQuote = true
				quoteChar = char
				buffer.WriteRune(char)

			case '(', '{':
				inParen = true
				parenStack = append(parenStack, matchingParen(char))
				buffer.WriteRune(char)

			default:
				buffer.WriteRune(char)
				if functionCallRegex.MatchString(buffer.String()) {
					inFunctionCall = true
				}
			}
		}
	}

	if buffer.Len() > 0 {
		tokens = append(tokens, buffer.String())
	}

	if haveSpace {
		// #関数名 と そのすぐ後にある()を結合する
		for i, token := range(tokens) {
			if strings.HasPrefix(token, "#") {
				if i+1 < len(tokens) {
					if strings.HasPrefix(tokens[i+1], "(") {
						tokens[i] = token + tokens[i+1]
						tokens = append(tokens[:i+1], tokens[i+2:]...)
					}
				}
			}
		}

		return tokens
	}

	return tokens
}

func matchingParen(char rune) rune {
	switch char {
	case '(':
		return ')'
	case '{':
		return '}'
	case ')':
		return '('
	case '}':
		return '{'
	}
	return char
}

/*
func replaceSymbols(input string) string {
	input = strings.ReplaceAll(input, "\n", "\\n")
	input = strings.ReplaceAll(input, "\t", "\\t")
	input = strings.ReplaceAll(input, "\\", "\\\\")
	input = strings.ReplaceAll(input, "\"", "\\\"")
	input = strings.ReplaceAll(input, "'", "\\'")

	return input
}*/

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
	
	if mem == "" {
		p.parsed = []string{"None", ""}
		p.parsed_meaning = CodeType{"None", "", []string{}}
		return
	}

	// 先頭に//があったらコメントアウトを削除
	if strings.HasPrefix(code, "//") {
		p.parsed = []string{"None", ""}
		p.parsed_meaning = CodeType{"None", "", []string{}}
		return
	}

	divided_code := split(code)

	p.parsed = divided_code

	// もしdivided_code[1]が記号ではなかったら、divided_code[1:]をvalueに格納しconjunctionには<を格納
	re = regexp.MustCompile(`^[><=]$`)
	if re.MatchString(divided_code[1]) {
		p.parsed_meaning.name = strings.ReplaceAll(divided_code[0], " ", "")
		p.parsed_meaning.conjunction = divided_code[1]
		p.parsed_meaning.value = divided_code[2:]
	} else {
		p.parsed_meaning.name = strings.ReplaceAll(divided_code[0], " ", "")
		p.parsed_meaning.conjunction = "<"
		p.parsed_meaning.value = divided_code[1:]
	}

	return
}

func variables_replacer(variables *map[string]string, target string) string {
	val, ok := (*variables)[target]
	if ok {
		return val
	}

	return target
}

func take_off_quotation(target string) string {
	if strings.HasPrefix(target, "'") && strings.HasSuffix(target, "'") {
		// 一番外側のシングルクォーテーションを取り除く
		re := regexp.MustCompile(`^'|'$`)
		return re.ReplaceAllString(target, "")
	} else if strings.HasPrefix(target, "\"") && strings.HasSuffix(target, "\"") {
		re := regexp.MustCompile(`^"|"$`)
		return re.ReplaceAllString(target, "")
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

func sharp_functions(func_name string, args []string, variables *map[string]string) string {
	if variables == nil {
		variables = &map[string]string{}
	}

	reformat_args := strings.Join(args, " ")
	args = split(reformat_args)

	// // もしargsのどこかに波括弧がり、かつその中身が計算式だったら、その計算式を計算してから格納する
	// for i, arg := range(args) {
	// 	if strings.HasPrefix(arg, "{") && strings.HasSuffix(arg, "}") {
	// 		// 波括弧を取り除く
	// 		arg = strings.Trim(arg, "{}")
	// 
	// 		// 一度空白で区切り、変数名があったら、その変数名の値を取得
	// 		divided_arg := split(arg)
	// 		for i, arg := range(divided_arg) {
	// 			divided_arg[i] = variables_replacer(variables, arg)
	// 		}
	// 
	// 		arg = strings.Join(divided_arg, " ")
	// 
	// 		if strings.Contains(arg, "+") || strings.Contains(arg, "-") || strings.Contains(arg, "*") || strings.Contains(arg, "/") {
	// 			args[i] = calc_expression(variables_replacer(variables, arg))
	// 		}
	// 	}
	// }
	
	if func_name != "" {
		for i, arg := range(args) {
			args[i] = take_off_quotation(arg)
		}
	}

	// もしfunc_nameに.が含まれていたら、そのパッケージを実行
	if strings.Contains(func_name, ".") {
		package_name := strings.Split(func_name, ".")[0]
		val, ok := packages_sharp_functions[package_name]
		if !ok {
			fmt.Println("The error occurred in sharp_functions. [1] / パッケージが見つからないためエラーが発生しました。")
			fmt.Println("The package is not found.")
			return ""
		}

		return val(strings.Split(func_name, ".")[1], args, variables)
	}

	if func_name == "arrayAt" {
		// args[0] は配列名、args[1] はインデックス
		// インデックスは数値である必要がある
		index, err := strconv.Atoi(variables_replacer(variables, args[1]))
		if err != nil {
			fmt.Println("The error occurred in arrayAt. [1] / インデックスが数値でないためエラーが発生しました。")
			fmt.Println(err)
		}

		// 配列名が存在するか確認
		_, ok := (*variables)[args[0]]
		if !ok {
			fmt.Println("The error occurred in arrayAt. [2] / 配列名が見つからないためエラーが発生しました。")
			fmt.Println("The array name is not found.")
			return ""
		}

		// 配列名が存在する場合は、その配列を取得
		array := split((*variables)[args[0]])

		// インデックスが配列の範囲内にあるか確認
		if index < 0 || index >= len(array) {
			fmt.Println("The error occurred in arrayAt. [3] / インデックスが範囲外のためエラーが発生しました。")
			fmt.Println("The index is out of range.")
			return ""
		}

		return array[index]
	} else if func_name == "arrayLen" {
		// args[0] は配列名
		// 配列名が存在するか確認
		_, ok := (*variables)[args[0]]
		if !ok {
			fmt.Println("The error occurred in arrayLen. [1] / 配列名が見つからないためエラーが発生しました。")
			fmt.Println("The array name is not found.")
			return ""
		}

		// 配列名が存在する場合は、その配列を取得
		array := strings.Split((*variables)[args[0]], " ")

		return strconv.Itoa(len(array))
	} else if func_name == "calc" {
		// args[0] は計算式
		// args[1:] は変数名

		if len(args) == 1 {
			return calc_expression(variables_replacer(variables, args[0]), map[string]interface{}{})
		} else if len(args) > 1 {
			return calc_expression(variables_replacers((*variables), args[0], args[1:]), map[string]interface{}{})
		}
	} else if func_name == "from" {
		// 引数の1番目の数字から2番目の数字までの連番の配列を作成
		// args[0] は開始番号、args[1] は終了番号

		// もし引数が変数名である場合は、その変数名の値を取得
		start_arg := variables_replacer(variables, args[0])
		end_arg := variables_replacer(variables, args[1])

		// fmt.Println(start_arg, end_arg)

		start, err := strconv.Atoi(start_arg)
		if err != nil {
			fmt.Println("The error occurred in from. [1] / 開始番号が(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
			return ""
		}

		end, err := strconv.Atoi(end_arg)
		if err != nil {
			fmt.Println("The error occurred in from. [2] / 終了番号が(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
			return ""
		}

		var result string = " "
		for i := start; i <= end; i++ {
			result += strconv.Itoa(i) + " "
		}

		// 最後のスペースを削除
		result = result[:len(result) - 1]

		return result
	} else if func_name == "format" {
		// args[0]は接続する文字列フォーマット、args[1:]は変数名
		// printfと同じように変数を埋め込む
		return variables_replacers((*variables), args[0], args[1:])
	} else if func_name == "var" {
		val, ok := (*variables)[args[0]]
		if ok {
			return val
		}

		return ""
	} else if func_name == "remove" {
		// args[0] は配列名、args[1] は削除する値のインデックス
		_, ok := (*variables)[args[0]]
		if !ok {
			fmt.Println("The error occurred in remove. [1] / 配列名が見つからないためエラーが発生しました。")
			fmt.Println("The array name is not found.")
		}

		// 配列名が存在する場合は、その配列を取得
		array := strings.Split((*variables)[args[0]], " ")

		// インデックスが数値であるか確認
		index, err := strconv.Atoi(variables_replacer(variables, args[1]))
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
		return strings.Join(array, " ")
	} else if func_name == "rand" {
		// args[0] は開始値、args[1] は終了値
		start, err := strconv.Atoi(variables_replacer(variables, args[0]))
		if err != nil {
			fmt.Println("The error occurred in rand. [1] / 開始値が(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
		}

		end, err := strconv.Atoi(variables_replacer(variables, args[1]))
		if err != nil {
			fmt.Println("The error occurred in rand. [2] / 終了値が(変数の中の値が)数値でないためエラーが発生しました。")
			fmt.Println(err)
		}

		return strconv.Itoa(rand.IntN(end - start) + start)
	} else if func_name == "convert" {
		// args[0] は変換する値、args[1] は変換後の型
		switch args[1] {
		case "int":
			// 小数点があったら、それを取り除く
			val, err := strconv.Atoi(strings.Split(variables_replacer(variables, args[0]), ".")[0])
			if err != nil {
				fmt.Println("The error occurred in convert. [1] / 変換する値が(変数の中の値が)数値でないためエラーが発生しました。")
				fmt.Println(err)
			}

			return strconv.Itoa(val)
		case "float":
			val, err := strconv.ParseFloat(variables_replacer(variables, args[0]), 64)
			if err != nil {
				fmt.Println("The error occurred in convert. [2] / 変換する値が(変数の中の値が)数値でないためエラーが発生しました。")
				fmt.Println(err)
			}

			return strconv.FormatFloat(val, 'f', 10, 64)
		}
	} else if func_name == "" {
		for i, arg := range(args) {
			args[i] = variables_replacer(variables, arg)
		}

		var formula string = strings.Join(args, " ")

		return calc_expression(formula, map[string]interface{}{})
	}

	return ""
}

func (p Parse) runner(variables *map[string]string, functions *map[string][]string, before_func_name *string, top bool) int {
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

    if name == "None" {
        return 0
    }

    if name == "$" {
        code := *before_func_name + " " + conjunction + " " + strings.Join(value, " ")
        parsedCode := parser(code)
        status := parsedCode.runner(variables, functions, before_func_name, false)
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
		if name != "while" && name != "if" {
			if strings.HasPrefix(val, "#") {
				func_name := strings.Split(val, "(")[0][1:]
				re := regexp.MustCompile(`^#[^\(]*\(|\)$`)
				args := split(re.ReplaceAllString(val, ""))
				value[i] = sharp_functions(func_name, args, variables)
			}
		}
    }
	
    mem := &value
    not_remove_list := []string{"if"}
    if !contains(not_remove_list, name) {
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
			val(strings.Split(name, ".")[1], value, variables)
		}

		// もしもpackages_with_functionsにnameが含まれていたら、そのパッケージを実行
		val2, ok2 := packages_with_functions[package_name]
		if ok2 {
			val2(strings.Split(name, ".")[1], value, variables, functions)
		}

		if name == "print" {
			fmt.Print(take_off_quotation(giveSymbols(variables_replacer(variables, value[0]))))
		}

		if name == "printf" {
			fmt.Print(take_off_quotation(giveSymbols(variables_replacers((*variables), value[0], value[1:]))))
		}

		if name == "println" {
			fmt.Println(take_off_quotation(giveSymbols(variables_replacer(variables, value[0]))))
		}

		if name == "if" {
			for i, val := range(value) {
				if strings.HasPrefix(val, "#") {
					func_name := strings.Split(val, "(")[0][1:]
					re := regexp.MustCompile(`^#[^\(]*\(|\)$`)
					args := split(re.ReplaceAllString(val, ""))
					value[i] = sharp_functions(func_name, args, variables)
				}
			}

			if len(value) == 4 && value[2] == "else" {
				if value[0] == "true" || value[0] == "1" {
					// value[1]と[3] は {} で囲まれた部分。[2]はelse
					codes := strings.Join(split(value[1]), " ")
					parsed_codes := splitOutsideSemicolons(codes)

					for _, code := range(parsed_codes) {
						p := parser(code)
						status := p.runner(variables, functions, before_func_name, false)
						if status == 1 {
							break
						} else if status == 2 {
							return 2
						} else if status == 3 {
							return 3
						}
					}
				} else if value[0] == "false" || value[0] == "0" {
					// value[1]と[3] は {} で囲まれた部分。[2]はelse
					codes := strings.Join(split(value[3]), " ")
					parsed_codes := splitOutsideSemicolons(codes)
					
					for _, code := range(parsed_codes) {
						p := parser(code)
						status := p.runner(variables, functions, before_func_name, false)
						if status == 1 {
							break
						} else if status == 2 {
							return 2
						} else if status == 3 {
							return 3
						}
					}
				}
			} else if value[0] == "true" || value[0] == "1" {
				codes := strings.Join(split(value[1]), " ")
				parsed_codes := splitOutsideSemicolons(codes)

				for _, code := range(parsed_codes) {
					p := parser(code)
					status := p.runner(variables, functions, before_func_name, false)
					if status == 1 {
						break
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
							codes := strings.Join(split(value[i + 2]), " ")
							parsed_codes := splitOutsideSemicolons(codes)

							for _, code := range(parsed_codes) {
								p := parser(code)
								status := p.runner(variables, functions, before_func_name, false)
								if status == 1 {
									break
								} else if status == 2 {
									return 2
								} else if status == 3 {
									return 3
								}
							}

							break
						}
					} else if value[i] == "else" {
						codes := strings.Join(split(value[i + 1]), " ")
						parsed_codes := splitOutsideSemicolons(codes)

						for _, code := range(parsed_codes) {
							p := parser(code)
							status := p.runner(variables, functions, before_func_name, false)
							if status == 1 {
								break
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

		if name == "while" {
			// value[0] は条件式、value[1] は中身
			// value[0]に#を含む関数があった場合は、その関数を実行してから条件式を評価する
			var sharp string = ""
			var value_mem string = value[0]

			if strings.HasPrefix(value[0], "#") {
				func_name := strings.Split(value[0], "(")[0][1:]
				args := split(strings.Trim(strings.Split(value[0], "(")[1], ")"))
				sharp = sharp_functions(func_name, args, variables)
			}

			if sharp != "" {
				value[0] = sharp
			}

			eval := calc_expression(variables_replacer(variables, value[0]), map[string]interface{}{})
			
			var control string = ""
			for eval == "true" || eval == "1" {
				codes := splitOutsideSemicolons(value[1])

				for _, code := range(codes) {
					p := parser(code)
					status := p.runner(variables, functions, before_func_name, false)
					if status == 2 {
						control = "break"
						break
					} else if status == 3 {
						control = "continue"
						break
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
					sharp = sharp_functions(func_name, args, variables)
				}

				if sharp != "" {
					value[0] = sharp
				}

				eval = calc_expression(variables_replacer(variables, value[0]), map[string]interface{}{})
			}
		}

		if name == "calc" {
			(*variables)[value[0]] = calc_expression(variables_replacers((*variables), value[1], value[2:]), map[string]interface{}{})
		}

		if name == "input" {
			var input string
			fmt.Scanln(&input)
			(*variables)[value[0]] = input
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
					array := strings.Split(variables_replacer(variables, value[2]), " ")
					// arrayからランダムに1つ選ぶ
					result = array[rand.IntN(len(array))]
				} else {
					if value[2] == "int" {
						start, _ = strconv.Atoi(variables_replacer(variables, value[3]))
						end, _ = strconv.Atoi(variables_replacer(variables, value[4]))
						result = strconv.Itoa(rand.IntN(end - start) + start)
					} else if value[2] == "float" {
						start, _ = strconv.Atoi(variables_replacer(variables, value[3]))
						end, _ = strconv.Atoi(variables_replacer(variables, value[4]))
						result = strconv.FormatFloat(rand.Float64() * float64(end - start) + float64(start), 'f', 10, 64)
					} else if value[2] == "array" {
						array := strings.Split(variables_replacer(variables, value[3]), " ")
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

		if name == "each" {
			// value[0] は配列名、value[1] は回している配列の値を入れる変数名、value[2] は実行する中身
			// 配列名が存在するか確認
			var control string = ""

			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in each(function). [1] / 配列名が見つからないためエラーが発生しました。")
				fmt.Println("The array name is not found.")
				return -1
			}

			// value[1]の変数名が既に存在する場合は、警告を出す
			_, ok = (*variables)[value[1]]
			if ok {
				fmt.Println("The warning issued in each(function). [1] / [警告] 変数名がすでに存在するので、上書きされます。")
				fmt.Println("The variable name is already exist.")
			}

			// 配列名が存在する場合は、その配列を取得
			array := strings.Split((*variables)[value[0]], " ")

			for _, val := range(array) {
				(*variables)[value[1]] = take_off_quotation(val)

				codes := splitOutsideSemicolons(value[2])

				for _, code := range(codes) {
					p := parser(code)
					status := p.runner(variables, functions, before_func_name, false)
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

			// eachが終わったら、value[1]の変数を削除
			delete((*variables), value[1])
		}

		if name == "return" {
			(*variables)["0__return__"] = value[0]
			return 1
		}

		if name == "runmyself" {
			// value[0]を再度パースして実行
			for _, code := range(splitOutsideSemicolons(variables_replacer(variables, value[0]))) {
				p := parser(code)
				status := p.runner(variables, functions, before_func_name, false)
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
			// value[0] は変数名
			// 変数が存在するか確認
			_, ok := (*variables)[value[0]]
			if !ok {
				fmt.Println("The error occurred in delete. [1] / 変数名が見つからないためエラーが発生しました。")
				fmt.Println("The variable name is not found.")
				return -1
			}
			
			delete((*variables), value[0])
		}

		// もしオリジナル関数名がnameに存在する場合は、その関数を実行
		// value[0]は引数の値のリスト、value[1]は「to」など、value[2]は返り値を格納する変数名
		// toがあったら、その変数に返り値を格納する
		if _, ok := (*functions)[name]; ok {
			args := strings.Split((*functions)[name][0], " ")
			splited_value := split(value[0])
			
			if len(args) != 0 && len(splited_value) != 0 && len(args) == len(splited_value) {
				for i, arg := range(args) {
					args[i] = take_off_quotation(arg)
					splited_value[i] = take_off_quotation(splited_value[i])

					// もし#があったら、その関数を実行してから格納する
					if strings.HasPrefix(splited_value[i], "#") {
						func_name := strings.Split(splited_value[i], "(")[0][1:]
						re := regexp.MustCompile(`^#[^\(]*\(|\)$`)
						args := split(re.ReplaceAllString(splited_value[i], ""))
						splited_value[i] = sharp_functions(func_name, args, variables)
					}

					// もしsplited_value[i]に波括弧があったら、その中身を取り出す
					if strings.HasPrefix(splited_value[i], "{") && strings.HasSuffix(splited_value[i], "}") {
						splited_value[i] = strings.Trim(splited_value[i], "{}")
						splited_value[i] = (*variables)[splited_value[i]]
					}
					(*variables)[arg] = splited_value[i]
				}
			}

			// 関数の中身を実行
			codes := splitOutsideSemicolons((*functions)[name][1])
			
			for _, code := range(codes) {
				p := parser(code)
				status := p.runner(variables, functions, before_func_name, false)
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
					(*variables)[value[2]] = variables_replacer(variables, (*variables)["0__return__"])
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

		if value[0] == "value" {
			(*variables)[name] = variables_replacer(variables, value[1])
		} else if value[0] == "array" {
			array_value := split(value[1])
			for i := 1; i < len(array_value); i++ {
				array_value[i] = variables_replacer(variables, array_value[i])
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

			value1, err3 = strconv.Atoi(variables_replacer(variables, value[1]))
			if err3 != nil {
				// 少数として格納されている場合
				value1_float, err4 = strconv.ParseFloat(variables_replacer(variables, value[1]), 64)
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

	// \\;という文字列を__SEMICOLON__に置き換える
	input = strings.ReplaceAll(input, "\\;", "__SEMICOLON__")

	var result []string
	var count int = 0
	var mem string = ""
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

	// セミコロンを元に戻す
	for i, val := range(result) {
		result[i] = strings.ReplaceAll(val, "__SEMICOLON__", ";")
	}

	return result
}

func main() {
	// コードが書かれたファイルを読み込む
	// コマンドラインでファイル名を指定する
	// 例: このファイル run --path ファイル名.wg
	var code string = ""
	var path string

	// もし run なら
	args := os.Args
	if len(args) == 1 {
		fmt.Println("The error occurred in main. [2] / コマンドが指定されていないためエラーが発生しました。")
		fmt.Println("The command is not found.")
		return
	}

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
	var before_func_name string = ""

	for _, val := range(splitOutsideSemicolons(code)) {
		parsed := parser(val)
		status := parsed.runner(&variables, &functions, &before_func_name, true)
		if status == -1 {
			break
		}
	}
}