package utils

import (
	"fmt"
	"math/rand/v2"
	"strconv"
)

type Computable interface {
	int | int8 | int16 | int32 | int64 | float32 | float64
}

// GenerateRandomSlice generate random slice, you can set 'length' and 'max value' of slice
//
//	@param length: length of slice, min is 10
//	@param maxValue: max value of slice element, in fact, slice[i] is random in the area [-'max value', 'max value')
//	@param specialValues: special values in test, for each method, it may need some special case when test,
//	                      more values than 'length' will be ignored
func GenerateRandomSlice[T Computable](length int, maxValue T, specialValues ...T) []T {
	length = max(length, 10)

	slice := make([]T, length)

	i := 0
	for ; i < len(slice) && i < len(specialValues); i++ { // special values if given
		slice[i] = specialValues[i]
	}

	for ; i < len(slice); i++ {
		// random values
		// 2a*f(x) - a : [-a, a)
		//  - a: maxValue
		//  - f(x): Float64() : [0.0, 1.0)
		randomFloat := 2*float64(maxValue)*rand.Float64() - float64(maxValue)
		formatFloat, _ := strconv.ParseFloat(fmt.Sprintf("%.3f", randomFloat), 64)
		slice[i] = T(formatFloat)
	}

	return slice
}
