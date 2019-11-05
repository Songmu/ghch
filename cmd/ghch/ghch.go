package main

import (
	"context"
	"log"
	"os"

	"github.com/Songmu/ghch"
)

func main() {
	log.SetFlags(0)
	if err := ghch.Run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
