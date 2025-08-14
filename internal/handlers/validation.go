package handlers

import (
	"bufio"
	"compress/gzip"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
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

	resultChan := make(chan bool, len(couponFiles))
	wg := sync.WaitGroup{}
	for _, path := range couponFiles {
		path = dataFilePath(path)
		wg.Add(1)
		go containsCodeInGzip(path, code, &wg, resultChan)
	}
	wg.Wait()
	close(resultChan)
	count := 0
	for res := range resultChan {
		if res {
			count++
		}
	}

	return count > 1
}

func containsCodeInGzip(path, code string, wg *sync.WaitGroup, result chan bool) {
	defer wg.Done()
	log.Println("checking in file", path)
	f, err := os.Open(path)
	if err != nil {
		log.Println(err)
		result <- false
		return
	}
	defer func() { _ = f.Close() }()

	gr, err := gzip.NewReader(f)
	if err != nil {
		log.Println(err)
		result <- false
		return
	}
	defer func() { _ = gr.Close() }()

	scanner := bufio.NewScanner(gr)
	scanner.Buffer(nil, 1024*1024)
	for scanner.Scan() {
		c := strings.TrimSpace(scanner.Text())
		if c == code {
			log.Println(err)
			result <- true
			return
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	result <- false
}
