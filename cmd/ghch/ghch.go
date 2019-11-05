package main

import (
	"context"
	"log"
	"os"

	"github.com/Songmu/ghch"
)

func main() {
	log.SetFlags(0)
	err := (&ghch.CLI{ErrStream: os.Stderr, OutStream: os.Stdout}).Run(context.Background(), os.Args[1:])
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
