package main

import "fmt"

func main() {
	funcs := []struct {
		name      string
		generator func() []func() int
	}{
		{name: "close on loop index", generator: loopBasic},
		{name: "close on loop index value", generator: loopCapture},
		{name: "close on range iterator", generator: rangeBasic},
		{name: "close on range iterator value", generator: rangeCapture},
	}
	for _, s := range funcs {
		fmt.Printf("%s:\n", s.name)
		closures := s.generator()
		for _, closure := range closures {
			fmt.Printf("  %d\n", closure())
		}
	}
}

func loopBasic() []func() int {
	var ret []func() int
	for i := 1; i <= 3; i++ {
		ret = append(ret, func() int { return i })
	}
	return ret
}

func loopCapture() []func() int {
	var ret []func() int
	for i := 1; i <= 3; i++ {
		i := i
		ret = append(ret, func() int { return i })
	}
	return ret
}

func rangeBasic() []func() int {
	var ret []func() int
	for _, v := range []int{1, 2, 3} {
		ret = append(ret, func() int { return v })
	}
	return ret
}

func rangeCapture() []func() int {
	var ret []func() int
	for _, v := range []int{1, 2, 3} {
		v := v
		ret = append(ret, func() int { return v })
	}
	return ret
}
