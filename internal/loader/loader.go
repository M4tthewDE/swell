package loader

import (
	"archive/zip"
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/logger"
)

type Loader struct {
	classes map[string]class.Class
}

func NewLoader() Loader {
	return Loader{classes: make(map[string]class.Class)}
}

func (l *Loader) Load(ctx context.Context, className string) (*class.Class, error) {
	log := logger.FromContext(ctx)

	c, ok := l.classes[className]
	if ok {
		return &c, nil
	}

	log.Infof("loading %s", className)

	r, err := getReader(className)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	reader := bufio.NewReader(r)

	class, err := class.NewClass(ctx, reader, className)
	if err != nil {
		return nil, err
	}

	l.classes[className] = *class

	log.Infof("loaded %s", className)

	return class, nil
}

func getReader(className string) (io.ReadCloser, error) {
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		return nil, errors.New("JAVA_HOME not set")
	}

	reader, err := zip.OpenReader(javaHome + "jmods/java.base.jmod")
	if err != nil {
		return nil, err
	}

	for _, f := range reader.File {
		if f.Name == "classes/"+strings.ReplaceAll(className, ".", "/")+".class" {
			r, err := f.Open()
			if err != nil {
				return nil, err
			}

			return r, nil
		}
	}

	dirEntries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}

	for _, entry := range dirEntries {
		if entry.Name() == className+".class" {
			if !entry.IsDir() {
				return os.Open("./" + entry.Name())
			}
		}
	}

	return nil, errors.New("class not found: " + className)
}
