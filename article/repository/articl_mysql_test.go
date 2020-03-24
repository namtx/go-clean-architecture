package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/namtx/go-clean-architecture/article/repository"
	"github.com/namtx/go-clean-architecture/models"
	"github.com/stretchr/testify/assert"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expcted when opening a stub database connection", err)
	}

	mockArticles := []models.Article{
		models.Article{
			ID: 1, Title: "Title 1", Content: "Content 1", UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		models.Article{
			ID: 2, Title: "Title 2", Content: "Content 2", UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "updated_at", "created_at"}).
		AddRow(mockArticles[0].ID, mockArticles[0].Title, mockArticles[0].Content, mockArticles[0].UpdatedAt, mockArticles[0].CreatedAt).
		AddRow(mockArticles[1].ID, mockArticles[1].Title, mockArticles[1].Content, mockArticles[1].UpdatedAt, mockArticles[1].CreatedAt)
	query := "SELECT id, title, content, updated_at, created_at FROM articles WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := repository.NewMysqlArticlesRepository(db)
	cursor := repository.EncodeCursor(mockArticles[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)

	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}
