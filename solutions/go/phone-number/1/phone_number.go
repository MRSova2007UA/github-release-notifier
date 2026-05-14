package phonenumber

import ("strings"
        "errors")

func Number(phoneNumber string) (string, error) {
    digits := strings.Map(func(r rune) rune {
        if r >= '0' && r <= '9' {
            return r
        }
        return -1
    }, phoneNumber)
    if len(digits) == 11 {
        if digits[0] != '1' {
            return "", errors.New("11 digits must start with 1")
        }
        digits = digits[1:]
    }
    if len(digits) != 10 {
        return "", errors.New("must be 10 digits")
    }
    if digits[0] < '2' || digits[3] < '2' {
        return "", errors.New("area or exchange code cannot start with 0 or 1")
    }
    return digits, nil
}


func AreaCode(phoneNumber string) (string, error) {
	s, err := Number(phoneNumber)
    if err != nil {
        return "", err
    }
    return s[:3], nil
}

func Format(phoneNumber string) (string, error) {
    digits, err := Number(phoneNumber)
    if err != nil {
        return "", err
    }
	return "(" + digits[:3] + ") " + digits[3:6] + "-" + digits[6:], nil
}
