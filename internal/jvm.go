package internal

import (
	"fmt"
	"github.com/m4tthewde/swell/internal/class"
)

func Run(className string) error {
	classFilePath := fmt.Sprintf("%s.class", className)

	_, err := class.NewClass(classFilePath)
	if err != nil {
		return err
	}

	return nil
}
