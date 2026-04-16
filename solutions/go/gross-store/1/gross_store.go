package gross

// Units stores the Gross Store unit measurements.
func Units() map[string]int {
    return map[string]int {
		"quarter_of_a_dozen": 3,
    	"half_of_a_dozen": 6,
    	"dozen": 12,
    	"small_gross": 120,
    	"gross": 144,
    	"great_gross": 1728,
	}
}
// NewBill creates a new bill.
func NewBill() map[string]int {
	return make(map[string]int)
}

// AddItem adds an item to customer bill.
func AddItem(bill, units map[string]int, item, unit string) bool {
    val, ok := units[unit]
    if !ok {
        return false
    }else {
        bill[item] += val
        return true
    }
}

// RemoveItem removes an item from customer bill.
func RemoveItem(bill, units map[string]int, item, unit string) bool {
	quantity, ok1 := bill[item]
    val, ok2 := units[unit]
    if !ok1 {
        return false
    }
    if !ok2 {
        return false
    }
    newQuantity := quantity - val
    if newQuantity < 0 {
    	return false
    }
    if newQuantity == 0 {
        delete(bill, item)
        return true
    }
    bill[item] = newQuantity
    return true
}
// GetItem returns the quantity of an item that the customer has in his/her bill.
func GetItem(bill map[string]int, item string) (int, bool) {
	val, ok := bill[item]
    if !ok {
        return val, ok
    }else {
        return val, ok
    }
}
