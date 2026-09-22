package utils

import (
	"path/filepath"
	"runtime"
)

func GetRootPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return "."
	}
	return filepath.Join(filepath.Dir(filename), "..")
}
