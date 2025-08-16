package util

import (
	"path/filepath"
	"runtime"
	"strconv"
)

// RelativeFilePath constructs a relative file path based on the current file's location.
// It takes a file name and a location offset, which is typically used to navigate to the
// desired directory structure relative to the current file's directory.
// The function returns the absolute path to the specified file.
func RelativeFilePath(file, locOffset string) string {
	_, thisFile, _, _ := runtime.Caller(0) // 0 = this function
	baseDir := filepath.Join(filepath.Dir(thisFile), locOffset)
	return filepath.Join(baseDir, file)
}

func StringToUint(s string) (uint, error) {
	t, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}

	return uint(t), nil
}
