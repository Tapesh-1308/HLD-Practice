package users

import (
	"database/sql"
	"time"
)

type User struct {
	ID        int
	Name      string
	Email     string
	Age       int
	CreatedAt time.Time
}

func AddUser(db *sql.DB, name, email string, age int) error {
	query := `INSERT INTO persons (name, email, age, created_at) VALUES (?, ?, ?, ?)`
	_, err := db.Exec(query, name, email, age, time.Now())
	return err
}

func GetUserByID(db *sql.DB, id int) (User, error) {
	var user User
	query := `SELECT id, name, email, age, created_at FROM persons WHERE id = ?`
	row := db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt)
	return user, err
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`SELECT id, name, email, age, created_at FROM persons`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Age, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func UpdateUser(db *sql.DB, user User) error {
	query := `UPDATE persons SET name = ?, email = ?, age = ? WHERE id = ?`
	_, err := db.Exec(query, user.Name, user.Email, user.Age, user.ID)
	return err
}

func DeleteUser(db *sql.DB, id int) error {
	query := `DELETE FROM persons WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}
