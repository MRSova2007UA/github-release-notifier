package techpalace

import ("fmt"
        "strings")

// WelcomeMessage returns a welcome message for the customer.
func WelcomeMessage(customer string) string {
    UpLetters := strings.ToUpper(customer)
	return fmt.Sprintf("Welcome to the Tech Palace, %s", UpLetters)	
}

// AddBorder adds a border to a welcome message.
func AddBorder(welcomeMsg string, numStarsPerLine int) string {
	stars := strings.Repeat("*", numStarsPerLine)
    return fmt.Sprintf("%s\n%s\n%s", stars, welcomeMsg, stars)
}

// CleanupMessage cleans up an old marketing message.
func CleanupMessage(oldMsg string) string {
	clearStars := strings.ReplaceAll(oldMsg, "*", "")
    cleanSpace := strings.TrimSpace(clearStars)
    return cleanSpace
}
