package rotationalcipher

func RotationalCipher(plain string, shiftKey int) string {
	if len(plain) == 0 {
        return ""
    }
    var res []rune
    shift := rune(shiftKey % 26)
    for _, r := range plain {
        if r >= 'a' && r <= 'z' {
            shifted := 'a' + (r-'a'+shift)%26
            res = append(res, shifted)
        }else if r >= 'A' && r <= 'Z'{
            shifted := 'A' + (r-'A'+shift)%26
            res = append(res, shifted)
        }else {
            res = append(res, r)
        }
    }
    return string(res)
}
