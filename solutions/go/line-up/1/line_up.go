package lineup

import "fmt"

func Format(name string, number int) string {
	lastTwo := number % 100
    if lastTwo >= 11 && lastTwo <= 13 {
        return fmt.Sprintf("%s, you are the %dth customer we serve today. Thank you!", name, number)
    }
    lastDigit := number % 10
    switch lastDigit {
        case 1:
        return fmt.Sprintf("%s, you are the %dst customer we serve today. Thank you!", name, number)
        case 2:
        return fmt.Sprintf("%s, you are the %dnd customer we serve today. Thank you!", name, number)
        case 3:
        return fmt.Sprintf("%s, you are the %drd customer we serve today. Thank you!", name, number)
        default:
        return fmt.Sprintf("%s, you are the %dth customer we serve today. Thank you!", name, number)
    }
}
