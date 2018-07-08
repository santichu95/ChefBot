package currency

import (
	"database/sql"
	"fmt"
	"log"
)

// CheckForCurrency will query the database and return the value of the targetUserID's wallet
// If that user is not in the database they will be added and given the starting amount of currency
func CheckForCurrency(db *sql.DB, targetUserID string) (int, error) {
	// Get info from database.
	stmtOut, err := db.Prepare("SELECT Value FROM Currency WHERE ID = ?")

	if err != nil {
		log.Printf("Error preparing query, %v", err.Error())
		return 0, err
	}

	var wealth int
	err = stmtOut.QueryRow(targetUserID).Scan(&wealth)

	// If the user is not in the database add them
	if err == sql.ErrNoRows {
		log.Print("No rows found")
		// Insert the value and print 0
		// TODO add a default value for users to start with
		wealth = 0
		_, err = db.Exec(fmt.Sprintf("INSERT INTO Currency (ID, Value) VALUES(%v, %v)", targetUserID, wealth))
		if err != nil {
			log.Printf("Error inserting new user into database, %v", err.Error())
			return 0, err
		}
	} else if err != nil {
		log.Printf("Error querying database, %v", err.Error())
		return 0, err
	}

	return wealth, nil
}
