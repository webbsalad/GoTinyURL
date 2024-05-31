package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/webbsalad/GoTinyURL/db"
)

func FetchKeyByValue(dbConn *db.DBConnection, tableName string, searchValue string) (string, error) {
	querySelect := fmt.Sprintf(`SELECT key FROM "%s" WHERE value = $1`, tableName)
	row := dbConn.Conn.QueryRow(context.Background(), querySelect, searchValue)
	var key string
	err := row.Scan(&key)
	if err != nil {
		return "", err
	}

	queryUpdate := fmt.Sprintf(`UPDATE "%s" SET count = count + 1 WHERE value = $1;`, tableName)
	_, err = dbConn.Conn.Exec(context.Background(), queryUpdate, searchValue)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(map[string]string{"key": key})
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func FetchCountByValue(dbConn *db.DBConnection, tableName string, searchValue string) (string, error) {
	querySelect := fmt.Sprintf(`SELECT count FROM "%s" WHERE value = $1`, tableName)
	row := dbConn.Conn.QueryRow(context.Background(), querySelect, searchValue)
	var count string
	err := row.Scan(&count)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(map[string]string{"count": count})
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
