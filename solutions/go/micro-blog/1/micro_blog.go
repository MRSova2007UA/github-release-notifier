package microblog

func Truncate(phrase string) string {
	r := []rune(phrase)
    if len(r) > 5 {
        r = r[:5]
    }
    finalPhrase := string(r)
    return finalPhrase
}
