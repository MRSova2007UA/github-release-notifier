package booking

import "time"
import "fmt"

func Schedule(date string) time.Time {
	layout := "1/2/2006 15:04:05" 
	
	t, _ := time.Parse(layout, date)
	
	return t
}


// HasPassed повертає true, якщо дата вже минула.
func HasPassed(dateStr string) bool {
	now := time.Now()
    layout := "January 2, 2006 15:04:05"
    date, _ := time.Parse(layout, dateStr)
    if now.After(date) {
        return true
    }
    return false
}

// IsAfternoonAppointment returns whether a time is in the afternoon.
func IsAfternoonAppointment(date string) bool {
	layout := "Monday, January 2, 2006 15:04:05"
    t, _ := time.Parse(layout, date)
    return t.Hour() >= 12 && t.Hour() < 18
}

// Description returns a formatted string of the appointment time.
func Description(date string) string {
	layout := "1/2/2006 15:04:05"
    t, _ := time.Parse(layout, date)
    newLayout := "Monday, January 2, 2006, at 15:04."
    return fmt.Sprintf("You have an appointment on %s", t.Format(newLayout))
    
}

// AnniversaryDate returns a Time with this year's anniversary.
func AnniversaryDate() time.Time {
	thisYear := time.Now().Year()
    return time.Date(thisYear, time.September, 15, 0, 0, 0, 0, time.UTC)
}
