package largestseriesproduct

import "errors"

func LargestSeriesProduct(digits string, span int) (int64, error) {
    if span < 0 || span > len(digits) {
            return 0, errors.New("span must be smaller than string length and not negative")
        }
    if span == 0 {
        return 1, nil
    }
    var maxProduct int64
    for i := 0; i <= len(digits) - span; i++ {
        var currentProduct int64 = 1
        for j := 0; j < span; j++ {
            if digits[i+j] < '0' || digits[i+j] > '9' {
    			return 0, errors.New("digits input must only contain digits")
			}
            val := int64(digits[i+j] - '0')
            currentProduct *= val
        }
        if maxProduct < currentProduct {
            maxProduct = currentProduct
        }
    }
    return maxProduct, nil
}
