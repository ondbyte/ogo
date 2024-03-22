package ogo

import (
	"regexp"
)

var pathParamsRX = regexp.MustCompile(`{(\w+)}`)

func PathParams(path string) (m map[string]int) {
	m = map[string]int{}
	matches := pathParamsRX.FindAllStringSubmatch(path, -1)
	for _, v := range matches {
		m[v[1]] = 1
	}
	return
}
