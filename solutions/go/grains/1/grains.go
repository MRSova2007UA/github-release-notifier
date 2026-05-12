package grains

import ("math"
        "errors")

func Square(number int) (uint64, error) {
	if number == 0 || number > 64 || number < 0 {
        return 0, errors.New("недійсні числа") 
    }
    res := math.Pow(2, float64(number - 1))
    return uint64(res), nil
}

func Total() uint64 {
	var total uint64 = 0
    var current uint64 = 1
    for i := 0; i <64; i++ {
        total += current
        current *= 2
    }
    return total
}
