package romannumerals

import "errors"
func ToRomanNumeral(input int) (string, error) {
	conv := []struct {
        val int
        sym string
    }{	
        {1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
		{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
		{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
    }
    if input <= 0 || input > 3999 {
        return "", errors.New("input must be between 1 and 3999")
    }
    var result string
    for _, c := range conv {
        for input >= c.val {
            result += c.sym
            input -= c.val
        }
    }
    return result, nil
}
