package cars

// CalculateWorkingCarsPerHour calculates how many working cars are
// produced by the assembly line every hour.
func CalculateWorkingCarsPerHour(productionRate int, successRate float64) float64 {
	return float64 (productionRate) * (successRate / 100)
}

// CalculateWorkingCarsPerMinute calculates how many working cars are
// produced by the assembly line every minute.
func CalculateWorkingCarsPerMinute(productionRate int, successRate float64) int {
	carsPerMinute := float64(productionRate) / 60.0
	workingCars := carsPerMinute * (successRate / 100.0)
	return int(workingCars)
} 

// CalculateCost works out the cost of producing the given number of cars.
func CalculateCost(carsCount int) uint {
	carsGroup := carsCount / 10
    carsnoGroup := carsCount % 10
    return uint(carsGroup * 95000 + carsnoGroup * 10000)
}
