package main

import (
	"log"
	"os"

	"github.com/m4tthewde/swell/internal"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if len(os.Args) < 2 {
		log.Fatalln("no class provided")
	}

	className := os.Args[1]

	log.Printf("running class %s\n", className)
	err := internal.Run(className)
	if err != nil {
		log.Fatalln(err)
	}
}
