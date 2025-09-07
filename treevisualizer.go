package main

import (
	"fmt"
	"image/color"
	"math"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/diagramwidget"
)

type TreeVisualizer struct {
	widget.BaseWidget
	diagram *diagramwidget.DiagramWidget
	tree    *TreeNode
	nodes   map[*TreeNode]*diagramwidget.DiagramNode
	links   []*diagramwidget.BaseDiagramLink
}

func NewTreeVisualizer(tree *TreeNode, diagram *diagramwidget.DiagramWidget) *TreeVisualizer {
	v := &TreeVisualizer{
		diagram: diagram,
		tree:    tree,
		nodes:   make(map[*TreeNode]*diagramwidget.DiagramNode),
		links:   make([]*diagramwidget.BaseDiagramLink, 0),
	}
	v.ExtendBaseWidget(v)
	v.createDiagram()
	return v
}

func (v *TreeVisualizer) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(v.diagram)
}

func (v *TreeVisualizer) createDiagram() {
	// Create all nodes first
	levelOffsets := make(map[int]float32)
	v.createNodes(v.tree, 0, 0, 80, 0, levelOffsets)

	// Then create all links
	v.createLinks(v.tree)

	// Apply force-directed layout
	// for i := 0; i < 100; i++ {
	// 	diagramwidget.StepForceLayout(v.diagram, 300)
	// }
}

func (v *TreeVisualizer) createNodes(node *TreeNode, level, index int, xPos, yPos float32, levelOffsets map[int]float32) {
	if node == nil {
		return
	}

	// Ensure levelOffsets is initialized for this level
	if _, exists := levelOffsets[level]; !exists {
		levelOffsets[level] = 0
	}

	// Adjust xPos to avoid overlap
	if xPos < levelOffsets[level] {
		xPos = levelOffsets[level]
	}
	levelOffsets[level] = xPos + 200 // Update level offset with spacing for the next node

	// Create label content based on node type
	labelText := v.getNodeLabel(node)
	label := widget.NewLabel(labelText)
	label.Alignment = fyne.TextAlignCenter

	// Create diagram node
	diagNode := diagramwidget.NewDiagramNode(v.diagram, label, fmt.Sprintf("node-%p", node))
	// Set to red color if Leaf node
	if node.ExpKind == ConstK || node.ExpKind == IdK {
		diagNode.SetForegroundColor(color.RGBA{255, 0, 0, 255})
	}
	v.nodes[node] = &diagNode

	diagNode.Move(fyne.NewPos(xPos, yPos))

	// Process children
	numOfChildNodes := getNumChildNodes(node)
	var childSpacing float32 = 150.0 // Base spacing between children
	childXStart := xPos - (float32(numOfChildNodes-1) * childSpacing / 2)

	for i := 0; i < 3; i++ {
		if node.Children[i] != nil {
			childXPos := childXStart + float32(i)*childSpacing
			v.createNodes(node.Children[i], level+1, index, childXPos, yPos+100, levelOffsets)
		}
	}

	// Process siblings
	if node.Sibling != nil {
		siblingXPos := xPos + 360 // Sibling spacing
		v.createNodes(node.Sibling, level, index+1, siblingXPos, yPos, levelOffsets)
	}
}

func (v *TreeVisualizer) getNodeLabel(node *TreeNode) string {
	var nodeType, details string

	// Determine node type
	switch node.NodeKind {
	case StmtK:
		nodeType = "Statement"
		switch node.StmtKind {
		case IfK:
			details = "If"
		case RepeatK:
			details = "Repeat"
		case AssignK:
			details = fmt.Sprintf("Assign\n%s", node.Name)
		case ReadK:
			details = fmt.Sprintf("Read\n%s", node.Name)
		case WriteK:
			details = "Write"
		}
	case ExpK:
		nodeType = "Expression"
		switch node.ExpKind {
		case OpK:
			details = fmt.Sprintf("Op\n%s", node.Op)
		case ConstK:
			details = fmt.Sprintf("Const\n%d", node.Value)
		case IdK:
			details = fmt.Sprintf("Id\n%s", node.Name)
		}
	}

	return fmt.Sprintf("%s\n%s", nodeType, details)
}

func (v *TreeVisualizer) createLinks(node *TreeNode) {
	if node == nil {
		return
	}

	currentNode := *v.nodes[node]

	// Create links to children
	for i := 0; i < 3; i++ {
		if node.Children[i] != nil {
			childNode := *v.nodes[node.Children[i]]
			link := diagramwidget.NewDiagramLink(v.diagram, fmt.Sprintf("link-%p-%d", node, i))
			link.SetSourcePad(currentNode.GetEdgePad())
			link.SetTargetPad(childNode.GetEdgePad())
			v.links = append(v.links, link)
		}
	}

	// Create link to sibling
	if node.Sibling != nil {
		siblingNode := *v.nodes[node.Sibling]
		link := diagramwidget.NewDiagramLink(v.diagram, fmt.Sprintf("sibling-link-%p", node))
		link.SetSourcePad(currentNode.GetEdgePad())
		link.SetTargetPad(siblingNode.GetEdgePad())
		link.SetForegroundColor(color.RGBA{0, 128, 0, 255}) // Green for sibling links
		v.links = append(v.links, link)
	}

	// Recursively create links for children
	for i := 0; i < 3; i++ {
		v.createLinks(node.Children[i])
	}

	// Recursively create links for sibling
	v.createLinks(node.Sibling)
}

// Helper function to avoid overlap
func (v *TreeVisualizer) avoidOverlap(x, y float32, currentNode diagramwidget.DiagramNode) (float32, float32) {
	for _, existingNode := range v.nodes {
		if existingNode == &currentNode {
			continue
		}
		existingX, existingY := (*existingNode).Position().X, (*existingNode).Position().Y
		distance := math.Sqrt(math.Pow(float64(x-existingX), 2) + math.Pow(float64(y-existingY), 2))
		if distance < 50 { // Adjust 50 based on node size
			// Move slightly to avoid collision
			x += 10
			y += 10
		}
	}
	return x, y
}
