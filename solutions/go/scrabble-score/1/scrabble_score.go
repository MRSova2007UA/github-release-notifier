package scrabblescore

import "strings"
func Score(word string) int {
	word = strings.ToUpper(word)
    counts := 0
    for _, n := range word {
    	switch n{
        	case 'A', 'E', 'I', 'O', 'U', 'L', 'N', 'R', 'S', 'T':
			counts += 1
			case 'D', 'G':
			counts += 2
			case 'B', 'C', 'M', 'P':
			counts += 3
			case 'F', 'H', 'V', 'W', 'Y':
			counts += 4
			case 'K':
			counts += 5
			case 'J', 'X':
			counts += 8
			case 'Q', 'Z':
			counts += 10
    	}
	}
    return counts
}