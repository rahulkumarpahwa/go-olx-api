package handlers

import "database/sql"

type Handlers struct {
	DB *sql.DB
}

func NewHanlders(db *sql.DB) *Handlers {
	return &Handlers{
		DB: db,
	}
}
