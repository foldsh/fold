package main

import (
	"fmt"

	"github.com/foldsh/fold/version"
)

func main() {
	fmt.Println(version.FoldVersion.String())
}
