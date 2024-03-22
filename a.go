package cb

import "fmt"

// "ABA" ->
// "xyzyx" ->
// "Xy zyx" -> true
// A man a plan a canal Panama

func IsPalindrome(s string) bool {
	i := 0
	j := len(s) - 1
	centre := j / 2
	for {
		if i >= j {
			break
		}
		if s[i] == ' ' {
			i++
			continue
		}
		if s[j] == ' ' {
			j--
			continue
		}
		if s[i] == s[j] || s[i]-32 == s[j] || s[i]+32 == s[j] {
			centre--
		} else {
			return false
		}
		i++
		j--
	}
	return true
}

func main() {
	fmt.Println(IsPalindrome("Xy 1zyx"))
}
