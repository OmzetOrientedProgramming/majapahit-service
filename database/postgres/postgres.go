package postgres

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func Init() *sqlx.DB {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", os.Getenv("DB_USERNAME"), os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		logrus.Error(err)
		return nil
	}

	logrus.Info("PostgreSQL Connected Successfully")
	return db
}
