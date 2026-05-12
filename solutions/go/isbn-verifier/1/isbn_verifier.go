package isbnverifier

import ("strings"
        "strconv")

func IsValidISBN(isbn string) bool {
	norm := strings.ReplaceAll(isbn, "-", "")
    if len(norm) != 10 {
        return false
    }
    count := 0
    for i, char := range norm {
        var num int
        if i == 9 && char == 'X' {
            num = 10
        }else {
			n, err := strconv.Atoi(string(char))
			if err != nil {
				return false
			}
			num = n
		}
        count += num * (10 - i)
    }
    return count%11 == 0
}
