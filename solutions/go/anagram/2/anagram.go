package anagram

import ("strings"
        "slices")
func Sort (s string) string {
    runes := []rune(s)
    slices.Sort(runes)
    return string(runes)
}
func Detect(subject string, candidates []string) []string {
	var res []string
    subjectToL := strings.ToLower(subject)
    subjectSorted := Sort(subjectToL)
    for _, candidate := range candidates {
        candToL := strings.ToLower(candidate)
        if candToL == subjectToL {
            continue
        }
        if Sort(candToL) == subjectSorted {
            res = append(res, candidate)
        }
    }
    return res
}

