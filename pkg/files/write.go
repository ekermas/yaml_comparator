package files

import (
	"os"
)

func CreateFile(path string) (*os.File, error) {
	return os.Create(path)
}
