package common

import (
	"database/sql"
	"os"
)

func CreateTables(db *sql.DB) error {
	file, err := os.ReadFile("data/create_tables.sql")
	if err != nil {
		return err
	}
	sql := string(file)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func InsertData(db *sql.DB) error {
	file, err := os.ReadFile("data/insert_test_data.sql")
	if err != nil {
		return err
	}
	sql := string(file)
	_, err = db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
