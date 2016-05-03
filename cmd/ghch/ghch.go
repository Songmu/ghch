package main

import (
	"os"

	"github.com/Songmu/ghch"
)

func main() {
	os.Exit((&ghch.CLI{ErrStream: os.Stderr, OutStream: os.Stdout}).Run(os.Args[1:]))
}
