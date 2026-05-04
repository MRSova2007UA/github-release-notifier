package armstrongnumbers

import ("strconv"
    	"math")
func IsNumber(n int) bool {
	strN := strconv.Itoa(n)
    numDigits := len(strN)
	sum := 0
    temp := n
    for temp > 0{
        digit := temp % 10
        sum += int(math.Pow(float64(digit), float64(numDigits)))
        temp /= 10
    }
    return sum == n
}
