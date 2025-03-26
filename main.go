package main

import (
	"context"
	"errors"
	"os"
	"strings"

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

	mainClassName, err := getMainClassName()
	if err != nil {
		log.Fatalln(err)
	}

	classPath, err := getClassPath()
	if err != nil {
		log.Fatalln(err)
	}

	log.Infow("executing main", "mainClass", mainClassName, "classPath", classPath)

	ctx := logger.OnContext(context.Background(), log)
	runner := jvm.NewRunner(classPath)

	err = runner.RunMain(ctx, mainClassName)
	if err != nil {
		log.Fatalln(err)
	}
}

func getMainClassName() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("no main class provided")
	}

	return os.Args[1], nil
}

func getClassPath() ([]string, error) {
	if len(os.Args) < 3 {
		return nil, errors.New("no class path provided")
	}

	return strings.Split(os.Args[2], ":"), nil
}
