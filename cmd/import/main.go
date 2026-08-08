package main

import (
	"context"
	"log"
	"os"

	"github.com/petrickS/BOP/internal/core/config"
	"github.com/petrickS/BOP/internal/core/db"
	"github.com/petrickS/BOP/internal/features/parser"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("usage: %s <path-to-prices.csv>", os.Args[0])
	}

	ctx := context.Background()
	cfg := config.Load()

	pool, err := db.NewPool(ctx, cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("connect to db: %v", err)
	}
	defer pool.Close()

	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("open file: %v", err)
	}
	defer f.Close()

	count, err := parser.ImportCSV(ctx, pool, f)
	if err != nil {
		log.Fatalf("import: %v", err)
	}

	log.Printf("imported %d offers", count)
}
