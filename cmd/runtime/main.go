package main

import (
	"fmt"
	"runtime"
)

type trace struct{}

func (t *trace) one(a int) {
	t.two()
}

func (t *trace) two() {
	t.three()
}

func (t *trace) three() {
	pcs := make([]uintptr, 10)
	n := runtime.Callers(0, pcs)
	fmt.Printf("n: %d, len(callers): %d\n", n, len(pcs))
	pcs = pcs[:n]
	frames := runtime.CallersFrames(pcs)
	fmt.Println()
	fmt.Println("====================== CALLERS ======================")
	for i := n; n > 0; i-- {
		frame, more := frames.Next()

		fmt.Printf("[%d] %s\n\t%s:%d\n", i, frame.Function, frame.File, frame.Line)

		if !more {
			break
		}
	}

	buf := make([]byte, 1024)
	n = runtime.Stack(buf, true)
	buf = buf[:n]
	fmt.Println()
	fmt.Println("======================  STACK  ======================")
	fmt.Println(string(buf))

	fmt.Println()
	panic("Dump Trace")
}

func main() {
	t := &trace{}
	t.one(4)
}
