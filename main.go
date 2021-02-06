package main

import (
	"log"

	"github.com/gopherty/blog/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}
