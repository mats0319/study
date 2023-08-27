package utils

import (
    "math/rand"
)

// GenerateRandomIntSlice generate random int slice, you can set 'length' and 'max value' of slice
//   @param length: length of slice, min is 10
//   @param maxValue: max value of slice element, in fact, slice[i] is random in the area [-'max value', 'max value')
//   @param specialValues: special values in test, for each method, it may need some special case when test,
//                         more values than 'length' will be ignored
func GenerateRandomIntSlice(length int, maxValue int, specialValues ...int) []int {
    if length < 10 {
        length = 10
    }

    intSlice := make([]int, length) // length: big(10, 'length')

    i := 0
    for ; i < len(intSlice) && i < len(specialValues); i++ { // special values if given
        intSlice[i] = specialValues[i]
    }

    for ; i < len(intSlice); i++ { // random values
        intSlice[i] = 2*rand.Intn(maxValue+1) - maxValue // item value: [-'max value', 'max value']
    }

    return intSlice
}
