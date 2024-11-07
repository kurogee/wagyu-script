package system_split

import (
	"fmt"
	"regexp"
	"strings"
)

type TokenFormula struct {
	Pattern string
	Regex *regexp.Regexp
}

var replace_chars map[string]string = map[string]string{
	"\\;" : "__SEMICOLON__",
	"\\\"" : "__DOUBLE_QUOTATION__",
	"\\'" : "__SINGLE_QUOTATION__",
	"\\`" : "__BACK_QUOTATION__",
	"\\\\\\" : "__BACKSLASH__",
    "\\(" : "__LEFT_PAREN__",
    "\\)" : "__RIGHT_PAREN__",
};

func Take_off_quotation(target string) string {
	if strings.HasPrefix(target, "'") && strings.HasSuffix(target, "'") {
		// 一番外側のシングルクォーテーションを取り除く
		re := regexp.MustCompile(`^'|'$`)
		return re.ReplaceAllString(target, "")
	} else if strings.HasPrefix(target, "\"") && strings.HasSuffix(target, "\"") {
		re := regexp.MustCompile(`^"|"$`)
		return re.ReplaceAllString(target, "")
	} else if strings.HasPrefix(target, "`") && strings.HasSuffix(target, "`") {
		re := regexp.MustCompile(`^`+"`"+`|`+"`"+`$`)
		return re.ReplaceAllString(target, "")
	}

	return target
}

func Split(input string) []map[string]bool {
    var tokens []string
    var buffer strings.Builder
    var quoteChar rune
    var parenStack []rune
    var bracketStack []rune

    re := regexp.MustCompile(`([\(\{\[\)\}\]])\s+([\(\{\[\)\}\]])`)
    input = re.ReplaceAllString(input, "$1 $2")

    for key, value := range replace_chars {
        input = strings.ReplaceAll(input, key, value)
    }

    inQuote := false
    inQuotes := []bool{}

    inParen := false
    inBracket := false
    inFunctionCall := false

    parenRegex := regexp.MustCompile(`^[\{\(\)\}]$`)
    bracketRegex := regexp.MustCompile(`^[\[\]]$`)
    functionCallRegex := regexp.MustCompile(`^#[^\(]*\(`)

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
            if len(parenStack) > 0 && char == parenStack[len(parenStack)-1] {
                parenStack = parenStack[:len(parenStack)-1]
                if len(parenStack) == 0 {
                    inParen = false

                    tokens = append(tokens, buffer.String())
                    buffer.Reset()

                    inQuotes = append(inQuotes, false)
                }
            } else if parenRegex.MatchString(string(char)) {
                parenStack = append(parenStack, matchingParen(char))
            }

        case inBracket:
            if len(bracketStack) > 0 && char == bracketStack[len(bracketStack)-1] {
                buffer.WriteRune(' ')
                bracketStack = bracketStack[:len(bracketStack)-1]
                inBracket = false
                tokens = append(tokens, buffer.String())
                buffer.Reset()
                inQuotes = append(inQuotes, false)
            } else if bracketRegex.MatchString(string(char)) {
                buffer.WriteRune(char)
                bracketStack = append(bracketStack, matchingParen(char))
            } else {
                buffer.WriteRune(char)
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
                if buffer.Len() > 0 && !inParen && !inBracket {
                    tokens = append(tokens, buffer.String())
                    buffer.Reset()
                    inQuotes = append(inQuotes, false)
                }

            case '\'', '"', '`':
                inQuote = true
                quoteChar = char
                buffer.WriteRune(char)

            case '(':
                inParen = true
                parenStack = append(parenStack, matchingParen(char))
                buffer.WriteRune(char)
            
            case '{':
                inParen = true
                parenStack = append(parenStack, matchingParen(char))
                buffer.WriteRune(char)
            
            case '[':
                inBracket = true
                bracketStack = append(bracketStack, matchingParen(char))

            case ')', '}', ']':
                // 追加: 閉じ括弧が直接現れた場合の処理
                if inParen || inBracket {
                    if inParen {
                        buffer.WriteRune(char)
                        if len(parenStack) > 0 && char == parenStack[len(parenStack)-1] {
                            parenStack = parenStack[:len(parenStack)-1]
                            if len(parenStack) == 0 {
                                inParen = false
                                tokens = append(tokens, buffer.String())
                                buffer.Reset()
                                inQuotes = append(inQuotes, false)
                            }
                        }
                    } else if inBracket {
                        buffer.WriteRune(char)
                        if len(bracketStack) > 0 && char == bracketStack[len(bracketStack)-1] {
                            bracketStack = bracketStack[:len(bracketStack)-1]
                            if len(bracketStack) == 0 {
                                inBracket = false
                                tokens = append(tokens, buffer.String())
                                buffer.Reset()
                                inQuotes = append(inQuotes, false)
                            }
                        }
                    }
                } else {
                    buffer.WriteRune(char)
                }

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

    for key, value := range replace_chars {
        for i, token := range tokens {
            tokens[i] = strings.ReplaceAll(token, value, key)
        }
    }

    var tokens_map []map[string]bool
    
    for i, token := range(tokens) {
        tokens_map = append(tokens_map, map[string]bool{token: inQuotes[i]})
    }

    fmt.Print("")
    return tokens_map
}

// splitの返り値から、文字列とboolを別々にする関数
func Divide_split(tokens []map[string]bool) ([]string, []bool) {
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
    case '[':
        return ']'
    case ')':
        return '('
    case '}':
        return '{'
    case ']':
        return '['
    }
    return char
}