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
		case char == ' ' || char == '\n':
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

		case char == ';':
			tokens = append(tokens, Token{";", "SEMICOLON"})
			char, err = r.ReadByte()

		case char == '<':
			tokens = append(tokens, Token{"<", "LESSTHAN"})
			char, err = r.ReadByte()

		case char == '>':
			tokens = append(tokens, Token{">", "GREATERTHAN"})
			char, err = r.ReadByte()

		case char == '(':
			tokens = append(tokens, Token{"(", "OPENBRACKET"})
			char, err = r.ReadByte()

		case char == ')':
			tokens = append(tokens, Token{")", "CLOSEDBRACKET"})
			char, err = r.ReadByte()

		case char == '+':
			tokens = append(tokens, Token{"+", "PLUS"})
			char, err = r.ReadByte()

		case char == '-':
			tokens = append(tokens, Token{"-", "MINUS"})
			char, err = r.ReadByte()

		case char == '*':
			tokens = append(tokens, Token{"*", "MULT"})
			char, err = r.ReadByte()

		case char == '/':
			tokens = append(tokens, Token{"/", "DIV"})
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

		case char == '=':
			tokens = append(tokens, Token{string("="), "EQUAL"})
			char, err = r.ReadByte()

		case char >= '0' && char <= '9':
			var number []byte
			for char >= '0' && char <= '9' {
				number = append(number, char)
				char, err = r.ReadByte()
				if err != nil {
					panic(err)
				}
			}
			// TODO: store token in tokens struct {number, TokenType::NUMBER}
			tokens = append(tokens, Token{string(number), "NUMBER"})

		case (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z'):
			var identifier []byte
			for (char >= 'A' && char <= 'Z') || (char >= 'a' && char <= 'z') {
				identifier = append(identifier, char)
				// Next character
				char, err = r.ReadByte()
				if err != nil {
					panic(err)
				}
			}
			// reserved words
			switch string(identifier) {
			case "if":
				tokens = append(tokens, Token{"if", "IF"})
			case "then":
				tokens = append(tokens, Token{"then", "THEN"})
			case "end":
				tokens = append(tokens, Token{"end", "END"})
			case "repeat":
				tokens = append(tokens, Token{"repeat", "REPEAT"})
			case "until":
				tokens = append(tokens, Token{"until", "UNTIL"})
			case "read":
				tokens = append(tokens, Token{"read", "READ"})
			case "write":
				tokens = append(tokens, Token{"write", "WRITE"})
			default:
				tokens = append(tokens, Token{string(identifier), "IDENTIFIER"})
			}
		default:
			fmt.Println("compilation error: undefined character entered")
			return
		}
	}
	for _, token := range tokens {
		fmt.Printf("%v, %v\n", token.TokenValue, token.TokenType)
	}
}
