package repository

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/tankubopa777/assessment-tax/module/models"
)

func TestAdminRepository(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    repo := NewPostgresAdminRepository(db)

    tests := []struct {
        name    string
        test    func()
        expect  func()
        wantErr bool
    }{
        {
            name: "GetAdminSettings Success",
            test: func() {
                rows := sqlmock.NewRows([]string{"personal_deduction", "k_receipt_limit"}).
                    AddRow(50000.0, 30000.0)
                mock.ExpectQuery("SELECT personal_deduction, k_receipt_limit FROM admin_settings WHERE id = 1").
                    WillReturnRows(rows)
            },
            expect: func() {
                settings, err := repo.GetAdminSettings()
                assert.NoError(t, err)
                assert.Equal(t, float64(50000), settings.PersonalDeduction)
                assert.Equal(t, float64(30000), settings.KReceiptLimit)
            },
            wantErr: false,
        },
        {
            name: "GetAdminSettings Error",
            test: func() {
                mock.ExpectQuery("SELECT personal_deduction, k_receipt_limit FROM admin_settings WHERE id = 1").
                    WillReturnError(errors.New("query error"))
            },
            expect: func() {
                settings, err := repo.GetAdminSettings()
                assert.Error(t, err)
                assert.Equal(t, models.AdminSettings{}, settings)
            },
            wantErr: true,
        },
        {
            name: "SetPersonalDeduction Success",
            test: func() {
                mock.ExpectExec("UPDATE admin_settings SET personal_deduction = \\$1 WHERE id = 1").
                    WithArgs(60000.0).
                    WillReturnResult(sqlmock.NewResult(0, 1))
            },
            expect: func() {
                err := repo.SetPersonalDeduction(60000.0)
                assert.NoError(t, err)
            },
            wantErr: false,
        },
        {
            name: "SetKReceiptLimit Success",
            test: func() {
                mock.ExpectExec("UPDATE admin_settings SET k_receipt_limit = \\$1 WHERE id = 1").
                    WithArgs(35000.0).
                    WillReturnResult(sqlmock.NewResult(0, 1))
            },
            expect: func() {
                err := repo.SetKReceiptLimit(35000.0)
                assert.NoError(t, err)
            },
            wantErr: false,
        },
    }

    // Run the tests
    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            tc.test()
            tc.expect()
        })
    }

    // Verify all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}
