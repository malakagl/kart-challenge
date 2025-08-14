package handlers

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var couponFiles = []string{
	"couponbase1.gz",
	"couponbase2.gz",
	"couponbase3.gz",
}

func dataFilePath(name string) string {
	_, thisFile, _, _ := runtime.Caller(0) // 0 = this function
	baseDir := filepath.Join(filepath.Dir(thisFile), "../../promocodes/")
	return filepath.Join(baseDir, name)
}

func IsPromoCodeValid(code string) bool {
	log.Println("Validating promo code:", code)
	if len(code) < 8 || len(code) > 10 {
		return false
	}

	foundCount := 0
	for _, path := range couponFiles {
		path = dataFilePath(path)
		ok, err := containsCodeInGzip(path, code)
		if err != nil {
			log.Println(err)
			return false
		}

		if ok {
			foundCount++
			if foundCount >= 2 {
				return true
			}
		}
	}

	return false
}

func containsCodeInGzip(path, code string) (bool, error) {
	log.Println("checking in file", path)
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer func() { _ = f.Close() }()

	gr, err := gzip.NewReader(f)
	if err != nil {
		return false, err
	}
	defer func() { _ = gr.Close() }()

	scanner := bufio.NewScanner(gr)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	for scanner.Scan() {
		c := strings.TrimSpace(scanner.Text())
		if c == code {
			return true, nil
		}
	}

	if scanner.Err() != nil {
		return false, err
	}

	return false, nil
}
