package etl

import "strings"
func Transform(in map[int][]string) map[string]int {
	res := make(map[string]int)
    for score, letters := range in {
        for _, letter := range letters {
            lowLetter := strings.ToLower(letter)
            res[lowLetter] = score
        }    
    }
    return res
}
