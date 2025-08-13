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
	promoCodes = make(map[string][]string)
	mu         sync.Mutex
)

func init() {
	promoCodesDir := os.Getenv("PROMO_CODES_DIR")
	if promoCodesDir == "" {
		promoCodesDir = "../../promocodes" // Default value if not set
	}
	files, err := os.ReadDir(promoCodesDir)
	if err != nil {
		log.Fatalf("Failed to read promocodes directory: %v", err)
	}

	var wg sync.WaitGroup
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".gz" {
			filePath := filepath.Join(promoCodesDir, file.Name())
			wg.Add(1)
			go readFile(file.Name(), filePath, &wg)
		}
	}
	wg.Wait()
}

func readFile(fileName, filePath string, wg *sync.WaitGroup) {
	defer wg.Done()
	promos := readPromosFromFile(filePath)
	if promos != nil {
		mu.Lock()
		defer mu.Unlock()
		promoCodes[fileName] = promos
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
	if len(code) == 0 || len(code) < 5 {
		return false
	}

	// TODO: Implement logic to check if the code exists in the loaded promo codes
	return true
}
