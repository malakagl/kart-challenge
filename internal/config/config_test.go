package config

import (
	"os"
	"testing"
)

func TestLoadConfig_ValidFile(t *testing.T) {
	yamlContent := `
server:
  host: "localhost"
  port: 8080
database:
  host: "dbhost"
  port: 5432
  name: "testdb"
  user: "testuser"
  password: "testpass"
  sslEnabled: true
  debug: false
  type: "postgres"
logging:
  debug: true
couponCodeConfig:
  unzipped: false
  filePaths:
    - "file1.gz"
    - "file2.gz"
    - "file3.gz"
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.Write([]byte(yamlContent)); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	cfg := LoadConfig(tmpFile.Name())
	if cfg.Server.Port != 8080 {
		t.Errorf("expected server port 8080, got %d", cfg.Server.Port)
	}
	if cfg.Database.Type != "postgres" {
		t.Errorf("expected database type postgres, got %s", cfg.Database.Type)
	}
	if !cfg.Logging.Debug {
		t.Errorf("expected logging debug true")
	}
	if cfg.CouponCodeConfig.Unzipped != false {
		t.Errorf("expected unzipped false, got true")
	}
	if len(cfg.CouponCodeConfig.FilePaths) != 3 {
		t.Errorf("expected 3 coupon code files, got %d", len(cfg.CouponCodeConfig.FilePaths))
	}
	if cfg.CouponCodeConfig.FilePaths[0] != "file1.gz" {
		t.Errorf("expected first coupon code file to be 'file1.gz', got %s", cfg.CouponCodeConfig.FilePaths[0])
	}
	if cfg.CouponCodeConfig.FilePaths[1] != "file2.gz" {
		t.Errorf("expected second coupon code file to be 'file2.gz', got %s", cfg.CouponCodeConfig.FilePaths[1])
	}
	if cfg.CouponCodeConfig.FilePaths[2] != "file3.gz" {
		t.Errorf("expected third coupon code file to be 'file3.gz', got %s", cfg.CouponCodeConfig.FilePaths[2])
	}
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	invalidPath := "nonexistent.yaml"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("expected panic for missing file")
		}
	}()
	LoadConfig(invalidPath)
}
