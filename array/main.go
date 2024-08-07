package array

import (
	"fmt"
	"strings"
	"strconv"
	"regexp"
)

var replace_chars map[string]string = map[string]string{
	"\\;" : "__SEMICOLON__",
	"\\\"" : "__DOUBLE_QUOTATION__",
	"\\'" : "__SINGLE_QUOTATION__",
	"\\`" : "__BACK_QUOTATION__",
};

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

func split(input string) []map[string]bool {
	var tokens []string
	var buffer strings.Builder
	var quoteChar rune
	var parenStack []rune
	
	// 違う種類の括弧が重なっているところがあったら、それを離す
	// 例: ({ → ( { にする など
	re := regexp.MustCompile(`([\(\{\[\)\}\]])\s+([\(\{\[\)\}\]])`)
	input = re.ReplaceAllString(input, "$1 $2")

	// \\～を一時的にreplace_charsに置き換える
	for key, value := range(replace_chars) {
		input = strings.ReplaceAll(input, key, value)
	}

	inQuote := false
	inQuotes := []bool{}

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

				inQuotes = append(inQuotes, true)

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

						inQuotes = append(inQuotes, false)
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

				inQuotes = append(inQuotes, false)
			}

		default:
			switch char {
			case ' ', '\t', '\n':
				if buffer.Len() > 0 {
					tokens = append(tokens, buffer.String())
					buffer.Reset()

					inQuotes = append(inQuotes, false)
				}

			case '\'', '"', '`':
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

		if inQuote {
			inQuotes = append(inQuotes, true)
		} else {
			inQuotes = append(inQuotes, false)
		}
	}
	
	// replace_charsを元に戻す
	for i := 0; i < len(tokens); i++ {
		for key, value := range(replace_chars) {
			tokens[i] = strings.ReplaceAll(tokens[i], value, key)
		}
	}

	// tokensをmapに変換し、inQuotesと合わせる
	var tokens_map []map[string]bool
	for i, token := range(tokens) {
		tokens_map = append(tokens_map, map[string]bool{take_off_quotation(token): inQuotes[i]})
	}

	return tokens_map
}

// splitの返り値から、文字列とboolを別々にする関数
func divide_split(tokens []map[string]bool) ([]string, []bool) {
	var divided_tokens []string
	var divided_tokens_bool []bool
	for _, token := range(tokens) {
		for key := range(token) {
			divided_tokens = append(divided_tokens, key)
		}
		for _, val := range(token) {
			divided_tokens_bool = append(divided_tokens_bool, val)
		}
	}

	return divided_tokens, divided_tokens_bool
}

func take_off_quotation(target string) string {
	if strings.HasPrefix(target, "'") && strings.HasSuffix(target, "'") {
		return strings.Trim(target, "'")
	} else if strings.HasPrefix(target, "\"") && strings.HasSuffix(target, "\"") {
		return strings.Trim(target, "\"")
	}

	return target
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

	return target
}

func Run(name string, value []string, value_in_quotes []bool, variables *map[string]string) {
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

		(*variables)[value[0]] = strings.Join(strings.Split(variables_replacer(variables, value[1], value_in_quotes[1], false), value[2]), " ")
	} else if name == "join" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in join function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		(*variables)[value[0]] = strings.Join(strings.Split(variables_replacer(variables, value[1], value_in_quotes[1], false), " "), value[2])
	} else if name == "addbeg" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addbeg function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(variables, value[1], value_in_quotes[1], true)
		// val = add_quotation(val)

		// もしvalue[0]の変数が空文字列なら、スペースを追加しない
		if (*variables)[value[0]] == "" {
			(*variables)[value[0]] = val
		} else {
			(*variables)[value[0]] = val + " " + (*variables)[value[0]]
		}
	} else if name == "addend" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in addend function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		val := variables_replacer(variables, value[1], value_in_quotes[1], true)

		// もしvalue[0]の変数が空文字列なら、スペースを追加しない
		if (*variables)[value[0]] == "" {
			(*variables)[value[0]] = val
		} else {
			(*variables)[value[0]] = (*variables)[value[0]] + " " + val
		}
	} else if name == "addnth" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in add function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		index, err := strconv.Atoi(variables_replacer(variables, value[1], value_in_quotes[1], false))
		if err != nil {
			fmt.Println("The error occurred in add function in array package. [2]")
			fmt.Println("The index is not integer.")
		}

		val := variables_replacer(variables, value[2], value_in_quotes[2], true)
		// val = add_quotation(val)

		(*variables)[value[0]] = (*variables)[value[0]][:index + 1] + " " + val + (*variables)[value[0]][index + 1:]
	} else if name == "replace" {
		// value[0] = target, value[1] = old index, value[2] = new value
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in replace function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		index, err := strconv.Atoi(variables_replacer(variables, value[1], value_in_quotes[1], false))
		if err != nil {
			fmt.Println("The error occurred in replace function in array package. [2]")
			fmt.Println("The index is not integer.")
		}

		// 一回配列に変換してから置換する
		slice := strings.Split((*variables)[value[0]], " ")

		val := variables_replacer(variables, value[2], value_in_quotes[2], true)
		// val = add_quotation(val)

		slice[index] = val
		(*variables)[value[0]] = strings.Join(slice, " ")
	} else if name == "delnth" {
		_, ok := (*variables)[value[0]]
		if !ok {
			fmt.Println("The error occurred in delnth function in array package. [1]")
			fmt.Println("The variable is not found.")
		}

		index, err := strconv.Atoi(variables_replacer(variables, value[1], value_in_quotes[1], false))
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
			if val == variables_replacer(variables, value[2], value_in_quotes[2], false) {
				index = i
				break
			}
		}

		(*variables)[value[0]] = strconv.Itoa(index)
	}
}

func Sharp(func_name string, args []string, args_in_quote []bool, variables *map[string]string) (string, bool) {
	if func_name == "search" {
		// args[0] = 対象の配列 args[1] = 検索する文字列 -> インデックスを返す(なければ-1)

		args[0] = variables_replacer(variables, args[0], args_in_quote[0], false)
		args[1] = variables_replacer(variables, args[1], args_in_quote[1], false)

		// 配列に変換
		slice, _ := divide_split(split(args[0]))

		// 配列から値を検索し、そのインデックスを返す なければ-1
		index := -1
		for i, val := range(slice) {
			if val == args[1] {
				index = i
				break
			}
		}

		return strconv.Itoa(index), false
	}

	return "", false
}