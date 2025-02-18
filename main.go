package main

import (
	"context"
	"os"

	"github.com/m4tthewde/swell/internal/jvm"
	"github.com/m4tthewde/swell/internal/logger"
)

func main() {
	log, err := logger.NewLogger()
	if err != nil {
		panic(err)
	}

	if len(os.Args) < 2 {
		log.Fatalln("no class provided")
	}

	className := os.Args[1]

	log.Infof("running class %s", className)

	ctx := logger.OnContext(context.Background(), log)

	err = jvm.Run(ctx, className)
	if err != nil {
		log.Fatalln(err)
	}
}
