// This is a "stub" file.  It's a little start on your solution.
// It's not a complete solution though; you have to write some code.

// Package acronym should have a package comment that summarizes what it's about.
// https://golang.org/doc/effective_go.html#commentary
package acronym

import "strings"
// Abbreviate should have a comment documenting it.
func Abbreviate(s string) string {
	s = strings.ReplaceAll(s, "-", " ")
    s = strings.ReplaceAll(s, "_", " ")
    word := strings.Fields(s)
	var runeTLA []rune
    for i := 0; i < len(word); i++ {
        sRunes := []rune(word[i])
        if len(sRunes) > 0 {
            runeTLA = append(runeTLA, sRunes[0])
        }
    } 
	return strings.ToUpper(string(runeTLA))
}
