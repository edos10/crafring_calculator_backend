package common

import (
	"database/sql"

	"github.com/tanimutomo/sqlfile"
)

func CreateTables(db *sql.DB) error {
	s := sqlfile.New()
	err := s.File("../common/data/create_tables.sql")
	if err != nil {
		return err
	}
	_, err = s.Exec(db)
	if err != nil {
		return err
	}
	return nil
}

func InsertData(db *sql.DB) error {
	s := sqlfile.New()
	err := s.File("../common/data/insert_test_data.sql")
	if err != nil {
		return err
	}
	_, err = s.Exec(db)
	if err != nil {
		return err
	}
	return nil
}
