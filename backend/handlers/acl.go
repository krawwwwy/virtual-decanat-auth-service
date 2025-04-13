package handlers

import (
	"database/sql"
)

func CheckAccess(db *sql.DB, userID int, resource, permission string) (bool, error) {
	var storedPermission string
	query := "SELECT permission FROM acl WHERE user_id = $1 AND resource = $2"
	err := db.QueryRow(query, userID, resource).Scan(&storedPermission)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return storedPermission == permission, nil
}
