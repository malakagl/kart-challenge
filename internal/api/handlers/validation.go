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

	var count int
	err = conn.QueryRow(context.Background(), "select count(distinct file_id) from coupon_codes where code=$1;", code).Scan(&count)
	if err != nil {
		log.Println("Error querying promo code:", err)
		return false
	}

	return count > 1
}
