package blackjack

	// ParseCard returns the integer value of a card following blackjack ruleset.
func ParseCard(card string) int {
	switch card {
	case "ace":
		return 11
	case "two":
		return 2
	case "three":
		return 3
	case "four":
		return 4
	case "five":
		return 5
	case "six":
		return 6
	case "seven":
		return 7
	case "eight":
		return 8
	case "nine":
		return 9
	case "ten", "jack", "queen", "king":
		return 10
	default:
		return 0
	}
}

// FirstTurn returns the decision for the first turn, given two cards of the
// player and one card of the dealer.
func FirstTurn(card1, card2, dealerCard string) string {
	val1 := ParseCard(card1)
	val2 := ParseCard(card2)
	dealerVal := ParseCard(dealerCard)
	
	sum := val1 + val2

	switch {
	case card1 == "ace" && card2 == "ace":
		return "P" // Split (Розділити, якщо два тузи)
		
	case sum == 21:
		if dealerVal < 10 {
			return "W" // Win (Перемога, якщо у дилера немає карти на 10+ очок)
		}
		return "S" // Stand (Інакше просто чекаємо)
		
	case sum >= 17 && sum <= 20:
		return "S" // Stand (Завжди чекаємо)
		
	case sum >= 12 && sum <= 16:
		if dealerVal >= 7 {
			return "H" // Hit (Беремо карту, якщо у дилера 7 або більше)
		}
		return "S" // Stand (Інакше чекаємо)
		
	default:
		return "H" // Hit (Якщо 11 очок або менше - завжди беремо карту)
	}
}

