package utils

import (
	"fmt"
	"testing"
)

func TestGenerateRandomSlice(t *testing.T) {
	for range 3 {
		fmt.Println(GenerateRandomSlice[int](10, 20))
		fmt.Println(GenerateRandomSlice[float64](10, 20.0))
	}

	// [-13 4 7 -4 -15 12 -9 4 -8 0]
	// [4.662 14.249 -19.734 14.437 -8.219 -11.027 -0.751 -2.469 7.318 17.657]
	// [18 15 14 16 -9 -12 16 -18 19 18]
	// [4.937 6.893 -13.351 19.661 -6.618 -2.81 -14.096 5.774 -10.319 -10.349]
	// [-10 9 -14 1 -2 2 -1 8 5 -16]
	// [15.717 -11.257 17.042 -14.607 -16.935 -7.369 -2.736 9.796 -11.652 -0.121]
}

func TestGenerateRandomBytes(t *testing.T) {
	for range 3 {
		fmt.Println(string(GenerateRandomBytes(32)))
	}

	//ZZSABFOJNES3TMKGVOKZMGXFNS7L6YXW
	//PMWLBXHFEBBVHBK7K5EUMCDMJ7NICQRI
	//4WGAIU4FKYLQOFB337ZMXFASQ3DPKLUK
}
