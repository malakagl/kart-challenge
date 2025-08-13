package handlers

import (
	"os"
	"testing"
)

func init() {
	err := os.Setenv("PROMO_CODES_DIR", "../../promocodes")
	if err != nil {
		panic("Failed to set environment variable: " + err.Error())
	}
}

func TestValidatePromoCode(t *testing.T) {
	tests := []struct {
		name      string
		promoCode string
		expected  bool
	}{
		{"Empty code", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidatePromoCode(tt.promoCode)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
