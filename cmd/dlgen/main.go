package main

import (
	"fmt"
	"graphy/cmd/dlgen/generator"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("usage: return type")
		fmt.Println("dlgen graphy/transport/graphql/dataloader.GradeRoundLoader graphy/pkg/competition/rounds.Service:FindGradeRounds destination_gen.go")
		os.Exit(1)
	}

	//wd, err := os.Getwd()
	//if err != nil {
	//	_, _ = fmt.Fprintln(os.Stderr, err.Error())
	//	os.Exit(2)
	//}

	if err := generator.Generate(os.Args[1], os.Args[2], os.Args[3]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(2)
	}
}