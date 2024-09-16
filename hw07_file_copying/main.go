package main

import (
	"flag"
	"fmt"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	from = "D:/Coding/GitHub/1/hw07_file_copying/testdata/input.txt"
	to = "D:/Coding/GitHub/1/hw07_file_copying/testdata/out.txt"
	flag.Parse()

	err := Copy(from, to, limit, offset)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("cmp out.txt testdata/out_offset%d_limit%d.txt", limit, offset)
}
