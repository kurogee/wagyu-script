package file

import (
	"fmt"
	"strings"
	"strconv"
	"regexp"
	"os"
)

var replace_chars map[string]string = map[string]string{
	"\\;" : "__SEMICOLON__",
	"\\\"" : "__DOUBLE_QUOTATION__",
	"\\'" : "__SINGLE_QUOTATION__",
	"\\`" : "__BACK_QUOTATION__",
	"\\\\\\" : "__BACKSLASH__",
};

func take_off_quotation(target string) string {
	if strings.HasPrefix(target, "'") && strings.HasSuffix(target, "'") {
		return strings.Trim(target, "'")
	} else if strings.HasPrefix(target, "\"") && strings.HasSuffix(target, "\"") {
		return strings.Trim(target, "\"")
	}

	return target
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
	
	// tokensのエスケープシーケンスを元に戻す
	for key, value := range(replace_chars) {
		for i, token := range(tokens) {
			tokens[i] = strings.ReplaceAll(token, value, key)
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

func replaceSymbols(input string) string {
	input = strings.ReplaceAll(input, "\n", "\\n")
	input = strings.ReplaceAll(input, "\t", "\\t")
	input = strings.ReplaceAll(input, "\\", "\\\\")
	input = strings.ReplaceAll(input, "\"", "\\\"")
	input = strings.ReplaceAll(input, "'", "\\'")

	return input
}

func giveSymbols(input string) string {
	input = strings.ReplaceAll(input, "\\n", "\n")
	input = strings.ReplaceAll(input, "\\t", "\t")
	input = strings.ReplaceAll(input, "\\\\", "\\")
	input = strings.ReplaceAll(input, "\\\"", "\"")
	input = strings.ReplaceAll(input, "\\'", "'")
	
	return input
}

func Run(name string, value []string, value_in_quotes []bool, variables *map[string]string) {
	if name == "read" {
		// value[0] ... ファイルのパス value[1] ... 値を入れる変数名
		data, err := os.ReadFile(variables_replacer(variables, value[0], value_in_quotes[0], false))
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

		(*variables)[value[1]] = string(giveSymbols(string(data)))
	} else if name == "write" {
		// value[0] ... 書き込むファイルのパス value[1] ... 書き込む値（変数もあり）
		write_data := []byte(giveSymbols(variables_replacer(variables, value[1], value_in_quotes[1], false)))
		err := os.WriteFile(variables_replacer(variables, value[0], value_in_quotes[0], false), write_data, os.ModePerm)
		if err != nil {
			fmt.Println("The error occurred in write function in file package. [1]")
			fmt.Println(err)
			return
		}
	} else if name == "addend" {
		// value[0] ... 追記するファイルのパス value[1] ... 追記する値（変数もあり）
		addend_data := giveSymbols(variables_replacer(variables, value[1], value_in_quotes[1], false))
		file_data, err := os.ReadFile(variables_replacer(variables, value[0], value_in_quotes[0], false))
		if err != nil {
			fmt.Println("The error occurred in addend function in file package. [1]")
			fmt.Println(err)
			return
		}

		all_data := []byte(string(file_data) + addend_data)

		err2 := os.WriteFile(variables_replacer(variables, value[0], value_in_quotes[0], false), all_data, os.ModePerm)
		if err2 != nil {
			fmt.Println("The error occurred in addend function in file package. [2]")
			fmt.Println(err)
			return
		}
	} else if name == "remove" {
		// value[0] ... 削除するファイルのパス
		err := os.Remove(variables_replacer(variables, value[0], value_in_quotes[0], false))
		if err != nil {
			fmt.Println("The error occurred in remove function in file package. [1]")
			fmt.Println(err)
			return
		}
	} else if name == "rename" {
		// value[0] ... リネームするファイルのパス value[1] ... 新しいファイル名
		err := os.Rename(variables_replacer(variables, value[0], value_in_quotes[0], false), variables_replacer(variables, value[1], value_in_quotes[1], false))
		if err != nil {
			fmt.Println("The error occurred in rename function in file package. [1]")
			fmt.Println(err)
			return
		}
	} else if name == "readline" {
		// value[0] ... ファイルのパス value[1] ... 値を入れる変数名 (value[2] ... 読み込む行数 デフォルトは全部)
		data, err := os.ReadFile(variables_replacer(variables, value[0], value_in_quotes[0], false))
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

			(*variables)[value[1]] = strings.Join(lines[:], " ")
		} else {
			how_many, err := strconv.Atoi(variables_replacer(variables, value[2], value_in_quotes[2], false))
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

func Sharp(func_name string, args []string, args_in_quote []bool, variables *map[string]string) (string, bool) {
	if func_name == "read" {
		// args[0] ... ファイルのパス -> ファイルの中身を返す
		data, err := os.ReadFile(variables_replacer(variables, args[0], args_in_quote[0], false))
		
		if err != nil {
			fmt.Println("The error occurred in read sharp function in file package. [1]")
			fmt.Println(err)
			return "", false
		}

		return string(data), true
	}

	return "", false
}