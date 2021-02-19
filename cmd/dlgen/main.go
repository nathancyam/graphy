package main

import (
	"fmt"
	"graphy/cmd/dlgen/generator"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: return type")
		fmt.Println("dlgen []example.com/rounds.Round")
		os.Exit(1)
	}

	wd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}

	if err = generator.Generate(wd, os.Args[1]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}