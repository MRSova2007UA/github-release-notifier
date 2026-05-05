package anagram

import ("strings"
        "slices")
func Sort (s string) string {
    runes := []rune(strings.ToLower(s))
    slices.Sort(runes)
    return string(runes)
}
func Detect(subject string, candidates []string) []string {
	var res []string
    for _, candidate := range candidates {
        if strings.ToLower(candidate) == strings.ToLower(subject) {
            continue
        }
        if Sort(candidate) == Sort(subject) {
            res = append(res, candidate)
        }
    }
    return res
}

