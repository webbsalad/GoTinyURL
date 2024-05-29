package operations

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/webbsalad/GoTinyURL/db"
)

func FetchKeyByValue(dbConn *db.DBConnection, tableName string, searchValue string) (string, error) {
	query := fmt.Sprintf(`SELECT "key" FROM "%s" WHERE "value" = $1`, tableName)
	row := dbConn.Conn.QueryRow(context.Background(), query, searchValue)

	var key string
	err := row.Scan(&key)
	if err != nil {
		return "", err
	}

	jsonData, err := json.Marshal(map[string]string{"key": key})
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
