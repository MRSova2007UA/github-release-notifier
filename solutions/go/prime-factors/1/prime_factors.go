package primefactors

func Factors(n int64) []int64 {
	var libr []int64
    num := n
    for num % 2 == 0 {
        libr = append(libr, 2)
        num /= 2
    }
    for i := int64(3); i*i <= num; i += 2 {
        for num % i == 0 {
            libr = append(libr, i)
            num /= i
        }
    }
    if num > 1 {
        libr = append(libr, num)
    }
	return libr
}

 