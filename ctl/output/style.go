package output

import "github.com/fatih/color"

var (
	Red   = color.New(color.FgRed).SprintFunc()
	Green = color.New(color.FgGreen).SprintFunc()
	Blue  = color.New(color.FgBlue).SprintFunc()

	Bold = color.New(color.Bold).SprintFunc()

	BoldRed   = color.New(color.Bold, color.FgRed).SprintFunc()
	BoldGreen = color.New(color.Bold, color.FgGreen).SprintFunc()
)
