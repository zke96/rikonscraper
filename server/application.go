package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"rikonscraper/pkg/alerts"
	"rikonscraper/pkg/config"
	"rikonscraper/pkg/products"
	"rikonscraper/pkg/routes"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

var cfg = config.LoadConfig()

func main() {
	log.SetOutput(os.Stdout)                     // Set output to the console
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Include date, time, and short file name

	cfgJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal config: %v", err)
	} else {
		log.Printf("Loaded config:\n%s", string(cfgJSON))
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s",
		cfg.PGHost, cfg.PGPort, cfg.PGUser, cfg.PGPassword, cfg.PGDBName)

	dbpool, err := pgxpool.New(context.Background(), psqlInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()
	if err := dbpool.Ping(context.Background()); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected!")
	go alerts.HandleAlerts(dbpool)
	go products.HandleUpdateProducts(dbpool)

	routes.RunServer()
}
