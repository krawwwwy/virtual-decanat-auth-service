package handlers

import (
	"3lab/models"
	"database/sql"
)

func GetUserRole(db *sql.DB, userID int) (*models.Role, error) {
	var role models.Role
	query := `
        SELECT r.id, r.name
        FROM roles r
        JOIN user_roles ur ON r.id = ur.role_id
        WHERE ur.user_id = $1
    `
	err := db.QueryRow(query, userID).Scan(&role.ID, &role.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &role, nil
}
