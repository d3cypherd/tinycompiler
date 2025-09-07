package main

import (
	"bufio"
	"fmt"
	"io"
)

// TokenType represents different types of tokens using an enum
type TokenType int

// Token type constants defined as enum values
const (
	// Special tokens
	ERROR TokenType = iota // Using 0 as ERROR helps detect uninitialized tokens
	EOF

	// Reserved words
	IF
	THEN
	ELSE
	END
	REPEAT
	UNTIL
	READ
	WRITE

	// Special symbols
	SEMICOLON     // ;
	LESSTHAN      // <
	OPENBRACKET   // (
	CLOSEDBRACKET // )
	PLUS          // +
	MINUS         // -
	MULT          // *
	DIV           // /
	EQUAL         // =
	ASSIGN        // :=

	// Multi-character tokens
	NUMBER
	IDENTIFIER
)

// Token struct now uses the enum type
type Token struct {
	Value   string
	Type    TokenType
	LineNum int // Added for better error reporting
	CharNum int
}

// Scanner struct remains similar but with better organization
type Scanner struct {
	r       bufio.Reader
	tokens  []Token
	errors  []string
	CharNum int
	LineNum int
}

func (s *Scanner) PrintTokens() string {
	var output string
	for _, token := range s.tokens {
		output += fmt.Sprintf("%v, %v\n", token.Value, token.Type.String())
	}
	return output
}

// String method for TokenType provides readable token types for debugging
func (t TokenType) String() string {
	return [...]string{
		"ERROR",
		"EOF",
		"IF",
		"THEN",
		"ELSE",
		"END",
		"REPEAT",
		"UNTIL",
		"READ",
		"WRITE",
		"SEMICOLON",
		"LESSTHAN",
		"OPENBRACKET",
		"CLOSEDBRACKET",
		"PLUS",
		"MINUS",
		"MULT",
		"DIV",
		"EQUAL",
		"ASSIGN",
		"NUMBER",
		"IDENTIFIER",
	}[t]
}

// Helper functions remain the same
func isWhitespace(c byte) bool {
	return c == ' ' || c == '\n'
}

func isSingleOperator(c byte) bool {
	return c == ';' || c == '<' || c == '>' || c == '(' || c == ')' ||
		c == '+' || c == '-' || c == '*' || c == '/' || c == '='
}

func isAlphabet(c byte) bool {
	return (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z')
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

// getTokenType now returns TokenType instead of string
func getTokenType(c string) TokenType {
	// Map of reserved words to their token types
	reservedWords := map[string]TokenType{
		"if":     IF,
		"then":   THEN,
		"else":   ELSE,
		"end":    END,
		"repeat": REPEAT,
		"until":  UNTIL,
		"read":   READ,
		"write":  WRITE,
	}

	// First check if it's a reserved word
	if tokenType, ok := reservedWords[c]; ok {
		return tokenType
	}

	// Then check single operators
	singleOperators := map[string]TokenType{
		";": SEMICOLON,
		"<": LESSTHAN,
		"(": OPENBRACKET,
		")": CLOSEDBRACKET,
		"+": PLUS,
		"-": MINUS,
		"*": MULT,
		"/": DIV,
		"=": EQUAL,
	}

	if tokenType, ok := singleOperators[c]; ok {
		return tokenType
	}

	// Default case
	return IDENTIFIER
}

func newScanner(reader bufio.Reader) *Scanner {
	return &Scanner{
		r:       reader,
		CharNum: 0,
		LineNum: 1,
	}
}

func (s *Scanner) Read() (byte, error) {
	char, err := s.r.ReadByte()

	if char == '\n' {
		s.CharNum = 1
		s.LineNum++
	}

	s.CharNum++
	return char, err
}

func (s *Scanner) addToken(value string, tokenType TokenType) {
	s.tokens = append(s.tokens, Token{
		Value:   value,
		Type:    tokenType,
		LineNum: s.LineNum,
		CharNum: s.CharNum,
	})
}

func (s *Scanner) Scan() bool {
	char, err := s.Read()
	if err != nil {
		if err == io.EOF {
			return true
		}
		panic(err)
	}

	for {
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		switch {
		case isWhitespace(char):
			char, err = s.Read()
			continue

		case char == '{':
			for char != '}' {
				char, err = s.Read()
				if err != nil {
					if err == io.EOF {
						return s.error(fmt.Sprintf("unmatched '{'"))
					}
					panic(err)
				}
			}
			char, err = s.Read()

		case isSingleOperator(char):
			s.addToken(string(char), getTokenType(string(char)))
			char, err = s.Read()

		case char == ':':
			char, err = s.Read()
			if err != nil {
				panic(err)
			}
			if char == '=' {
				s.addToken(":=", ASSIGN)
				char, err = s.Read()
			} else {
				return s.error("':' not followed by '='")
			}

		case isNumber(char):
			var number []byte
			for isNumber(char) {
				number = append(number, char)
				char, err = s.Read()
				if err != nil && err != io.EOF {
					panic(err)
				}
			}
			s.addToken(string(number), NUMBER)

		case isAlphabet(char):
			var identifier []byte
			for isAlphabet(char) {
				identifier = append(identifier, char)
				char, err = s.Read()
				if err != nil && err != io.EOF {
					panic(err)
				}
			}
			word := string(identifier)
			s.addToken(word, getTokenType(word))

		default:
			return s.error("undefined character entered '" + string(char) + "'")
		}
	}
	return true
}

func (s *Scanner) error(msg string) bool {
	fmt.Printf("[%d:%d] compilation error: %s\n", s.LineNum, s.CharNum, msg)
	s.errors = append(s.errors, fmt.Sprintf("[%d:%d] compilation error: %s\n", s.LineNum, s.CharNum, msg))
	return false
}

func (s *Scanner) addError(msg string) {
	s.errors = append(s.errors, msg)
}
