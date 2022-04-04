package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

type ConnectionInfo struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func NewPostgresConnection(info ConnectionInfo) (*sqlx.DB, error) {
	return sqlx.Connect(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%d user=%s dbname=%s sslmode=%s password=%s",
			info.Host, info.Port, info.Username, info.DBName, info.SSLMode, info.Password,
		),
	)
}
