package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var couponFiles = []string{
	"couponbase1.gz",
	"couponbase2.gz",
	"couponbase3.gz",
}

var ctx = context.Background()

const poolSize = 8

func dataFilePath(name string) string {
	_, thisFile, _, _ := runtime.Caller(0) // 0 = this file
	baseDir := filepath.Join(filepath.Dir(thisFile), "../../promocodes/")
	return filepath.Join(baseDir, name)
}

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DB:           0,
		PoolSize:     poolSize,         // fewer but larger pipelines
		PoolTimeout:  60 * time.Second, // wait longer for free conn
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	})

	var wg sync.WaitGroup
	for i, f := range couponFiles {
		wg.Add(1)
		go func(fname string, idx int) {
			defer wg.Done()
			fileID := fmt.Sprintf("file%d", idx+1)
			log.Println("loading", fname)
			if err := processFile(rdb, fname, fileID); err != nil {
				log.Println("error processing", fname, ":", err)
			}
		}(f, i)
	}

	wg.Wait()

	// Example query
	word := "FIFTYOFF"
	count, err := rdb.SCard(ctx, "word:"+word).Result()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Word %q found in %d files\n", word, count)
}

func processFile(rdb *redis.Client, path, fileID string) error {
	path = dataFilePath(path)
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	// Channel for word batches
	wordCh := make(chan []string, poolSize)
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for words := range wordCh {
				if err := batchInsert(rdb, words, fileID); err != nil {
					log.Println("error inserting batch for", fileID, ":", err)
				}
			}
		}()
	}

	// Producer: scan file and send batches
	scanner := bufio.NewScanner(gz)
	const batchSize = 1000000
	words := make([]string, 0, batchSize)

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" {
			continue
		}
		words = append(words, word)

		if len(words) >= batchSize {
			wordCh <- words
			words = make([]string, 0, batchSize) // reset
		}
	}

	// leftover
	if len(words) > 0 {
		wordCh <- words
	}

	close(wordCh) // no more work
	wg.Wait()     // wait for workers to finish

	return scanner.Err()
}

func batchInsert(rdb *redis.Client, words []string, fileID string) error {
	pipe := rdb.Pipeline()
	for _, w := range words {
		pipe.SAdd(ctx, "word:"+w, fileID)
	}
	_, err := pipe.Exec(ctx)
	return err
}
