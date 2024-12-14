package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main <input-filename> <output-filename>")
		return
	}

	inputFile, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	s := newScanner(*bufio.NewReader(inputFile))

	if !s.Scan() {
		fmt.Println("Scanning Failed.")
		return
	}

	tokens := s.tokens

	// Create and use the parser
	parser := NewParser(tokens)
	tree, errors := parser.Parse()

	// Check for errors
	if len(errors) > 0 {
		// Handle errors
		for _, err := range errors {
			fmt.Println(err)
		}
	} else {
		// Process the syntax tree
		fmt.Println("Syntax tree created successfully")
		PrintSyntaxTree(tree, 0)
	}

	// outputFile, err := os.Create(os.Args[2])
	// if err != nil {
	// 	panic(err)
	// }
	// defer outputFile.Close()
	//
	// writer := bufio.NewWriter(outputFile)
	//
	// for _, token := range s.tokens {
	// 	output := fmt.Sprintf("%v, %v\n", token.Value, token.Type)
	// 	if _, err := writer.WriteString(output); err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Print(output)
	// }
	//
	// writer.Flush()
}
