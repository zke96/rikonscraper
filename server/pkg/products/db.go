package products

import (
	"context"
	"log"
	"rikonscraper/pkg/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func GetAllProducts() ([]types.Product, error) {
	products := []types.Product{}
	rows, err := pool.Query(context.Background(), "SELECT * FROM rikonscraper.products;")
	if err != nil {
		log.Printf("Failed to get products from db, %s", err)
		return products, err
	}
	products, err = pgx.CollectRows(rows, pgx.RowToStructByName[types.Product])
	if err != nil {
		log.Printf("Failed to scan products, %s", err)
		return products, err
	}

	return products, nil
}

func GetAllPartsForProduct(id uuid.UUID) ([]types.Product, error) {
	parts := []types.Product{}
	rows, err := pool.Query(context.Background(), "SELECT id, display, url, product_code FROM rikonscraper.parts WHERE parent=$1;", id)
	if err != nil {
		log.Printf("Failed to get products from db, %s", err)
		return parts, err
	}
	parts, err = pgx.CollectRows(rows, pgx.RowToStructByName[types.Product])
	if err != nil {
		log.Printf("Failed to scan products, %s", err)
		return parts, err
	}

	return parts, nil
}

func GetPartByID(id uuid.UUID) (types.Product, error) {
	part := types.Product{}

	rows, err := pool.Query(context.Background(), "SELECT id, display, url, product_code FROM rikonscraper.parts WHERE id=$1;", id)
	if err != nil {
		log.Printf("Failed to get products from db, %s", err)
		return part, err
	}
	part, err = pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[types.Product])
	if err != nil {
		log.Printf("Failed to scan products, %s", err)
		return part, err
	}

	return part, nil
}
