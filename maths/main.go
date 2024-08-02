package math

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

func variables_replacer(variables map[string]string, target string) string {
	val, ok := variables[target]
	if ok {
		return val
	}

	return target
}

func Sharp(func_name string, args []string, variables *map[string]string) string {
	if func_name == "pi" {
		return "3.1415926535"
	} else if func_name == "e" {
		return "2.7182818284"
	} else if func_name == "sin" {
		// args[0] ... 角度
		degree, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in sin function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Sin(degree * math.Pi / 180))
	} else if func_name == "cos" {
		// args[0] ... 角度
		degree, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in cos function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Cos(degree * math.Pi / 180))
	} else if func_name == "tan" {
		// args[0] ... 角度
		degree, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in tan function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Tan(degree * math.Pi / 180))
	} else if func_name == "asin" {
		// args[0] ... 角度
		degree, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in asin function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Asin(degree) * 180 / math.Pi)
	} else if func_name == "acos" {
		// args[0] ... 角度
		degree, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in acos function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Acos(degree) * 180 / math.Pi)
	} else if func_name == "atan" {
		// args[0] ... 角度
		degree, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in atan function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Atan(degree) * 180 / math.Pi)
	} else if func_name == "sqrt" {
		// args[0] ... 数値
		num, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in sqrt function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Sqrt(num))
	} else if func_name == "log" {
		// args[0] ... 数値
		num, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in log function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Log(num))
	} else if func_name == "log10" {
		// args[0] ... 数値
		num, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in log10 function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Log10(num))
	} else if func_name == "pow" {
		// args[0] ... 底 args[1] ... 指数
		base, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in pow function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		exponent, err := strconv.ParseFloat(variables_replacer(*variables, args[1]), 64)
		if err != nil {
			fmt.Println("The error occurred in pow function in math package. [2]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Pow(base, exponent))
	} else if func_name == "abs" {
		// args[0] ... 数値
		num, err := strconv.ParseFloat(variables_replacer(*variables, args[0]), 64)
		if err != nil {
			fmt.Println("The error occurred in abs function in math package. [1]")
			fmt.Println(err)
			return ""
		}

		return fmt.Sprintf("%f", math.Abs(num))
	} else if func_name == "median" {
		// args[0] ... 数値の配列(スペース区切り)
		nums := strings.Split(variables_replacer(*variables, args[0]), " ")
		for i, num := range nums {
			nums[i] = variables_replacer(*variables, num)
		}

		nums_list := []float64{}
		for _, num := range nums {
			num, err := strconv.ParseFloat(num, 64)
			if err != nil {
				fmt.Println("The error occurred in median function in math package. [1]")
				fmt.Println(err)
				return ""
			}

			nums_list = append(nums_list, num)
		}

		// ソート
		for i := 0; i < len(nums_list); i++ {
			for j := i + 1; j < len(nums_list); j++ {
				if nums_list[i] > nums_list[j] {
					tmp := nums_list[i]
					nums_list[i] = nums_list[j]
					nums_list[j] = tmp
				}
			}
		}

		// 中央値を求める
		if len(nums_list) % 2 == 0 {
			// 偶数
			return fmt.Sprintf("%f", (nums_list[len(nums_list) / 2 - 1] + nums_list[len(nums_list) / 2]) / 2)
		} else {
			// 奇数
			return fmt.Sprintf("%f", nums_list[len(nums_list) / 2])
		}
	} else if func_name == "average" {
		// args[0] ... 数値の配列(スペース区切り)
		nums := strings.Split(variables_replacer(*variables, args[0]), " ")
		for i, num := range nums {
			nums[i] = variables_replacer(*variables, num)
		}

		nums_list := []float64{}
		for _, num := range nums {
			num, err := strconv.ParseFloat(num, 64)
			if err != nil {
				fmt.Println("The error occurred in average function in math package. [1]")
				fmt.Println(err)
				return ""
			}

			nums_list = append(nums_list, num)
		}

		// 合計を求める
		sum := 0.0
		for _, num := range nums_list {
			sum += num
		}

		return fmt.Sprintf("%f", sum / float64(len(nums_list)))
	} else if func_name == "mode" {
		// args[0] ... 数値の配列(スペース区切り)
		nums := strings.Split(variables_replacer(*variables, args[0]), " ")
		for i, num := range nums {
			nums[i] = variables_replacer(*variables, num)
		}

		nums_list := []float64{}
		for _, num := range nums {
			num, err := strconv.ParseFloat(num, 64)
			if err != nil {
				fmt.Println("The error occurred in mode function in math package. [1]")
				fmt.Println(err)
				return ""
			}

			nums_list = append(nums_list, num)
		}

		// 出現回数を数える
		counts := map[float64]int{}
		for _, num := range nums_list {
			counts[num]++
		}

		// 最大値を求める
		max := 0
		for _, count := range counts {
			if count > max {
				max = count
			}
		}

		// 最大値を持つものを列挙
		modes := []float64{}
		for num, count := range counts {
			if count == max {
				modes = append(modes, num)
			}
		}

		// すべての要素が1回ずつ出現している場合は、一番最初の要素を返す
		if len(modes) == len(nums_list) {
			return "none"
		}

		return fmt.Sprintf("%v", modes)
	}

	return ""
}
