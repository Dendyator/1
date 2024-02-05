package main

import (
	"fmt"

	"github.com/ozgio/strutil"
)

func main() {
	a := "Hello, OTUS!"
	a = strutil.Reverse(a)
	fmt.Print(a)
}
