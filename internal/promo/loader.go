package promo

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var couponFiles = []string{
	"couponbase1.gz",
	"couponbase2.gz",
	"couponbase3.gz",
}

type CodeRow struct {
	FileId int
	Code   string
}

func dataFilePath(name string) string {
	_, thisFile, _, _ := runtime.Caller(0) // 0 = this function
	baseDir := filepath.Join(filepath.Dir(thisFile), "../../promocodes/")
	return filepath.Join(baseDir, name)
}

func LoadCouponCodes() {
	user := "user"
	pass := "password"
	host := "localhost"
	port := "5432"
	name := "test"

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, pass, host, port, name,
	)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	files := couponFiles
	concurrency := 1 // number of concurrent files
	batchSize := 1000000

	jobs := make(chan struct {
		fileID int
		path   string
	}, len(files))

	var wg sync.WaitGroup

	// Start worker pool
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for job := range jobs {
				fmt.Printf("Processing file: %s\n", job.path)
				if err := copyFromGzipBatched(ctx, pool, job.fileID, job.path, batchSize); err != nil {
					log.Fatalf("Failed to copy file %s: %v", job.path, err)
				}
				fmt.Printf("Finished file: %s\n", job.path)
			}
		}()
	}

	// Send jobs
	for fileID, path := range files {
		jobs <- struct {
			fileID int
			path   string
		}{fileID + 1, path}
	}
	close(jobs)

	wg.Wait()
	fmt.Println("All files processed successfully.")
}

func copyFromGzipBatched(ctx context.Context, pool *pgxpool.Pool, fileID int, path string, batchSize int) error {
	path = dataFilePath(path)
	log.Printf("Copying %s to %d", path, fileID)
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	scanner := bufio.NewScanner(gzReader)
	batch := make([][]interface{}, 0, batchSize)

	for scanner.Scan() {
		batch = append(batch, []interface{}{fileID, scanner.Text()})

		if len(batch) >= batchSize {
			if err := copyBatch(ctx, pool, batch); err != nil {
				return err
			}
			batch = batch[:0] // reset batch
		}
	}

	// Insert remaining rows
	if len(batch) > 0 {
		if err := copyBatch(ctx, pool, batch); err != nil {
			return err
		}
	}

	return scanner.Err()
}

func copyBatch(ctx context.Context, pool *pgxpool.Pool, batch [][]interface{}) error {
	rows := &batchCopySource{batch: batch, idx: -1}
	_, err := pool.CopyFrom(ctx, pgx.Identifier{"coupon_codes"}, []string{"file_id", "code"}, rows)
	return err
}

// Implements CopyFromSource for a batch
type batchCopySource struct {
	batch [][]interface{}
	idx   int
}

func (b *batchCopySource) Next() bool {
	b.idx++
	return b.idx < len(b.batch)
}

func (b *batchCopySource) Values() ([]interface{}, error) {
	return b.batch[b.idx], nil
}

func (b *batchCopySource) Err() error {
	return nil
}
