package main

import (
	"fmt"
	"regexp"
)

func main() {
	const configMapKeyFmt = `[-._a-zA-Z0-9]+`
	var configMapKeyRegexp = regexp.MustCompile("^" + configMapKeyFmt + "$")
	fmt.Println(configMapKeyRegexp.MatchString("."))
}
