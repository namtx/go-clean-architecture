package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/namtx/go-clean-architecture/article"
	"github.com/namtx/go-clean-architecture/models"
	"github.com/sirupsen/logrus"
)

type mysqlArticlesRepository struct {
	Conn *sql.DB
}

func NewMysqlArticlesRepository(conn *sql.DB) article.Repository {
	return &mysqlArticlesRepository{Conn: conn}
}

func (m *mysqlArticlesRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Article, error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			logrus.Error(err)
		}
	}()

	result := make([]*models.Article, 0)
	for rows.Next() {
		t := new(models.Article)
		err = rows.Scan(&t.ID, &t.Title, &t.Content, &t.UpdatedAt, &t.CreatedAt)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlArticlesRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Article, string, error) {
	query := `SELECT id, title, content, updated_at, created_at FROM articles WHERE created_at > ? ORDER BY created_at LIMIT ?`
	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}
	res, err := m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	nextCursor := ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return res, nextCursor, nil
}

func (m *mysqlArticlesRepository) GetByID(ctx context.Context, id int64) (res *models.Article, err error) {
	query := `SELECT id, title, content, created_at, updated_at FROM articles WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return nil, models.ErrNotFound
	}

	return
}

func (m *mysqlArticlesRepository) GetByTitle(ctx context.Context, title string) (res *models.Article, err error) {
	query := `SELECT id, title, content, updated_at, created_at FROM articles where title = ?`
	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}
	if len(list) > 0 {
		return list[0], nil
	}

	return nil, models.ErrNotFound
}

func (m *mysqlArticlesRepository) Store(ctx context.Context, a *models.Article) error {
	query := `INSERT articles SET title=?, content=?, updated_at=?, created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	a.ID = lastID

	return nil
}

func (m *mysqlArticlesRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE from articles WHERE ID = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		err = fmt.Errorf("Weird Behavior. Total Affected: %d", rowsAffected)

		return err
	}

	return nil
}
func (m *mysqlArticlesRepository) Update(ctx context.Context, a *models.Article) error {
	query := `UPDATE articles SET title = ?, content = ?, updated_at = ? WHERE ID = ?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.UpdatedAt, a.ID)
	if err != nil {
		return err
	}
	affectedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows != 1 {
		err = fmt.Errorf("Weird Behavior. Total Affected: %d", affectedRows)

		return err
	}

	return nil
}

func DecodeCursor(encodedTime string) (time.Time, error) {
	data, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(data)
	t, err := time.Parse(time.RFC3339Nano, timeString)

	return t, err
}

func EncodeCursor(t time.Time) string {
	timeString := t.Format(time.RFC3339Nano)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
