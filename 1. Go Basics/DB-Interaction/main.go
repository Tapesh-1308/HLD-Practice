package main

import (
	"database/sql"
	users "db-interaction/user"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	db_url := "postgresql://tapesh:tapesh@localhost:5432/hoteldb?sslmode=disable"

	db, err := sql.Open("postgres", db_url)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot reach the database: ", err)
	}
	log.Println("Successfully connected to the database!")

	allUser, err := users.GetAllUsers(db)
	if err != nil {
		log.Fatal("Error fetching users: ", err)
	}

	for _, user := range allUser {
		log.Printf("User: %+v\n", user)
	}
}
