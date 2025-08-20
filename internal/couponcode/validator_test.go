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

	couponcode.SetCouponCodeFiles([]string{file1, file2, file3})

	// code present in 2 files → should return true
	if ok, err := couponcode.ValidateCouponCode(t.Context(), "ABC12345"); !ok || err != nil {
		t.Error("Expected true and no error, got ", ok, err)
	}

	// code present in 1 file → should return false
	if ok, err := couponcode.ValidateCouponCode(t.Context(), "QWE22222"); ok || err == nil {
		t.Error("Expected false and error, got ", ok, err)
	}

	// code not present → should return false
	if ok, err := couponcode.ValidateCouponCode(t.Context(), "NOTFOUND"); ok || err == nil {
		t.Error("Expected false and error, got ", ok, err)
	}

	// code too short → should return false
	if ok, err := couponcode.ValidateCouponCode(t.Context(), "SHORT"); ok || err == nil {
		t.Error("Expected false for short code, got ", ok, err)
	}
}

func TestValidateCouponCode_GzipFiles(t *testing.T) {
	file1 := createTempGzipFile(t, []string{"ABC12345", "XYZ98765"})
	defer os.Remove(file1)

	file2 := createTempGzipFile(t, []string{"ABC12345", "LMN11111"})
	defer os.Remove(file2)

	couponcode.SetCouponCodeFiles([]string{file1, file2})

	// code present in 2 files → should return true
	if ok, err := couponcode.ValidateCouponCode(t.Context(), "ABC12345"); !ok || err != nil {
		t.Error("Expected true and no error, got ", ok, err)
	}

	// code present in 1 file → should return false
	if ok, err := couponcode.ValidateCouponCode(t.Context(), "LMN11111"); ok || err == nil {
		t.Error("Expected false and error, got ", ok, err)
	}
}

func TestValidateCouponCode_MixedFiles(t *testing.T) {
	file1 := createTempFile(t, []string{"ABC12345", "XYZ98765"})
	defer os.Remove(file1)

	file2 := createTempGzipFile(t, []string{"ABC12345", "LMN11111"})
	defer os.Remove(file2)

	couponcode.SetCouponCodeFiles([]string{file1, file2})

	// code present in both → should return true
	if ok, err := couponcode.ValidateCouponCode(t.Context(), "ABC12345"); !ok || err != nil {
		t.Error("Expected true and no error, got ", ok, err)
	}

	// code present in only one → should return false
	if ok, err := couponcode.ValidateCouponCode(t.Context(), "LMN11111"); ok || err == nil {
		t.Error("Expected false and no error, got ", ok, err)
	}
}
