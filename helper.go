package main

import "fmt"

// Helper function to print the syntax tree
func PrintSyntaxTree(node *TreeNode, indent int) {
	if node == nil {
		return
	}

	// Print indentation
	for i := 0; i < indent; i++ {
		fmt.Print("  ")
	}

	// Print node information
	switch node.NodeKind {
	case StmtK:
		PrintStmtNode(node)
	case ExpK:
		PrintExpNode(node)
	}
	fmt.Println()

	// Print children
	for i := 0; i < 3; i++ {
		if node.Children[i] != nil {
			PrintSyntaxTree(node.Children[i], indent+1)
		}
	}

	// Print siblings
	if node.Sibling != nil {
		PrintSyntaxTree(node.Sibling, indent)
	}
}

func PrintStmtNode(node *TreeNode) {
	switch node.StmtKind {
	case IfK:
		fmt.Print("If")
	case RepeatK:
		fmt.Print("Repeat")
	case AssignK:
		fmt.Printf("Assign to: %s", node.Name)
	case ReadK:
		fmt.Printf("Read: %s", node.Name)
	case WriteK:
		fmt.Print("Write")
	}
}

func PrintExpNode(node *TreeNode) {
	switch node.ExpKind {
	case OpK:
		fmt.Printf("Op: %s", node.Op)
	case ConstK:
		fmt.Printf("Const: %d", node.Value)
	case IdK:
		fmt.Printf("Id: %s", node.Name)
	}
}

func getNumChildNodes(node *TreeNode) int {
	var count int
	for i := 0; i < 3; i++ {
		if node.Children[i] != nil {
			count++
		}
	}
	return count
}
