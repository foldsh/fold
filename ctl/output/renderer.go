package output

import "fmt"

// The Renderer interface is used to describe something which can be rendered to the
// terminal. The most obvious example is a line of text but this also includes things
// like structs that represents tables, or other reports.
type Renderer interface {
	Render() string
}

// Printing a line is a very common use case so the renderer for it is defined here.
type Line string

func (l Line) Render() string {
	return string(l)
}

// This is for printing a success message as the end of a command. It makes the formatting
// consistent across commands.
type Success string

func (s Success) Render() string {
	return fmt.Sprintf("\n%s", BoldGreen(s))
}

// This is for printing an error message. It makes the formatting consistent across commands.
type Error string

func (e Error) Render() string {
	return fmt.Sprintf("\n%s%s", BoldRed("Error: "), Bold(e))
}
