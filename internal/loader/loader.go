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
}

func NewLoader() Loader {
	return Loader{}
}

func (l *Loader) Load(className string) error {
	r, err := getReader(className)
	if err != nil {
		return err
	}

	defer r.Close()

	reader := bufio.NewReader(r)

	_, err = class.NewClass(reader)
	if err != nil {
		return err
	}

	return errors.New("not implemented: load")
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
