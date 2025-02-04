package loader

import (
	"archive/zip"
	"bufio"
	"errors"
	"io"
	"os"
	"strings"

	"github.com/m4tthewde/swell/internal/class"
)

type Loader struct {
	classes map[string]class.Class
}

func NewLoader() Loader {
	return Loader{classes: make(map[string]class.Class)}
}

func (l *Loader) Load(className string) (*class.Class, error) {
	r, err := getReader(className)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	reader := bufio.NewReader(r)

	class, err := class.NewClass(reader)
	if err != nil {
		return nil, err
	}

	l.classes[className] = *class

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

	return nil, errors.New("class not found in jmods: " + className)
}
