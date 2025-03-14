package loader

import (
	"archive/zip"
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/m4tthewde/swell/internal/class"
	"github.com/m4tthewde/swell/internal/logger"
)

type Loader struct {
	classPath []string
	classes   map[string]class.Class
}

func NewLoader(classPath []string) Loader {
	return Loader{
		classes:   make(map[string]class.Class),
		classPath: classPath,
	}
}

func (l *Loader) Load(ctx context.Context, className string) (*class.Class, error) {
	log := logger.FromContext(ctx)

	c, ok := l.classes[className]
	if ok {
		return &c, nil
	}

	log.Infof("loading %s", className)

	r, err := getReader(className, l.classPath)
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

func getReader(className string, classPath []string) (io.ReadCloser, error) {
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

	for _, classPathDir := range classPath {
		dirEntries, err := os.ReadDir(filepath.Join(".", classPathDir))
		if err != nil {
			return nil, fmt.Errorf("class path dir %s not found: %v", classPathDir, err)
		}

		for _, entry := range dirEntries {
			if entry.Name() == className+".class" {
				if !entry.IsDir() {
					return os.Open(filepath.Join(".", classPathDir, entry.Name()))
				}
			}
		}

	}

	return nil, errors.New("class not found: " + className)
}
