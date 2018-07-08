package currency

import (
	"database/sql"
	"fmt"
	"log"
)

// ChangeValue will alter the targetUserID's wallet by currencyDelta on the database pointed to by db
func ChangeValue(db *sql.DB, currencyDelta int, targetUserID string) {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO Currency (ID, Value) VALUES (%v, %v) ON DUPLICATE KEY UPDATE Value=Value + %v;",
		targetUserID, currencyDelta, currencyDelta))
	if err != nil {
		log.Printf("Error changing value(%v) for %v, %v", currencyDelta, targetUserID, err.Error())
		return
	}
}
