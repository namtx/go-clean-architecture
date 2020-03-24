package article

import (
	"context"

	"github.com/namtx/go-clean-architecture/models"
)

type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Article, nextCursor string, err error)
	GetByID(ctx context.Context, id int64) (*models.Article, error)
	GetByTitle(ctx context.Context, title string) (*models.Article, error)
	Update(ctx context.Context, article *models.Article) error
	Store(ctx context.Context, article *models.Article) error
	Delete(ctx context.Context, id int64) error
}
