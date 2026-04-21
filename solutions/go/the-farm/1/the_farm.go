package thefarm

import ("errors"
        "fmt")

// TODO: define the 'DivideFood' function
func DivideFood (f FodderCalculator, num int) (float64, error) {
    amount, err := f.FodderAmount(num)
    if err != nil {
        return 0, err
    }
    factor, err := f.FatteningFactor()
    if err != nil {
        return 0, err
    }
    result := (amount * factor) / float64(num)
    return result, nil
}
// TODO: define the 'ValidateInputAndDivideFood' function
func ValidateInputAndDivideFood(f FodderCalculator, num int) (float64, error) {
    if num <= 0 {
        var err = errors.New("invalid number of cows")
        return 0, err
    }
    res, err := DivideFood(f, num)
    return res, err
}
// TODO: define the 'ValidateNumberOfCows' function
type InvalidCowsError struct {
    cows int
    message string
}
func (e *InvalidCowsError) Error() string {
    return fmt.Sprintf("%d cows are invalid: %s", e.cows, e.message)
}

func ValidateNumberOfCows(num int) error {
    if num < 0 {
        return &InvalidCowsError{
            cows: num,
            message: "there are no negative cows",        
        }
    }
    if num == 0 {
        return &InvalidCowsError{
            cows: num,
            message: "no cows don't need food",
        }
    }
    return nil
}
// Your first steps could be to read through the tasks, and create
// these functions with their correct parameter lists and return types.
// The function body only needs to contain `panic("")`.
//
// This will make the tests compile, but they will fail.
// You can then implement the function logic one by one and see
// an increasing number of tests passing as you implement more
// functionality.
