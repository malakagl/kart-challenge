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
)

type CouponValidator interface {
	ValidateCouponCode(code string) bool
}

type Validator struct {
	CouponCodeFiles []string
}

func NewValidator(couponCodeFiles []string) *Validator {
	return &Validator{CouponCodeFiles: couponCodeFiles}
}

func worker(ctx context.Context, path, code string, count *atomic.Int32, wg *sync.WaitGroup, cancel context.CancelFunc) {
	defer wg.Done()

	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	var reader io.Reader = f

	// If file ends with .gz → wrap in gzip reader
	if strings.HasSuffix(strings.ToLower(filepath.Ext(path)), ".gz") {
		gz, err := gzip.NewReader(f)
		if err != nil {
			return
		}
		defer gz.Close()
		reader = gz
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
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

func (v *Validator) ValidateCouponCode(code string) bool {
	if len(code) < 8 || len(code) > 10 {
		return false
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	var count atomic.Int32

	for _, f := range v.CouponCodeFiles {
		wg.Add(1)
		go worker(ctx, f, code, &count, &wg, cancel)
	}

	wg.Wait()
	return count.Load() >= 2
}
