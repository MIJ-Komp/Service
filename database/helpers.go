package database

import (
	"fmt"

	"gorm.io/gorm"
)

func CreateSequenceIfNotExists(db *gorm.DB, sequenceName string) error {
	rawSQL := fmt.Sprintf(`
	DO $$
	BEGIN
		IF NOT EXISTS (
			SELECT 1
			FROM pg_class c
			JOIN pg_namespace n ON n.oid = c.relnamespace
			WHERE c.relkind = 'S' AND c.relname = '%s'
		) THEN
			CREATE SEQUENCE %s
			START WITH 1
			INCREMENT BY 1
			MINVALUE 1
			CACHE 1;
		END IF;
	END
	$$;
	`, sequenceName, sequenceName)

	return db.Exec(rawSQL).Error
}

func GetNextInvoiceNumber(db *gorm.DB) (int64, error) {
	var number int64
	err := db.Raw("SELECT nextval('invoice_code_seq')").Scan(&number).Error
	return number, err
}
