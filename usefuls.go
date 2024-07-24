package usefuls

func All(n []bool) bool {
    for _, val := range(n) {
        if !val {
            return false
        }
    }
        
    return true
}

func Is_prime(n int) bool {
    var checkers = []bool{}

	if n == 1 {
		return false
	}
        
    checkers = append(checkers, !(n % 2 == 0 && n != 2))
    checkers = append(checkers, !(n % 3 == 0 && n != 3))
    checkers = append(checkers, !(n % 5 == 0 && n != 5))
    checkers = append(checkers, !(n % 7 == 0 && n != 7))

    return All(checkers)
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