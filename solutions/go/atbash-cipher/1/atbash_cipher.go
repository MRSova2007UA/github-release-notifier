package atbashcipher

import "strings"
func Atbash(s string) string {
    s = strings.ToLower(s)
    runes := []rune(s)
    var res []rune
    count := 0
    for _, ss := range runes{
        if ss >= 'a' && ss <= 'z' {
            if count == 5 {
                res = append(res, ' ')
                count = 0
            }
            Arune := 'z' - ss + 'a'
            res = append(res, Arune)
            count++
        }
        if ss >= '0' && ss <= '9' {
            if count == 5 {
                res = append(res, ' ')
                count = 0
                }
            res = append(res, ss)
            count++
        }
    }
    return string(res)
}
