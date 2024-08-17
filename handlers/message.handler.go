package handler

import "github.com/jmoiron/sqlx"

type DB struct {
	conn *sqlx.DB
}

func (db *DB) CreateMessage(content string, userID int) (int, error) {
	var id int
	err := db.conn.QueryRow("INSERT INTO messages (content, user_id) VALUES ($1, $2) RETURNING id", content, userID).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}
