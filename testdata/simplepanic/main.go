package main

import (
	"flag"
	"fmt"
	"simplepanic/concurrent"
	"simplepanic/recovering"
	"simplepanic/simple"
)

func nothingToSeeHere() {
	fmt.Println("There is no panic here.")
	fmt.Println("Shoo shoo.")
}

func innocentOperation(mode int) {
	switch mode {
	case 1:
		simple.ILoveCrashing()
	case 2:
		concurrent.ILoveCrashing()
	case 3:
		recovering.ILoveCrashing()
	default:
		nothingToSeeHere()
	}
}

func main() {
	mode := flag.Int("mode", 0, "panics: 0=not, 1=simple, 2=concurrently, 3=when revocering")
	flag.Parse()

	innocentOperation(*mode)
}
