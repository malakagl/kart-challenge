package e2e

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	config2 "github.com/malakagl/kart-challenge/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var configPath string
var cfg *config2.Config
var dbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	loadConfig()
	if !waitForPostgres() {
		log.Fatal("timed out waiting for postgres")
	}
	seedPostgresData()

	// Run tests
	code := m.Run()
	os.Exit(code)
}

func loadConfig() {
	flag.StringVar(&configPath, "config", "./config/config.default.yaml", "Path to config file")
	flag.Parse()
	log.Println("Loading config from ", configPath)
	var err error
	cfg, err = config2.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}
}

func waitForPostgres() bool {
	var db *pgxpool.Pool
	poolMaxWait := 60 * time.Second
	poolStart := time.Now()
	for {
		dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.Database.User, cfg.Database.Password, "localhost", cfg.Database.Port, cfg.Database.Name)
		var err error
		db, err = pgxpool.New(context.Background(), dbURL)
		if err == nil {
			err = db.Ping(context.Background())
		}
		if err == nil {
			log.Println("Postgres is ready!")
			break
		}

		if time.Since(poolStart) > poolMaxWait {
			log.Fatalf("Postgres did not start in %v: %v", poolMaxWait, err)
		}

		log.Println("Waiting for Postgres...", err)
		time.Sleep(1 * time.Second)
	}
	dbPool = db
	return db != nil
}

func seedPostgresData() {
	ctx := context.Background()
	_, err := dbPool.Exec(ctx, "SET search_path TO kart_challenge_it")
	if err != nil {
		log.Fatal("failed to set test schema", err)
	}

	_, err = dbPool.Exec(ctx, `INSERT INTO products (name, price, category) 
								VALUES ('Chicken Waffle', 13.25, 'Waffle')`)
	if err != nil {
		log.Fatal("seed product failed with error, ", err)
	}

	_, err = dbPool.Exec(ctx, `INSERT INTO product_images (product_id, thumbnail, mobile, tablet, desktop) 
								VALUES (1, '1/thumbnail.jpg', '1/mobile.jpg', '1/tablet.jpg', '1/desktop.jpg')`)
	if err != nil {
		log.Fatal("seed product images failed with error, ", err)
	}
}

func TestProductsAPI(t *testing.T) {
	url := fmt.Sprintf("http://%s:%d/products", cfg.Server.Host, cfg.Server.Port)
	resp, err := http.Get(url)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.Contains(t, string(body), "Unauthorized")
}
