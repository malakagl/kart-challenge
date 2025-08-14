package handlers

import (
	"testing"
)

func TestValidatePromoCode(t *testing.T) {
	tests := []struct {
		name      string
		promoCode string
		expected  bool
	}{
		{"Empty code", "", false},
		{"Valid code 1", "HAPPYHRS", true},
		{"Valid code 2", "FIFTYOFF", true},
		{"Invalid code", "SUPER100", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPromoCodeValid(tt.promoCode)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
