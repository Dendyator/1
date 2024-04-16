package main

import (
	"fmt"
	"strconv"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func main() {

	g := func(_ string, f func(v interface{}) interface{}) func(in <-chan interface{}) (out <-chan interface{}) {

		return func(in <-chan interface{}) <-chan interface{} {
			out := make(chan interface{})
			go func() {
				defer close(out)
				for v := range in {
					out <- f(v)
				}
			}()
			return out
		}
	}

	//d := g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 })

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	generator := func(integers ...int) <-chan interface{} {
		intStream := make(chan interface{})
		go func() {
			defer close(intStream)
			for _, i := range integers {
				intStream <- i
			}
		}()
		return intStream
	}
	intStream := generator(2)

	street := func(z <-chan interface{}) <-chan interface{} {
		high := make(chan interface{})
		go func() {
			//defer close(high)
			high <- z
			var delta interface{}
			for _, v := range stages {
				delta = v(high)
				high <- delta
			}
		}()
		return high
	}

	s := street(intStream)

	for num := range s {
		fmt.Println(num)
	}

}
