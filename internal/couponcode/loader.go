package couponcode

import (
	"compress/gzip"
	"io"
	"os"
	"sync"
	"time"

	logging "github.com/malakagl/kart-challenge/pkg/logger"
)

// SetupCouponCodeFiles processes a list of gzip files, unzipping each one into a corresponding text file.
// experimental: This function is designed to handle multiple files concurrently, improving performance for large datasets.
func SetupCouponCodeFiles(filePaths []string) error {
	defer func(start time.Time) {
		logging.Logger.Info().Msgf("Coupon code files setup completed in %s", time.Since(start).String())
	}(time.Now())

	var wg sync.WaitGroup
	var errCh = make(chan error, len(filePaths))
	for _, filePath := range filePaths {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			if err := UnZipGzipFile(path); err != nil {
				logging.Logger.Error().Msgf("error unzipping file %s: %v", path, err)
			}
		}(filePath)
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			logging.Logger.Error().Msgf("error processing file: %v", err)
			return err
		}
	}

	return nil
}

func UnZipGzipFile(input string) error {
	output := input[:len(input)-3] + ".txt" // remove .gz extension
	// Open gzip file
	f, err := os.Open(input)
	if err != nil {
		logging.Logger.Error().Msgf("failed to open gzip file %s: %v", input, err)
		return err
	}
	defer f.Close()

	// Create gzip reader
	gr, err := gzip.NewReader(f)
	if err != nil {
		logging.Logger.Error().Msgf("failed to create gzip reader for %s: %v", input, err)
		return err
	}
	defer gr.Close()

	// Create output file
	out, err := os.Create(output)
	if err != nil {
		logging.Logger.Error().Msgf("failed to create output file %s: %v", output, err)
		return err
	}
	defer out.Close()

	// Copy decompressed data
	_, err = io.Copy(out, gr)
	if err != nil {
		logging.Logger.Error().Msgf("failed to copy data from %s to %s: %v", input, output, err)
		return err
	}

	logging.Logger.Error().Msgf("unzipped file %s to %s", input, output)
	return nil
}
