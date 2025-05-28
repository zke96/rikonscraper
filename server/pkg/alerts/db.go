package alerts

import (
	"context"
	"log"
	"rikonscraper/pkg/types"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func CreateAlert(partID uuid.UUID, email string) error {
	if _, err := pool.Exec(context.Background(), "INSERT INTO rikonscraper.alerts (email, part_id) VALUES ($1, $2) ON CONFLICT DO NOTHING;", email, partID); err != nil {
		log.Printf("Failed to insert alert into db, %s", err)
		return err
	}

	return nil
}

func GetAlertsByEmail(email string) ([]types.AlertRecord, error) {
	alerts := []types.AlertRecord{}
	rows, err := pool.Query(context.Background(),
		"SELECT alerts.id, alerts.email, alerts.part_id, parts.url, parts.display FROM rikonscraper.alerts LEFT JOIN rikonscraper.parts ON alerts.part_id = parts.id WHERE email=$1;", email)
	if err != nil {
		log.Printf("Failed to get alerts from db, %s", err)
		return alerts, err
	}
	alerts, err = pgx.CollectRows(rows, pgx.RowToStructByName[types.AlertRecord])
	if err != nil {
		log.Printf("Failed to scan alerts, %s", err)
		return alerts, err
	}

	return alerts, nil
}

func GetAllAlerts() ([]types.AlertRecord, error) {
	alerts := []types.AlertRecord{}
	rows, err := pool.Query(context.Background(),
		"SELECT alerts.id, alerts.email, alerts.part_id, parts.url, parts.display FROM rikonscraper.alerts LEFT JOIN rikonscraper.parts ON alerts.part_id = parts.id;")
	if err != nil {
		log.Printf("Failed to get alerts from db, %s", err)
		return alerts, err
	}
	alerts, err = pgx.CollectRows(rows, pgx.RowToStructByName[types.AlertRecord])
	if err != nil {
		log.Printf("Failed to scan alerts, %s", err)
		return alerts, err
	}

	return alerts, nil
}

func DeleteAlert(id uuid.UUID) error {
	_, err := pool.Exec(context.Background(), "DELETE FROM rikonscraper.alerts WHERE id=$1", id)
	if err != nil {
		log.Printf("Failed to remove alert from db, %s", err)
		return err
	}

	return nil
}
