package databases

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

type SqlDatabase struct {
	connector *sql.DB
}

func GetSqlDatabse(dbHost, dbPort, dbUser, dbPassword, dbName string) (*SqlDatabase, error) {
	dataSourceName := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	return &SqlDatabase{db}, err
}

func CreateDatabaseConnect() (Database, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	db, err := GetSqlDatabse(dbHost, dbPort, dbUser, dbPassword, dbName)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db *SqlDatabase) GetPath(id int) (string, error) {
	var imagePath string
	err := db.connector.QueryRow("SELECT path FROM paths WHERE id = $1", id).Scan(&imagePath)

	// Проверка наличия ID в базе данных
	if err == nil {
		return imagePath, nil
	} else if err == sql.ErrNoRows {
		return "", fmt.Errorf("Such path wasn't found")
	} else {
		return "", fmt.Errorf("Error in process request")
	}
}
