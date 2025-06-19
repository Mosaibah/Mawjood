package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type ContentData struct {
	db *sql.DB
}

type Interface interface {
	CreateContent(ctx context.Context, content Content) (*Content, error)
	// GetContent(ctx context.Context, id string) (*Content, error)
	// UpdateContent(ctx context.Context, content Content) (*Content, error)
	// DeleteContent(ctx context.Context, id string) error
	// ListContents(ctx context.Context, pageSize int32, pageToken string) ([]Content, string, error)
	// SearchContents(ctx context.Context, query string, pageSize int32, pageToken string) ([]Content, string, error)
}

func New(db *sql.DB) Interface {
	return &ContentData{db: db}
}

type Content struct {
	ID             string
	Title          string
	Description    string
	Tags           []string
	Language       string
	DurationSeconds int32
	PublishedAt    time.Time
	ContentType    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	ExternalURL    string
	PlatformName   string
}

func (cd *ContentData) CreateContent(ctx context.Context, content Content) (*Content, error) {
	// Start a transaction to handle content and tags insertion
	tx, err := cd.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert the content
	insertContentQuery := `
		INSERT INTO contents (title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at, updated_at`

	now := time.Now()
	content.CreatedAt = now
	content.UpdatedAt = now

	err = tx.QueryRowContext(ctx, insertContentQuery,
		content.Title,
		content.Description,
		content.Language,
		content.DurationSeconds,
		content.PublishedAt,
		content.ContentType,
		content.CreatedAt,
		content.UpdatedAt,
		content.ExternalURL,
		content.PlatformName,
	).Scan(&content.ID, &content.CreatedAt, &content.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to insert content: %w", err)
	}

	// Handle tags insertion
	if len(content.Tags) > 0 {
		for _, tagName := range content.Tags {
			// Insert tag if it doesn't exist, or get existing tag ID
			var tagID string
			upsertTagQuery := `
				INSERT INTO tags (id, name) 
				VALUES (gen_random_uuid(), $1)
				ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
				RETURNING id`
			
			err = tx.QueryRowContext(ctx, upsertTagQuery, tagName).Scan(&tagID)
			if err != nil {
				return nil, fmt.Errorf("failed to upsert tag %s: %w", tagName, err)
			}

			// Link content to tag
			insertContentTagQuery := `
				INSERT INTO content_tags (content_id, tag_id)
				VALUES ($1, $2)
				ON CONFLICT (content_id, tag_id) DO NOTHING`
			
			_, err = tx.ExecContext(ctx, insertContentTagQuery, content.ID, tagID)
			if err != nil {
				return nil, fmt.Errorf("failed to link content to tag %s: %w", tagName, err)
			}
		}
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Successfully created content with ID: %s", content.ID)
	return &content, nil
}