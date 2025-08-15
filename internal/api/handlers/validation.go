package handlers

import (
	"context"
	"log"

	"github.com/malakagl/kart-challenge/internal/db"
)

func IsPromoCodeValid(code string) bool {
	conn, err := db.Connect()
	if err != nil {
		log.Println(err)
		return false
	}

	rows, err := conn.Query(context.Background(), "SELECT DISTINCT file_id FROM coupon_codes WHERE code = $1;", code)
	if err != nil {
		log.Println("Error querying promo code:", err)
		return false
	}
	defer rows.Close()

	return len(rows.RawValues()) > 1
}
