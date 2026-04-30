package isogram

import "unicode"

func IsIsogram(word string) bool {
    seen := make(map[rune]struct{})
	for _, r := range word {
        if r == '-' || unicode.IsSpace(r) {
            continue
        }
        r = unicode.ToLower(r)
        if _, exist := seen[r]; exist {
            return false
        }
        seen[r] = struct{}{}
    }
    return true
}
