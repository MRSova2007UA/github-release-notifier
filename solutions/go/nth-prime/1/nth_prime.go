package nthprime

import "errors"
// Nth returns the nth prime number. An error must be returned if the nth prime number can't be calculated ('n' is equal or less than zero)
func isPrime(num int) bool{
	if num < 2 {
        return false
    }
    if num == 2 {
        return true
    }
    if num % 2 == 0 {
        return false
    }
    for i := 3; i*i <= num; i += 2 {
        if num % i == 0 {
            return false
        }
    }
	return true
}
func Nth(n int) (int, error) {
	if n < 1 {
        return 0, errors.New("хуня малята")
    }
    count := 0
    candidate := 1
    for count < n {
        candidate++
         if isPrime(candidate) {
             count++
         }
    }
    return candidate, nil
}
