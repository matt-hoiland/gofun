package main

import (
	"fmt"
	"math/rand"
)

func main() {
	for j := 0; j < 32; j++ {
		var row [80]rune
		for i := 0; i < len(row); i++ {
			switch rand.Intn(2) {
			case 0:
				row[i] = '/'
				break
			case 1:
				row[i] = '\\'
				break
			}
		}
		fmt.Println(string(row[:]))
	}
}
