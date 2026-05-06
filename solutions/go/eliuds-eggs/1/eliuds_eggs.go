package eliudseggs

func EggCount(displayValue int) int {
	count := 0
    for displayValue > 0 {
    	if displayValue & 1 == 1 {
        	count++
    	}
        displayValue >>= 1
	}
    return count
}