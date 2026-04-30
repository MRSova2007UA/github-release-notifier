package luhn

import ("strings"
       "unicode")

func Valid(id string) bool {
id = strings.ReplaceAll(id, " ", "")
	if len(id) <= 1 {
		return false
	}
	sum := 0
	shouldDouble := false	
	for i := len(id) - 1; i >= 0; i-- {
		r := rune(id[i])
		if !unicode.IsDigit(r) {
			return false
		}
		digit := int(r - '0')
		if shouldDouble {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
	}
	return sum%10 == 0
}
