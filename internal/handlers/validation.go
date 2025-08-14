package handlers

import (
	"bufio"
	"compress/gzip"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

var couponFiles = []string{
	"couponbase1.gz",
	"couponbase2.gz",
	"couponbase3.gz",
}

var dbInstance *sql.DB

type CodeRow struct {
	FileIdx int
	Code    string
}

func init() {
	if _, err := os.Stat("./codes.db"); err == nil {
		log.Println("Database already exists, skipping initialization.")
		return
	}

	db, err := sql.Open("sqlite3", "./codes.db?_journal_mode=WAL&_synchronous=OFF")
	if err != nil {
		log.Fatal(err)
	}
	dbInstance = db

	r, err := db.Exec("PRAGMA page_size = 4096")
	log.Println("Setting page size to 4096:", r, err)
	r, err = db.Exec("VACUUM")
	log.Println("Vacuuming database:", r, err)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS codes (
		file_idx INT NOT NULL ,
		code VARCHAR(10) NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Channel for all codes from all files
	codeCh := make(chan CodeRow, 50000)
	var wg sync.WaitGroup

	// Start DB writer goroutine
	go dbWriter(db, codeCh)

	// Process files in parallel
	numWorkers := runtime.NumCPU()
	fileCh := make(chan int, len(couponFiles))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range fileCh {
				if err := processGzipFile(f, codeCh); err != nil {
					log.Printf("Error processing %d: %v", i, err)
				}
			}
		}()
	}

	for i := range couponFiles {
		fileCh <- i
	}
	close(fileCh)

	wg.Wait()
	close(codeCh) // Signal DB writer to stop
	fmt.Println("All files processed.")
}

func processGzipFile(idx int, out chan<- CodeRow) error {
	filePath := dataFilePath(couponFiles[idx])
	fmt.Printf("Reading %s...\n", filePath)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gz.Close()

	scanner := bufio.NewScanner(gz)
	for scanner.Scan() {
		code := strings.TrimSpace(scanner.Text())
		if code != "" && len(code) >= 8 && len(code) <= 10 {
			out <- CodeRow{FileIdx: idx, Code: code}
		}
	}

	return scanner.Err()
}

func dbWriter(db *sql.DB, in <-chan CodeRow) {
	batch := make([]CodeRow, 0, 50000)

	for row := range in {
		batch = append(batch, row)
		if len(batch) >= 2000 {
			if err := insertBatch(db, batch); err != nil {
				log.Fatalf("DB insert failed: %v", err)
			}
			batch = batch[:0]
		}
	}

	// Insert any leftover
	if len(batch) > 0 {
		if err := insertBatch(db, batch); err != nil {
			log.Fatalf("Final DB insert failed: %v", err)
		}
	}
}

func insertBatch(db *sql.DB, batch []CodeRow) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	vals := make([]interface{}, 0, len(batch)*2)
	placeholders := make([]string, 0, len(batch))

	for _, r := range batch {
		placeholders = append(placeholders, "(?, ?)")
		vals = append(vals, r.FileIdx, r.Code)
	}
	sqlStr := "INSERT INTO codes (file_idx, code) VALUES " + strings.Join(placeholders, ",")
	_, err = tx.Exec(sqlStr, vals...)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func dataFilePath(name string) string {
	_, thisFile, _, _ := runtime.Caller(0) // 0 = this function
	baseDir := filepath.Join(filepath.Dir(thisFile), "../../promocodes/")
	return filepath.Join(baseDir, name)
}

func GetDb() *sql.DB {
	if dbInstance == nil {
		db, err := sql.Open("sqlite3", "./codes.db?_journal_mode=WAL&_synchronous=OFF")
		if err != nil {
			log.Fatal(err)
		}
		dbInstance = db
	}
	return dbInstance
}

func IsPromoCodeValid(code string) bool {
	var count int
	db := GetDb()
	err := db.QueryRow(`
		SELECT COUNT(DISTINCT file_idx) FROM codes WHERE code = ?
	`, code).Scan(&count)
	if err != nil {
		return false
	}
	return count >= 2
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
