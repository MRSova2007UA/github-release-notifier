package pangram

func IsPangram(input string) bool {
	found := make(map[rune]bool)
    for _, char := range input {
        if char >= 'A' && char <= 'Z' {
            char += 'a' - 'A'
        }
        if char >= 'a' && char <= 'z' {
            found[char] = true
        }
    }
    return len(found) == 26
}