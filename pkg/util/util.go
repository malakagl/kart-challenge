package util

import (
	"path/filepath"
	"runtime"
	"strconv"
)

// AbsoluteFilePath constructs an absolute file path based on relative path.
// The function returns the absolute path to the specified file.
func AbsoluteFilePath(file, relativePath string) string {
	_, thisFile, _, _ := runtime.Caller(0) // 0 = this function
	baseDir := filepath.Join(filepath.Dir(thisFile), relativePath)
	return filepath.Join(baseDir, file)
}

func StringToUint(s string) (uint, error) {
	t, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(t), nil
}
