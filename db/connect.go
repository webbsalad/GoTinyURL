package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/webbsalad/GoTinyURL/config"
)

type DBConnection struct {
	Config config.ConfigDatabase
	Conn   *pgx.Conn
}

func (db *DBConnection) Connect() error {
	connParams := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s",
		db.Config.Host, db.Config.Port, db.Config.Name, db.Config.User, db.Config.Password)

	conn, err := pgx.Connect(context.Background(), connParams)
	if err != nil {
		return err
	}

	db.Conn = conn
	return nil
}

func (db *DBConnection) Close() {
	if db.Conn != nil {
		db.Conn.Close(context.Background())
	}
}
