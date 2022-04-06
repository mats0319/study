package sort

func radixSort(intSlice []int) {
    moreSort := true
    splitDigit := 1
    for moreSort {
        moreSort = false

        buckets := make([][]int, 10)

        for i := range intSlice {
            if intSlice[i] / splitDigit * 10 != 0 { // if array need more sort
                moreSort = true
            }

            digit := intSlice[i] / splitDigit % 10

            if digit < 0 {
                digit *= -1
            }

            buckets[digit] = append(buckets[digit], intSlice[i])
        }

        for i := 0; i < len(intSlice); i++ {
            for bucketIndex := range buckets {
                for j := range buckets[bucketIndex] {
                    intSlice[i] = buckets[bucketIndex][j]
                    i++
                }
            }
        }

        splitDigit *= 10
    }

    res := make([]int, len(intSlice))
    left, right := 0, len(intSlice)-1
    for i := len(intSlice)-1; i >= 0; i-- {
        if intSlice[i] < 0 {
            res[left] = intSlice[i]
            left++
        } else {
            res[right] = intSlice[i]
            right--
        }
    }

    for i := range intSlice {
        intSlice[i] = res[i]
    }
}

func radixSortLSD_2(intSlice []int) {
    moreSort := true
    splitDigit := 1
    for moreSort {
        moreSort = false

        buckets := make([][]int, 20)

        for i := range intSlice {
            if intSlice[i] / splitDigit * 10 != 0 { // if array need more sort
                moreSort = true
            }

            digit := intSlice[i] / splitDigit % 10 + 9

            if intSlice[i] > 0 { // distinguish '-0' and '+0'
                digit++
            }

            buckets[digit] = append(buckets[digit], intSlice[i])
        }

        for i := 0; i < len(intSlice); i++ {
            for bucketIndex := range buckets {
                for j := range buckets[bucketIndex] {
                    intSlice[i] = buckets[bucketIndex][j]
                    i++
                }
            }
        }

        splitDigit *= 10
    }
}

func radixSortMSD_1(intSlice []int) {
    maxDigit := 0
    for _, v := range intSlice {
        digit := 0
        for splitDigit := 1; v > 0; splitDigit *= 10 {
            v /= splitDigit
            digit++
        }

        if maxDigit < digit {
            maxDigit = digit
        }
    }

    radixSortMSD_1_Recurse(intSlice, maxDigit)

    res := make([]int, len(intSlice))
    left, right := 0, len(intSlice)-1
    for i := len(intSlice)-1; i >= 0; i-- {
        if intSlice[i] < 0 {
            res[left] = intSlice[i]
            left++
        } else {
            res[right] = intSlice[i]
            right--
        }
    }

    for i := range intSlice {
        intSlice[i] = res[i]
    }
}

func radixSortMSD_1_Recurse(intSlice []int, maxDigit int) {
    if maxDigit < 1 {
        return
    }

    splitDigit := 1
    for i := 1; i < maxDigit; i++ {
        splitDigit *= 10
    }

    buckets := make([][]int, 10)

    for i := range intSlice {
        digit := intSlice[i] / splitDigit % 10

        if digit < 0 {
            digit *= -1
        }

        buckets[digit] = append(buckets[digit], intSlice[i])
    }

    for i := range buckets {
        radixSortMSD_1_Recurse(buckets[i], maxDigit-1)
    }

    for i := 0; i < len(intSlice); i++ {
        for bucketIndex := range buckets {
            for j := range buckets[bucketIndex] {
                intSlice[i] = buckets[bucketIndex][j]
                i++
            }
        }
    }
}

func radixSortMSD_2(intSlice []int) {
    maxDigit := 0
    for _, v := range intSlice {
        digit := 0
        for splitDigit := 1; v > 0; splitDigit *= 10 {
            v /= splitDigit
            digit++
        }

        if maxDigit < digit {
            maxDigit = digit
        }
    }

    radixSortMSD_2_Recurse(intSlice, maxDigit)
}

func radixSortMSD_2_Recurse(intSlice []int, maxDigit int) {
    if maxDigit < 1 {
        return
    }

    splitDigit := 1
    for i := 1; i < maxDigit; i++ {
        splitDigit *= 10
    }

    buckets := make([][]int, 20)

    for i := range intSlice {
        digit := intSlice[i] / splitDigit % 10 + 9

        if intSlice[i] > 0 {
            digit++
        }

        buckets[digit] = append(buckets[digit], intSlice[i])
    }

    for i := range buckets {
        radixSortMSD_2_Recurse(buckets[i], maxDigit-1)
    }

    for i := 0; i < len(intSlice); i++ {
        for bucketIndex := range buckets {
            for j := range buckets[bucketIndex] {
                intSlice[i] = buckets[bucketIndex][j]
                i++
            }
        }
    }
}
