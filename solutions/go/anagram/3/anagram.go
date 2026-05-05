package anagram

import ("strings"
        "slices")
func sort (s string) string {
    runes := []rune(s)
    slices.Sort(runes)
    return string(runes)
}
func Detect(subject string, candidates []string) []string {
	var res []string
    subjectToL := strings.ToLower(subject)
    subjectSorted := sort(subjectToL)
    for _, candidate := range candidates {
        candToL := strings.ToLower(candidate)
        if candToL == subjectToL {
            continue
        }
        if sort(candToL) == subjectSorted {
            res = append(res, candidate)
        }
    }
    return res
}

