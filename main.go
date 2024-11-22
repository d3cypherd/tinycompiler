package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Token struct {
	TokenValue string
	TokenType  string
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\n'
}

func isSingleOperator(c byte) bool {
	return c == ';' || c == '<' || c == '>' || c == '(' || c == ')' || c == '+' || c == '-' || c == '*' || c == '/' || c == '='
}

func isAlphabet(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func getTokenType(c string) string {
	switch c {
	case ";":
		return "SEMICOLON"
	case "<":
		return "LESSTHAN"
	case ">":
		return "GREATERTHAN"
	case "(":
		return "OPENBRACKET"
	case ")":
		return "CLOSEDBRACKET"
	case "+":
		return "PLUS"
	case "-":
		return "MINUS"
	case "*":
		return "MULT"
	case "/":
		return "DIV"
	case "=":
		return "EQUAL"
	case "if":
		return "IF"
	case "then":
		return "THEN"
	case "end":
		return "END"
	case "repeat":
		return "REPEAT"
	case "until":
		return "UNTIL"
	case "read":
		return "READ"
	case "write":
		return "WRITE"
	default:
		return "IDENTIFIER"
	}
}

func main() {
	// Example: ./main code.txt
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./main <filename>")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var tokens []Token

	r := bufio.NewReader(file)
	char, err := r.ReadByte()
	if err != nil {
		panic(err)
	}

	for {
		if err != nil {
			if err == io.EOF {
				break // Exit the loop when the end of the file is reached
			}
			panic(err)
		}
		switch {
		case isWhitespace(char):
			char, err = r.ReadByte()
			continue

		case char == '{':
			for char != '}' {
				char, err = r.ReadByte()
				if err != nil {
					if err == io.EOF {
						panic("unmatched '{'")
					}
					panic(err)
				}
			}
			// Read next character after '}'
			char, err = r.ReadByte()

		case isSingleOperator(char):
			tokens = append(tokens, Token{string(char), getTokenType(string(char))})
			char, err = r.ReadByte()

		case char == ':':
			char, err = r.ReadByte()
			if err != nil {
				panic(err)
			}
			if char == '=' {
				tokens = append(tokens, Token{":=", "ASSIGN"})
				char, err = r.ReadByte()
			} else {
				fmt.Println("compilation error (ASSIGN): ':' not followed by '='.")
				return
			}

		case isNumber(char):
			var number []byte
			for isNumber(char) {
				number = append(number, char)
				char, err = r.ReadByte()
				if err != nil {
					panic(err)
				}
			}
			tokens = append(tokens, Token{string(number), "NUMBER"})

		case isAlphabet(char):
			var identifier []byte
			for isAlphabet(char) {
				identifier = append(identifier, char)
				// Next character
				char, err = r.ReadByte()
				if err != nil {
					panic(err)
				}
			}
			// reserved words or identifier
			tokens = append(tokens, Token{string(identifier), getTokenType(string(identifier))})

		default:
			fmt.Println("compilation error: undefined character entered")
			return
		}
	}
	for _, token := range tokens {
		fmt.Printf("%v, %v\n", token.TokenValue, token.TokenType)
	}
}
