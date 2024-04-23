//go:build admin || development

package repository

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAdminSettings(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    rows := sqlmock.NewRows([]string{"personal_deduction", "k_receipt_limit"}).
        AddRow(50000.0, 30000.0)

    mock.ExpectQuery("SELECT personal_deduction, k_receipt_limit FROM admin_settings WHERE id = 1").
        WillReturnRows(rows)

    repo := NewPostgresAdminRepository(db)
    settings, err := repo.GetAdminSettings()
    assert.NoError(t, err)
    assert.Equal(t, float64(50000), settings.PersonalDeduction)
    assert.Equal(t, float64(30000), settings.KReceiptLimit)
}

func TestSetPersonalDeduction(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    mock.ExpectExec("UPDATE admin_settings SET personal_deduction = \\$1 WHERE id = 1").
        WithArgs(60000.0).
        WillReturnResult(sqlmock.NewResult(0, 1))

    repo := NewPostgresAdminRepository(db)
    err = repo.SetPersonalDeduction(60000.0)
    assert.NoError(t, err)
}

func TestSetKReceiptLimit(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    mock.ExpectExec("UPDATE admin_settings SET k_receipt_limit = \\$1 WHERE id = 1").
        WithArgs(35000.0).
        WillReturnResult(sqlmock.NewResult(0, 1))

    repo := NewPostgresAdminRepository(db)
    err = repo.SetKReceiptLimit(35000.0)
    assert.NoError(t, err)
}

