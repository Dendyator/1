package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

var fromPath = "D:/Coding/Projects/NinjaProject/Uroki/Les1/input.txt"
var toPath = "D:/Coding/Projects/NinjaProject/Uroki/Les1/output.txt"

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()
	err := Copy(fromPath, toPath, 3000, 1000)
	if err != nil {
		fmt.Println(err)
	}
}
