package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"fyne.io/x/fyne/widget/diagramwidget"
)

var SCAN_PARSE_FLAG int

const (
	SCAN int = iota
	PARSE
)

func ScanFromBox(box *widget.Entry, displayBox *widget.TextGrid) (bool, []string) {
	code := box.Text
	s := newScanner(*bufio.NewReader(strings.NewReader(code)))

	if !s.Scan() {
		fmt.Println("Scanning Failed.")
		displayBox.SetText(fmt.Sprintf("Scanning Failed:\n %v", s.errors))
		return false, s.errors
	}
	tokens := s.PrintTokens()
	displayBox.SetText(tokens)
	// fmt.Print(tokens)
	return true, []string{}
}

// ScanFromFile scans text from the uploaded file and displays tokens
func ScanFromFile(inputFile *os.File, displayBox *widget.TextGrid) {
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)
	s := newScanner(*reader)

	if !s.Scan() {
		fmt.Println("Scanning Failed.")
		displayBox.SetText("Scanning Failed.")
		return
	}
	tokens := s.PrintTokens()
	displayBox.SetText(tokens)
}

func ParseAndDisplayTree(code string, widget *diagramwidget.DiagramWidget) (bool, []string) {
	widget.DiagramElements = list.New()
	widget.Refresh()
	// First scan the code
	s := newScanner(*bufio.NewReader(strings.NewReader(code)))
	if !s.Scan() {
		fmt.Println("Scanning Failed:\n", s.errors)
		return false, s.errors
	}

	// Parse the tokens
	parser := NewParser(s.tokens)
	tree, errors := parser.Parse()

	if len(errors) > 0 {
		fmt.Println("Parsing errors:", errors)
		return false, errors
	}

	// Create and display the tree visualizer
	NewTreeVisualizer(tree, widget)
	return true, []string{}
}

func main() {
	// Create the app and main window
	myApp := app.NewWithID("com.mycompany.myapp")
	myWindow := myApp.NewWindow("tinycompiler")

	// Textbox for input
	leftEntry := widget.NewMultiLineEntry()
	leftEntry.SetPlaceHolder("Enter your text here...")

	// Placeholder for the tree diagram
	rightTextGrid := widget.NewTextGridFromString("Tree diagram will be displayed here")
	rightTextGrid.ShowLineNumbers = true

	diagramWidget := diagramwidget.NewDiagramWidget("diagram1")
	scrollContainer := container.NewScroll(diagramWidget)

	// Stack container to hold both components
	dynamicContainer := container.NewStack(rightTextGrid, scrollContainer)

	// Bottom buttons
	button1 := widget.NewButton("SCAN", func() {
		ScanFromBox(leftEntry, rightTextGrid)
		// Show the TextGrid and hide the diagram
		rightTextGrid.Show()
		scrollContainer.Hide()
	})

	button2 := widget.NewButton("Parse", func() {
		if status, errors := ParseAndDisplayTree(leftEntry.Text, diagramWidget); status {
			// Show the diagram and hide the TextGrid
			scrollContainer.Show()
			rightTextGrid.Hide()
		} else {
			rightTextGrid.SetText(fmt.Sprintf("Parsing Failed:\n %v\n", errors))
			rightTextGrid.Show()
			scrollContainer.Hide()
		}
	})

	// File upload button logic
	fileUploadButton := widget.NewButton("Upload File", func() {
		dialog.NewFileOpen(
			func(reader fyne.URIReadCloser, err error) {
				if err != nil || reader == nil {
					return // Handle cancel or error gracefully
				}

				file, err := os.Open(reader.URI().Path())
				if err != nil {
					fmt.Println("Failed to open file:", err)
					rightTextGrid.SetText("Failed to open file.")
					return
				}
				// Read the file contents
				data, err := io.ReadAll(file) // ReadAll reads the entire file
				if err != nil {
					fmt.Println("Error reading file:", err)
					return
				}
				// Convert bytes to string
				content := string(data)
				leftEntry.SetText(content)
			}, myWindow).Show()
	})

	// Buttons layout container
	buttonContainer := container.NewHBox(
		layout.NewSpacer(),
		button1,
		fileUploadButton,
		button2,
		layout.NewSpacer(),
	)

	// Split horizontally by half layout
	splitContainer := container.NewHSplit(leftEntry, dynamicContainer)
	splitContainer.SetOffset(0.3)

	// Add padding to the entire layout
	mainContainer := container.NewBorder(
		nil,             // No top widget
		buttonContainer, // Buttons at the bottom
		nil,             // No left widget
		nil,             // No right widget
		splitContainer,  // Padded horizontal box in the center
	)

	// Set window content and run the app
	myWindow.SetContent(mainContainer)
	myWindow.Resize(fyne.NewSize(1200, 900))
	myWindow.ShowAndRun()
}
