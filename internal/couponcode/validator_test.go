package couponcode_test

import (
	"compress/gzip"
	"os"
	"testing"

	"github.com/malakagl/kart-challenge/internal/couponcode"
)

// helper: create a plain text file with given lines
func createTempFile(t *testing.T, lines []string) string {
	t.Helper()
	file, err := os.CreateTemp("", "coupon_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	for _, line := range lines {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			t.Fatal(err)
		}
	}
	return file.Name()
}

// helper: create a gzipped file with given lines
func createTempGzipFile(t *testing.T, lines []string) string {
	t.Helper()
	file, err := os.CreateTemp("", "coupon_*.gz")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	gw := gzip.NewWriter(file)
	defer gw.Close()

	for _, line := range lines {
		_, err := gw.Write([]byte(line + "\n"))
		if err != nil {
			t.Fatal(err)
		}
	}
	return file.Name()
}

func TestValidateCouponCode_PlainText(t *testing.T) {
	// create files
	file1 := createTempFile(t, []string{"ABC12345", "XYZ98765"})
	defer os.Remove(file1)

	file2 := createTempFile(t, []string{"ABC12345", "LMN11111"})
	defer os.Remove(file2)

	file3 := createTempFile(t, []string{"QWE22222"})
	defer os.Remove(file3)

	v := couponcode.NewValidator([]string{file1, file2, file3})

	// code present in 2 files → should return true
	if !v.ValidateCouponCode("ABC12345") {
		t.Error("Expected true, got false")
	}

	// code present in 1 file → should return false
	if v.ValidateCouponCode("QWE22222") {
		t.Error("Expected false, got true")
	}

	// code not present → should return false
	if v.ValidateCouponCode("NOTFOUND") {
		t.Error("Expected false, got true")
	}

	// code too short → should return false
	if v.ValidateCouponCode("SHORT") {
		t.Error("Expected false for short code")
	}
}

func TestValidateCouponCode_GzipFiles(t *testing.T) {
	file1 := createTempGzipFile(t, []string{"ABC12345", "XYZ98765"})
	defer os.Remove(file1)

	file2 := createTempGzipFile(t, []string{"ABC12345", "LMN11111"})
	defer os.Remove(file2)

	v := couponcode.NewValidator([]string{file1, file2})

	// code present in 2 files → should return true
	if !v.ValidateCouponCode("ABC12345") {
		t.Error("Expected true, got false")
	}

	// code present in 1 file → should return false
	if v.ValidateCouponCode("LMN11111") {
		t.Error("Expected false, got true")
	}
}

func TestValidateCouponCode_MixedFiles(t *testing.T) {
	file1 := createTempFile(t, []string{"ABC12345", "XYZ98765"})
	defer os.Remove(file1)

	file2 := createTempGzipFile(t, []string{"ABC12345", "LMN11111"})
	defer os.Remove(file2)

	v := couponcode.NewValidator([]string{file1, file2})

	// code present in both → should return true
	if !v.ValidateCouponCode("ABC12345") {
		t.Error("Expected true, got false")
	}

	// code present in only one → should return false
	if v.ValidateCouponCode("LMN11111") {
		t.Error("Expected false, got true")
	}
}
