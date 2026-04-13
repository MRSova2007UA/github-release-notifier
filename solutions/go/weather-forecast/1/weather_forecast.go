// Package weather provides tools to format and return weather forecasts.
package weather

var (
    // CurrentCondition represents the current weather condition.
	CurrentCondition string

    // CurrentLocation represents the current city or location for the forecast.
	CurrentLocation  string
)
// Forecast updates the current location and condition, and returns a formatted weather string.
func Forecast(city, condition string) string {
    
	CurrentLocation, CurrentCondition = city, condition
	return CurrentLocation + " - current weather condition: " + CurrentCondition
}
