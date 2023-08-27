package utils

import (
	"fmt"
	"testing"
)

func TestGenerateRandomIntSlice(t *testing.T) {
	for range [3]struct{}{} {
		fmt.Println(GenerateRandomIntSlice(10, 10))
	}

	// tips: The math/rand package now automatically seeds the global random number generator
	// (used by top-level functions like Float64 and Int) with a random value,
	// and the top-level Seed function has been deprecated.
	// reference: https://tip.golang.org/doc/go1.20#minor_library_changes

	//[-6 0 10 -8 8 -10 -10 -8 -2 10]
	//[-10 -8 -8 -8 -4 6 6 -2 -4 6]
	//[-6 -8 0 6 -10 2 -6 -10 4 6]

	//[-2 4 10 6 -8 6 8 0 4 4]
	//[-4 -8 10 4 -6 6 -10 0 -10 0]
	//[-2 -6 -4 -6 -6 -6 -2 6 4 -4]
}
