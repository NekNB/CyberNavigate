package postgres

import (
	"database/sql"
	"errors"

	userService "github.com/NekNB/CyberNavigate/backend/user-service/internal/services/user"
	"github.com/NekNB/CyberNavigate/backend/user-service/internal/storage"
	"github.com/NekNB/CyberNavigate/swagger/gen/article"
	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/sirupsen/logrus"
)

type PostgresStorage struct {
	db *sql.DB
}

var _ userService.ArticleMetaProvider = (*PostgresStorage)(nil)

func New(log *logrus.Logger, uri string) (*PostgresStorage, error) {
	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Errorf("%s", err.Error())
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("%s", err.Error())
	}
	return &PostgresStorage{db: db}, nil
}

// Выборка все article metadata
func (p *PostgresStorage) Articles() (*[]article.ArticleMetaData, error) {
	rows, err := p.db.Query(
		`
			SELECT uuid, title, status 
			FROM metadata;
		`,
	)
	if err != nil {
		return nil, err
	}

	var articleSlice []article.ArticleMetaData

	for rows.Next() {
		metadata := article.ArticleMetaData{}
		if err := rows.Scan(
			&metadata.Id,
			&metadata.Title,
			&metadata.Status,
		); err != nil {
			return nil, err
		}

		articleSlice = append(articleSlice, metadata)
	}

	return &articleSlice, nil
}

// Получение конкретного article metadata по uuid
func (p *PostgresStorage) ArticleByUUID(articleUUID string) (*article.ArticleMetaData, error) {
	stmt, err := p.db.Prepare(
		`
			SELECT uuid, title, status 
			FROM metadata
			WHERE uuid = $1;
		`,
	)
	if err != nil {
		return nil, err
	}

	var articleMetadata article.ArticleMetaData

	if err = stmt.QueryRow(articleUUID).Scan(
		&articleMetadata.Id,
		&articleMetadata.Title,
		&articleMetadata.Status,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrArticleNotFound
		}
		return nil, err
	}

	return &articleMetadata, nil
}

// Получение конкретного article text по uuid
func (p *PostgresStorage) ArticleTextIDByUUID(articleUUID string) (string, error) {

	var textIdNull sql.NullString
	if err := p.db.QueryRow(
		`
			SELECT text_id
			FROM metadata
			WHERE uuid = $1;
		`,
		articleUUID,
	).Scan(&textIdNull); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrArticleNotFound
		}
		return "", err
	}
	if textIdNull.Valid {
		return textIdNull.String, nil
	}
	return "", storage.ErrArticleTextNotCreatedYet
}

// Создание сущности Article
func (p *PostgresStorage) CreateArticle(articleTitle *string) (*article.ArticleMetaData, error) {
	var metadata article.ArticleMetaData

	if err := p.db.QueryRow(`
		INSERT INTO metadata (title)
		VALUES ($1)
		RETURNING uuid, title, status;
	`, articleTitle).
		Scan(
			&metadata.Id,
			&metadata.Title,
			&metadata.Status,
		); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.UniqueViolation {
			return nil, storage.ErrArticleExists
		}
		return nil, err
	}

	return &metadata, nil
}

// Обновление сущности Article по UUID
func (p *PostgresStorage) UpdateArticleByUUID(uuid, title, textID, status, videoUrl string) (*article.ArticleMetaData, error) {
	var metadata article.ArticleMetaData

	if err := p.db.QueryRow(`
		UPDATE metadata
		SET
			title = COALESCE(NULLIF($2, ''), title),
			text_id = COALESCE(NULLIF($3, ''), text_id),
			status = COALESCE(NULLIF($4, '')::article_status, status),
			video_url = COALESCE(NULLIF($5, ''), video_url)
		WHERE uuid = $1
		RETURNING uuid, title, status;
	`, uuid, title, textID, status, videoUrl).
		Scan(
			&metadata.Id,
			&metadata.Title,
			&metadata.Status,
		); err != nil {
		var pgerr *pq.Error
		if errors.As(err, &pgerr) && pgerr.Code == pqerror.UniqueViolation {
			return nil, storage.ErrArticleNotFound
		}
		return nil, err
	}

	return &metadata, nil
}
