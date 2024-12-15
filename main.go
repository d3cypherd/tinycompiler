package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var SCAN_PARSE_FLAG int

const (
	SCAN int = iota
	PARSE
)

func ScanFromBox(box *widget.Entry, displayBox *widget.TextGrid) {
	code := box.Text
	s := newScanner(*bufio.NewReader(strings.NewReader(code)))

	if !s.Scan() {
		fmt.Println("Scanning Failed.")
		displayBox.SetText("Scanning Failed.")
		return
	}
	tokens := s.PrintTokens()
	displayBox.SetText(tokens)
	// fmt.Print(tokens)
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

func main() {
	// Create the app and main window
	myApp := app.NewWithID("com.mycompany.myapp")
	myWindow := myApp.NewWindow("GUI with Fyne")

	// Textbox for input
	leftEntry := widget.NewMultiLineEntry()
	leftEntry.SetPlaceHolder("Enter your text here...")

	// Placeholder for the tree diagram
	rightTextGrid := widget.NewTextGridFromString("Tree diagram will be displayed here")
	rightTextGrid.ShowLineNumbers = true

	// Horizontal box for textbox and tree diagram
	// horizontalBox := container.NewHBox(
	// 	container.NewGridWrap( // Left VBox for text box
	// 		fyne.NewSize(300, 600),
	// 		leftEntry, // Text input
	// 	),
	// 	container.NewVBox( // Right VBox for tree diagram
	// 		layout.NewSpacer(),                 // Space above
	// 		container.NewCenter(rightTextGrid), // Tree placeholder
	// 		layout.NewSpacer(),                 // Space below
	// 	),
	// )

	// Bottom buttons
	button1 := widget.NewButton("SCAN", func() {
		ScanFromBox(leftEntry, rightTextGrid)
	})
	button2 := widget.NewButton("Button 2", func() {
		rightTextGrid.SetText("You pressed Button 2")
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
				ScanFromFile(file, rightTextGrid)
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
	splitContainer := container.NewHSplit(leftEntry, rightTextGrid)
	splitContainer.SetOffset(0.5)
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
	myWindow.Resize(fyne.NewSize(800, 600))
	myWindow.ShowAndRun()
}
