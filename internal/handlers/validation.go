package handlers

import (
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var (
	promoCodesSet1 sync.Map
	promoCodesSet2 sync.Map
	promoCodesSet3 sync.Map
)

func init() {
	promoCodesDir := os.Getenv("PROMO_CODES_DIR")
	if promoCodesDir == "" {
		promoCodesDir = "../../promocodes" // Default value if not set
	}
	var wg sync.WaitGroup
	filePath := filepath.Join(promoCodesDir, "couponbase1.gz")
	wg.Add(1)
	go readFile(filePath, &promoCodesSet1, &wg)

	filePath = filepath.Join(promoCodesDir, "couponbase2.gz")
	wg.Add(1)
	go readFile(filePath, &promoCodesSet2, &wg)

	filePath = filepath.Join(promoCodesDir, "couponbase3.gz")
	wg.Add(1)
	go readFile(filePath, &promoCodesSet3, &wg)
	wg.Wait()
}

func readFile(filePath string, codes *sync.Map, wg *sync.WaitGroup) {
	defer wg.Done()
	promos := readPromosFromFile(filePath)
	for _, promo := range promos {
		addCode(promo, codes)
	}
}

func addCode(promo string, codes *sync.Map) {
	code := strings.TrimSpace(promo)
	length := len(code)
	if 8 <= length && length <= 10 {
		codes.Store(code, true)
	}
}

func readPromosFromFile(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file %s: %v", file.Name(), err)
		return nil
	}

	defer func() { _ = file.Close() }()
	gzReader, err := gzip.NewReader(file)
	if err != nil {
		log.Printf("Failed to create gzip reader for file %s: %v", file.Name(), err)
		return nil
	}

	defer func() { _ = gzReader.Close() }()
	data, err := io.ReadAll(gzReader)
	if err != nil {
		log.Printf("Failed to read gzip file %s: %v", file.Name(), err)
		return nil
	}

	return strings.Split(string(data), "\n")
}

func ValidatePromoCode(code string) bool {
	count := 0
	if _, exists := promoCodesSet1.Load(code); exists {
		count++
	}
	if _, exists := promoCodesSet2.Load(code); exists {
		count++
	}
	if _, exists := promoCodesSet3.Load(code); exists {
		count++
	}

	return count > 1
}
