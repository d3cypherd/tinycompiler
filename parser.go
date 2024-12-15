package main

import (
	"fmt"
	"strconv"
)

// TreeNode represents a node in the abstract syntax tree
type TreeNode struct {
	NodeKind NodeKind
	StmtKind StmtKind
	ExpKind  ExpKind
	Children [3]*TreeNode // Max 3 children needed for if-else statements
	Sibling  *TreeNode    // For statement sequences
	Value    int          // For number constants
	Name     string       // For identifiers
	Op       string       // For operators
	LineNum  int
}

// Enum declaration for NodeKind, StmtKind, ExpKind
type (
	NodeKind int
	StmtKind int
	ExpKind  int
)

// Enum NodeKind
const (
	StmtK NodeKind = iota
	ExpK
)

// Enum StmtKind
const (
	IfK StmtKind = iota
	RepeatK
	AssignK
	ReadK
	WriteK
)

// Enum ExpKind
const (
	OpK ExpKind = iota
	ConstK
	IdK
)

// Parser maintains the parsing state
type Parser struct {
	tokens  []Token
	current int
	errors  []string
}

// NewParser creates a new parser instance
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		errors:  make([]string, 0),
	}
}

// Parse initiates the parsing process
// program = stmt-sequence
func (p *Parser) Parse() (*TreeNode, []string) {
	tree := p.parseStmtSequence()
	if p.current < len(p.tokens)-1 { // -1 for EOF token
		p.addError("Extra tokens after program end")
	}
	return tree, p.errors
}

// parseStmtSequence implements stmt-sequence = statement {";" statement}
func (p *Parser) parseStmtSequence() *TreeNode {
	// Parse the first statement
	firstStmt := p.parseStatement()
	if firstStmt == nil {
		return nil
	}

	currentStmt := firstStmt

	// Parse any additional statements after semicolons
	for p.currentToken().Type == SEMICOLON {
		p.match(SEMICOLON)
		nextStmt := p.parseStatement()
		if nextStmt != nil {
			currentStmt.Sibling = nextStmt
			currentStmt = nextStmt
		}
	}

	return firstStmt
}

// parseStatement implements statement = if-stmt | repeat-stmt | assign-stmt | read-stmt | write-stmt
func (p *Parser) parseStatement() *TreeNode {
	switch p.currentToken().Type {
	case IF:
		return p.parseIfStmt()
	case REPEAT:
		return p.parseRepeatStmt()
	case IDENTIFIER:
		return p.parseAssignStmt()
	case READ:
		return p.parseReadStmt()
	case WRITE:
		return p.parseWriteStmt()
	default:
		p.addError(fmt.Sprintf("Unexpected token: %v at line %d",
			p.currentToken().Value, p.currentToken().LineNum))
		return nil
	}
}

// parseIfStmt implements if-stmt = "if" exp "then" stmt-sequence ["else" stmt-sequence] "end"
func (p *Parser) parseIfStmt() *TreeNode {
	node := &TreeNode{
		NodeKind: StmtK,
		StmtKind: IfK,
		LineNum:  p.currentToken().LineNum,
	}

	p.match(IF)
	node.Children[0] = p.parseExp()
	p.match(THEN)
	node.Children[1] = p.parseStmtSequence()

	// Handle optional else clause
	if p.currentToken().Type == ELSE {
		p.match(ELSE)
		node.Children[2] = p.parseStmtSequence()
	}

	p.match(END)
	return node
}

// parseRepeatStmt implements repeat-stmt = "repeat" stmt-sequence "until" exp
func (p *Parser) parseRepeatStmt() *TreeNode {
	node := &TreeNode{
		NodeKind: StmtK,
		StmtKind: RepeatK,
		LineNum:  p.currentToken().LineNum,
	}

	p.match(REPEAT)
	node.Children[0] = p.parseStmtSequence()
	p.match(UNTIL)
	node.Children[1] = p.parseExp()
	return node
}

// parseAssignStmt implements assign-stmt = identifier ":=" exp
func (p *Parser) parseAssignStmt() *TreeNode {
	node := &TreeNode{
		NodeKind: StmtK,
		StmtKind: AssignK,
		Name:     p.currentToken().Value,
		LineNum:  p.currentToken().LineNum,
	}

	p.match(IDENTIFIER)
	p.match(ASSIGN)
	node.Children[0] = p.parseExp()
	return node
}

// parseReadStmt implements read-stmt = "read" identifier
func (p *Parser) parseReadStmt() *TreeNode {
	node := &TreeNode{
		NodeKind: StmtK,
		StmtKind: ReadK,
		LineNum:  p.currentToken().LineNum,
	}

	p.match(READ)
	node.Name = p.currentToken().Value
	p.match(IDENTIFIER)
	return node
}

// parseWriteStmt implements write-stmt = "write" exp
func (p *Parser) parseWriteStmt() *TreeNode {
	node := &TreeNode{
		NodeKind: StmtK,
		StmtKind: WriteK,
		LineNum:  p.currentToken().LineNum,
	}

	p.match(WRITE)
	node.Children[0] = p.parseExp()
	return node
}

// parseExp implements exp = simple-exp [comparison-op simple-exp]
func (p *Parser) parseExp() *TreeNode {
	left := p.parseSimpleExp()

	// Check for optional comparison operator
	if p.isComparisonOp(p.currentToken().Type) {
		node := &TreeNode{
			NodeKind: ExpK,
			ExpKind:  OpK,
			Op:       p.currentToken().Value,
			LineNum:  p.currentToken().LineNum,
		}

		node.Children[0] = left
		p.advance()
		node.Children[1] = p.parseSimpleExp()
		return node
	}

	return left
}

// parseSimpleExp implements simple-exp = term {addop term}
func (p *Parser) parseSimpleExp() *TreeNode {
	node := p.parseTerm()

	// Handle repeated addop terms
	for p.isAddOp(p.currentToken().Type) {
		newNode := &TreeNode{
			NodeKind: ExpK,
			ExpKind:  OpK,
			Op:       p.currentToken().Value,
			LineNum:  p.currentToken().LineNum,
		}

		newNode.Children[0] = node
		p.advance()
		newNode.Children[1] = p.parseTerm()
		node = newNode
	}

	return node
}

// parseTerm implements term = factor {mulop factor}
func (p *Parser) parseTerm() *TreeNode {
	node := p.parseFactor()

	// Handle repeated mulop factors
	for p.isMulOp(p.currentToken().Type) {
		newNode := &TreeNode{
			NodeKind: ExpK,
			ExpKind:  OpK,
			Op:       p.currentToken().Value,
			LineNum:  p.currentToken().LineNum,
		}

		newNode.Children[0] = node
		p.advance()
		newNode.Children[1] = p.parseFactor()
		node = newNode
	}

	return node
}

// parseFactor implements factor = "(" exp ")" | number | identifier
func (p *Parser) parseFactor() *TreeNode {
	var node *TreeNode

	switch p.currentToken().Type {
	case OPENBRACKET:
		p.match(OPENBRACKET)
		node = p.parseExp()
		p.match(CLOSEDBRACKET)

	case NUMBER:
		node = &TreeNode{
			NodeKind: ExpK,
			ExpKind:  ConstK,
			Value:    p.parseNumber(p.currentToken().Value),
			LineNum:  p.currentToken().LineNum,
		}
		p.advance()

	case IDENTIFIER:
		node = &TreeNode{
			NodeKind: ExpK,
			ExpKind:  IdK,
			Name:     p.currentToken().Value,
			LineNum:  p.currentToken().LineNum,
		}
		p.advance()

	default:
		p.addError(fmt.Sprintf("Unexpected token in factor: %v at line %d",
			p.currentToken().Value, p.currentToken().LineNum))
	}

	return node
}

// Helper functions
func (p *Parser) currentToken() Token {
	if p.current >= len(p.tokens) {
		return Token{Type: EOF}
	}
	return p.tokens[p.current]
}

// Match token and consume
func (p *Parser) match(expected TokenType) bool {
	if p.currentToken().Type == expected {
		p.advance()
		return true
	}
	p.addError(fmt.Sprintf("Expected %v but got %v at line %d",
		expected.String(), p.currentToken().Type.String(), p.currentToken().LineNum))
	return false
}

func (p *Parser) advance() {
	p.current++
}

func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseNumber(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		p.addError(fmt.Sprintf("Invalid number: %s", s))
		return 0
	}
	return n
}

func (p *Parser) isComparisonOp(t TokenType) bool {
	return t == LESSTHAN || t == GREATERTHAN || t == EQUAL
}

func (p *Parser) isAddOp(t TokenType) bool {
	return t == PLUS || t == MINUS
}

func (p *Parser) isMulOp(t TokenType) bool {
	return t == MULT || t == DIV
}
