package mysql

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"article/config"
	"article/pkg/adding"
	"article/pkg/listing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/agiledragon/gomonkey/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func TestNewStorage(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	db.Close()

	tests := []struct {
		name    string
		cfgdb   config.MySQL
		open    func() (*sql.DB, error)
		want    *sql.DB
		wantErr bool
	}{
		{
			name:    "success mysql reachable",
			cfgdb:   config.MySQL{DSN: "success-dsn"},
			open:    func() (*sql.DB, error) { return db, nil },
			want:    db,
			wantErr: false,
		},
		{
			name:    "error mysql unreachable",
			cfgdb:   config.MySQL{DSN: "error-dsn"},
			open:    func() (*sql.DB, error) { return nil, errors.New("dial tcp: connection refused") },
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := gomonkey.ApplyFunc(sql.Open,
				func(driverName, dsn string) (*sql.DB, error) {
					return tt.open()
				})
			defer p.Reset()

			got, err := NewStorage(tt.cfgdb)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStorage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			require.NotNil(t, got)
			require.Equal(t, tt.want, got.db)
		})
	}
}

func TestStorage_CreateArticle(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	s := &Storage{db: db}

	tests := []struct {
		name      string
		ar        adding.Article
		mockSetup func()
		want      int
		wantErr   bool
	}{
		{
			name: "success",
			ar:   adding.Article{AuthorID: 1, Title: "testtitle", Body: "testbody"},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT article SET author_id = ?, title = ?, body = ? ").
					ExpectExec().
					WithArgs(1, "testtitle", "testbody").
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			want:    1,
			wantErr: false,
		},
		{
			name: "error prepare context",
			ar:   adding.Article{AuthorID: 1, Title: "testtitle", Body: "testbody"},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT article SET author_id = ?, title = ?, body = ? ").
					WillReturnError(sql.ErrConnDone)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error exec context",
			ar:   adding.Article{AuthorID: 1, Title: "testtitle", Body: "testbody"},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT article SET author_id = ?, title = ?, body = ? ").
					ExpectExec().
					WithArgs(1, "testtitle", "testbody").
					WillReturnError(errors.New("exec error"))
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "error get last insert id",
			ar:   adding.Article{AuthorID: 1, Title: "testtitle", Body: "testbody"},
			mockSetup: func() {
				mock.ExpectPrepare("INSERT article SET author_id = ?, title = ?, body = ? ").
					ExpectExec().
					WithArgs(1, "testtitle", "testbody").
					WillReturnResult(sqlmock.NewErrorResult(errors.New("get last insert id error")))
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := s.CreateArticle(context.Background(), tt.ar)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.CreateArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Storage.CreateArticle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStorage_ReadArticles(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer db.Close()

	s := &Storage{db: db}

	tests := []struct {
		name      string
		lfga      listing.FilterGetArticle
		mockSetup func()
		want      []listing.Article
		wantErr   bool
	}{
		{
			name: "success",
			lfga: listing.FilterGetArticle{AuthorName: "nametest", Search: "test", Limit: 10, Page: 1},
			mockSetup: func() {
				rows := mock.NewRows([]string{"ar.id", "au.id", "au.name", "ar.title", "ar.body"}).
					AddRow(1, 2, "nametest", "titletest", "testbody")
				mock.ExpectQuery("SELECT ar.id, au.id, au.name, ar.title, ar.body FROM article ar LEFT JOIN author au ON ar.author_id = au.id WHERE au.name LIKE CONCAT('%', ?, '%') AND MATCH(title, body) AGAINST(?) ORDER BY ar.created_at DESC LIMIT ? OFFSET  ?").
					WithArgs("nametest", "test", int64(10), int64(1)).
					WillReturnRows(rows)
			},
			want:    []listing.Article{{ID: 1, Author: listing.Author{ID: 2, Name: "nametest"}, Title: "titletest", Body: "testbody"}},
			wantErr: false,
		},
		{
			name: "success condition > 0",
			lfga: listing.FilterGetArticle{AuthorName: "nametest", Limit: 10, Page: 1},
			mockSetup: func() {
				rows := mock.NewRows([]string{"ar.id", "au.id", "au.name", "ar.title", "ar.body"}).
					AddRow(1, 2, "nametest", "titletest", "testbody")
				mock.ExpectQuery("SELECT ar.id, au.id, au.name, ar.title, ar.body FROM article ar LEFT JOIN author au ON ar.author_id = au.id WHERE au.name LIKE CONCAT('%', ?, '%') ORDER BY ar.created_at DESC LIMIT ? OFFSET  ?").
					WithArgs("nametest", int64(10), int64(1)).
					WillReturnRows(rows)
			},
			want:    []listing.Article{{ID: 1, Author: listing.Author{ID: 2, Name: "nametest"}, Title: "titletest", Body: "testbody"}},
			wantErr: false,
		},
		{
			name: "error query row context",
			lfga: listing.FilterGetArticle{AuthorName: "nametest", Limit: 10, Page: 1},
			mockSetup: func() {
				mock.ExpectQuery("SELECT ar.id, au.id, au.name, ar.title, ar.body FROM article ar LEFT JOIN author au ON ar.author_id = au.id WHERE au.name LIKE CONCAT('%', ?, '%') ORDER BY ar.created_at DESC LIMIT ? OFFSET  ?").
					WithArgs("nametest", int64(10), int64(1)).
					WillReturnError(errors.New("error query row context"))
			},
			wantErr: true,
		},
		{
			name: "error scan",
			lfga: listing.FilterGetArticle{AuthorName: "nametest", Limit: 10, Page: 1},
			mockSetup: func() {
				rows := mock.NewRows([]string{"ar.id", "au.id", "au.name", "ar.title", "ar.body"}).
					AddRow(1, "error", "nametest", "titletest", "testbody")
				mock.ExpectQuery("SELECT ar.id, au.id, au.name, ar.title, ar.body FROM article ar LEFT JOIN author au ON ar.author_id = au.id WHERE au.name LIKE CONCAT('%', ?, '%') ORDER BY ar.created_at DESC LIMIT ? OFFSET  ?").
					WithArgs("nametest", int64(10), int64(1)).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			got, err := s.ReadArticles(context.Background(), tt.lfga)
			if (err != nil) != tt.wantErr {
				t.Errorf("Storage.ReadArticles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Storage.ReadArticles() = %v, want %v", got, tt.want)
			}
		})
	}
}
