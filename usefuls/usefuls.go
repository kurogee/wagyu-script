package usefuls

func All(n []bool) bool {
    for i := 0; i < len(n); i++ {
        if !n[i] {
            return false
        }
    }
    
    return true
}

func Is_prime(n int) bool {
	if n == 1 {
		return false
	}

    var result bool = All([]bool{
        n % 2 != 0 || n == 2, n % 3 != 0 || n == 3, n % 5 != 0 || n == 5, n % 7 != 0 || n == 7,})

    return result
}

func Contains(slice []string, search string) bool {
    for _, val := range(slice) {
        if search == val {
            return true
        }
    }
    
    return false
}

func Delete_slice(slice []string, target string) []string {
    result := []string{}
    for _, val := range(slice) {
        if val != target {
            result = append(result, val)
        }
    }

    return result
}