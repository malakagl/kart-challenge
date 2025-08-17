package couponcode

import (
	"bufio"
	"compress/gzip"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/malakagl/kart-challenge/pkg/log"
)

var couponCodeFiles []string

func SetCouponCodeFiles(f []string) {
	couponCodeFiles = f
}

func worker(ctx context.Context, path, code string, count *atomic.Int32, wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	var reader io.Reader = f

	// If file ends with .gz → wrap in gzip reader
	if strings.HasSuffix(strings.ToLower(filepath.Ext(path)), ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			log.Error().Msgf("Error creating gzip reader: %v", err)
			return
		}
		defer func() { _ = gz.Close() }()
		reader = gz
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			log.Debug().Msgf("Context done: %v", path)
			return
		default:
			if strings.TrimSpace(scanner.Text()) == code {
				if count.Add(1) >= 2 { // found in 2 files
					cancel() // stop all other workers
					return
				}
				return
			}
		}
	}
}

func ValidateCouponCode(code string) bool {
	log.Debug().Msgf("validating coupon code %s", code)
	defer func(start time.Time) {
		log.Debug().Msgf("validated coupon code in %s", time.Since(start).String())
	}(time.Now())

	if len(code) < 8 || len(code) > 10 {
		return false
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var count atomic.Int32

	for _, f := range couponCodeFiles {
		wg.Add(1)
		go worker(ctx, f, code, &count, &wg, cancel)
	}

	wg.Wait()
	return count.Load() >= 2
}
