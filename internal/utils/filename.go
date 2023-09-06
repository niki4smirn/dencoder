package utils

import (
	"path/filepath"
	"strings"
)

func FilenameWithoutExt(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}
