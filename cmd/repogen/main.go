package main

import (
	"fmt"
	"graphy/cmd/repogen/repogen"
	"log"
	"os"
)

// Cypher repository generator tool. This tool should be run with the following command:
//
// > go run graphy/cmd/repogen/main.go path/to/repo.yaml
//
// This tool creates a basic repository for your methods which are based from an given interface
// as provided in the YAML file.
func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: return type")
		fmt.Println("repogen path/to/repository.yaml")
		os.Exit(1)
	}

	yamlPath := os.Args[1]
	gen, err := repogen.Generate(yamlPath)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(gen)
}
