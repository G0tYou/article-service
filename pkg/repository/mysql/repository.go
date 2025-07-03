package mysql

import (
	"context"
	"database/sql"
	"log"

	"article/config"
	"article/pkg/adding"
	"article/pkg/listing"

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

func (s *Storage) ReadArticles(ctx context.Context, lfga listing.FilterGetArticle) ([]listing.Article, error) {
	var lars []listing.Article
	var condition string

	if lfga.Search != "" {
		condition = "AND MATCH(title, body) AGAINST('" + lfga.Search + "')"
	}

	q := "SELECT ar.id, au.id, au.name, ar.title, ar.body FROM article ar LEFT JOIN author au ON ar.author_id = au.id WHERE  au.name LIKE CONCAT('%', ?, '%') " + condition + " ORDER BY ar.created_at DESC LIMIT ? OFFSET  ?"

	res, err := s.db.QueryContext(ctx, q, lfga.AuthorName, lfga.Limit, lfga.Page)
	if err != nil {
		return lars, err
	}
	defer res.Close()

	for res.Next() {
		var lar listing.Article

		if err := res.Scan(&lar.ID, &lar.Author.ID, &lar.Author.Name, &lar.Title, &lar.Title); err != nil {
			return lars, err
		}

		lars = append(lars, lar)
	}

	if err := res.Err(); err != nil {
		return lars, err
	}

	return lars, nil
}
