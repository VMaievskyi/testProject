package main

import (
	"fmt"
	"strings"
)

func ToString(delim string, array ...interface{}) string {
	return strings.Trim(strings.Replace(fmt.Sprint(array), " ", delim, -1), "[]")
}
