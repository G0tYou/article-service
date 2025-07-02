package mysql

import (
	"context"
	"database/sql"
	"log"

	"article/config"
	"article/pkg/adding"

	_ "github.com/go-sql-driver/mysql"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(cfgdb config.MySQL) (*Storage, error) {
	var err error

	s := new(Storage)

	s.db, err = sql.Open("mysql", cfgdb.DSN)
	if err != nil {
		return s, err
	}

	log.Println("MySQL connected")

	return s, nil
}

func (s *Storage) CreateArticle(ctx context.Context, ar adding.Article) (int, error) {
	q := "INSERT article SET author_id = ?, title = ?, body = ?"

	stmt, err := s.db.PrepareContext(ctx, q)
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, ar.AuthorID, ar.Title, ar.Body)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
