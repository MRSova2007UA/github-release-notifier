package lasagnamaster

// TODO: define the 'PreparationTime()' function
func PreparationTime(layers []string, timePerLayer int) int {
    if timePerLayer == 0 {
        timePerLayer = 2
    }
    return len(layers) *timePerLayer
}
// TODO: define the 'Quantities()' function
func Quantities (layers []string) (int, float64) {
    var noodles int = 0
    var sauce float64 = 0
    for i := 0; i < len(layers); i++ {
        if layers[i] == "noodles" {
            noodles += 50
        }
        if layers[i] == "sauce" {
            sauce +=0.2
        }
    }
    return noodles, sauce
}
// TODO: define the 'AddSecretIngredient()' function
func AddSecretIngredient(friendsList, myList []string) {
	lastIng := friendsList[len(friendsList)-1]
	myList[len(myList)-1] = lastIng
}

// TODO: define the 'ScaleRecipe()' function
func ScaleRecipe(quantities []float64, numPort int) []float64 {
    scaled := make([]float64, len(quantities))
    for i := 0; i < len(quantities); i++ {
        scaled[i] = (quantities[i] / 2) * float64(numPort)
    }
    return scaled
}
// Your first steps could be to read through the tasks, and create
// these functions with their correct parameter lists and return types.
// The function body only needs to contain `panic("")`.
//
// This will make the tests compile, but they will fail.
// You can then implement the function logic one by one and see
// an increasing number of tests passing as you implement more
// functionality.
