package operations

import (
	"context"
	"fmt"
	"unicode"

	"github.com/webbsalad/GoTinyURL/db"
)

func AddItem(dbConn *db.DBConnection, tableName string, key string) (string, error) {

	lastValue, err := GetLastValue(dbConn, tableName)
	if err != nil {
		return "", err
	}

	newValue := generateNextValue(lastValue)

	query := fmt.Sprintf(`INSERT INTO "%s" ("key", "value") VALUES ($1, $2)`, tableName)
	_, err = dbConn.Conn.Exec(context.Background(), query, key, newValue)
	if err != nil {
		return "", err
	}
	generatedURL := "https://tiny-url-nu.vercel.app/" + newValue
	return generatedURL, nil
}

func GetLastValue(dbConn *db.DBConnection, tableName string) (string, error) {
	var lastValue string
	query := fmt.Sprintf(`SELECT "value" FROM "%s" ORDER BY "id" DESC LIMIT 1`, tableName)
	row := dbConn.Conn.QueryRow(context.Background(), query)
	err := row.Scan(&lastValue)
	if err != nil {
		if err.Error() != "no rows in result set" {
			return "", err
		}
		lastValue = ""
	}
	return lastValue, nil
}

func generateNextValue(lastValue string) string {
	if lastValue == "" {
		return "00000"
	}

	runes := []rune(lastValue)
	length := len(runes)
	carry := 1

	for i := length - 1; i >= 0; i-- {
		if carry == 0 {
			break
		}

		switch {
		case runes[i] == 'Z':
			runes[i] = '0'
			carry = 1
		case runes[i] == 'z':
			runes[i] = 'A'
			carry = 1
		case runes[i] == '9':
			runes[i] = 'a'
			carry = 1
		case unicode.IsDigit(runes[i]) || unicode.IsLower(runes[i]) || unicode.IsUpper(runes[i]):
			runes[i]++
			carry = 0
		default:
			carry = 0
		}
	}

	if carry == 1 {
		runes = append([]rune{'0'}, runes...)
	}

	return string(runes)
}
