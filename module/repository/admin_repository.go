package repository

import (
	"database/sql"

	"github.com/tankubopa777/assessment-tax/module/models"
)


type PostgresAdminRepository struct {
	db *sql.DB
}

func NewPostgresAdminRepository(db *sql.DB) *PostgresAdminRepository {
	return &PostgresAdminRepository{db: db}
}

func (r *PostgresAdminRepository) GetAdminSettings() (models.AdminSettings, error) {
    settings := models.AdminSettings{}
    err := r.db.QueryRow("SELECT personal_deduction, k_receipt_limit FROM admin_settings WHERE id = 1").Scan(&settings.PersonalDeduction, &settings.KReceiptLimit)
    if err != nil {
        return settings, err
    }
    return settings, nil
}

func (r *PostgresAdminRepository) SetPersonalDeduction(deduction float64) error {
    _, err := r.db.Exec("UPDATE admin_settings SET personal_deduction = $1 WHERE id = 1", deduction)
    return err
}

func (r *PostgresAdminRepository) SetKReceiptLimit(limit float64) error {
    _, err := r.db.Exec("UPDATE admin_settings SET k_receipt_limit = $1 WHERE id = 1", limit)
    return err
}
