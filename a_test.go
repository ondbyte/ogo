package cb_test

import (
	"fmt"
	"testing"

	"github.com/ondbyte/cb"
)

func TestIsPalindrome(t *testing.T) {
	tests := []string{
		"A man a plan a canal Panama",
		"race a car",
		"Was it a car or a cat I saw",
	}

	for _, test := range tests {
		fmt.Printf("%v: %v\n", test, cb.IsPalindrome(test))
	}
}
