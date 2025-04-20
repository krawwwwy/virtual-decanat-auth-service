package utils

import (
	"3lab/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

//						 |
//						 v
//						func init() {
//    						sql.Register("postgres", &Driver{})
//						}

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", config.GetConnectionString())
	if err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected to database")
	return db, nil
}
