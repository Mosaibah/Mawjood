package store

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type ContentData struct {
	db *sql.DB
}

type Interface interface {
	CreateContent(ctx context.Context, content Content) (*Content, error)
	GetContent(ctx context.Context, id string) (*Content, error)
	UpdateContent(ctx context.Context, content Content) (*Content, error)
	DeleteContent(ctx context.Context, id string) error
	ListContents(ctx context.Context, pageSize int32, pageToken string) ([]Content, string, error)
	SearchContents(ctx context.Context, query string, pageSize int32, pageToken string) ([]Content, string, error)
}

func New(db *sql.DB) Interface {
	return &ContentData{db: db}
}

type Content struct {
	ID              string
	Title           string
	Description     string
	Tags            []string
	Language        string
	DurationSeconds int32
	PublishedAt     time.Time
	ContentType     string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	ExternalURL     string
	PlatformName    string
	DeletedAt       *time.Time
}

func (cd *ContentData) CreateContent(ctx context.Context, content Content) (*Content, error) {
	tx, err := cd.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

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

	if len(content.Tags) > 0 {
		for _, tagName := range content.Tags {
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

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &content, nil
}

func (cd *ContentData) GetContent(ctx context.Context, id string) (*Content, error) {
	getContentQuery := `
		SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name, deleted_at
		FROM contents 
		WHERE id = $1 AND deleted_at IS NULL`

	var content Content
	var publishedAt, createdAt, updatedAt time.Time
	var description, language, url, platformName sql.NullString
	var durationSeconds sql.NullInt32
	var deletedAt sql.NullTime

	err := cd.db.QueryRowContext(ctx, getContentQuery, id).Scan(
		&content.ID,
		&content.Title,
		&description,
		&language,
		&durationSeconds,
		&publishedAt,
		&content.ContentType,
		&createdAt,
		&updatedAt,
		&url,
		&platformName,
		&deletedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("content with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	content.Description = description.String
	content.Language = language.String
	content.ExternalURL = url.String
	content.PlatformName = platformName.String
	content.DurationSeconds = durationSeconds.Int32
	content.PublishedAt = publishedAt
	content.CreatedAt = createdAt
	content.UpdatedAt = updatedAt
	if deletedAt.Valid {
		content.DeletedAt = &deletedAt.Time
	}

	tags, err := cd.getContentTags(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get content tags: %w", err)
	}
	content.Tags = tags

	return &content, nil
}

func (cd *ContentData) UpdateContent(ctx context.Context, content Content) (*Content, error) {
	tx, err := cd.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	updateContentQuery := `
		UPDATE contents 
		SET title = $1, description = $2, language = $3, duration_seconds = $4, published_at = $5, content_type = $6, updated_at = $7, url = $8, platform_name = $9
		WHERE id = $10 AND deleted_at IS NULL
		RETURNING created_at, updated_at`

	now := time.Now()
	content.UpdatedAt = now

	err = tx.QueryRowContext(ctx, updateContentQuery,
		content.Title,
		content.Description,
		content.Language,
		content.DurationSeconds,
		content.PublishedAt,
		content.ContentType,
		content.UpdatedAt,
		content.ExternalURL,
		content.PlatformName,
		content.ID,
	).Scan(&content.CreatedAt, &content.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("content with ID %s not found", content.ID)
		}
		return nil, fmt.Errorf("failed to update content: %w", err)
	}

	deleteTagsQuery := `DELETE FROM content_tags WHERE content_id = $1`
	_, err = tx.ExecContext(ctx, deleteTagsQuery, content.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to remove existing tags: %w", err)
	}

	if len(content.Tags) > 0 {
		for _, tagName := range content.Tags {
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

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &content, nil
}

func (cd *ContentData) DeleteContent(ctx context.Context, id string) error {
	now := time.Now()
	softDeleteQuery := `
		UPDATE contents 
		SET deleted_at = $1, updated_at = $1
		WHERE id = $2 AND deleted_at IS NULL`

	result, err := cd.db.ExecContext(ctx, softDeleteQuery, now, id)
	if err != nil {
		return fmt.Errorf("failed to soft delete content: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("content with ID %s not found or already deleted", id)
	}

	return nil
}

func (cd *ContentData) ListContents(ctx context.Context, pageSize int32, pageToken string) ([]Content, string, error) {
	if pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	var query string
	var args []interface{}

	if pageToken == "" {
		query = `
			SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name, deleted_at
			FROM contents 
			WHERE deleted_at IS NULL
			ORDER BY created_at DESC 
			LIMIT $1`
		args = []interface{}{pageSize + 1}
	} else {
		query = `
			SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name, deleted_at
			FROM contents 
			WHERE deleted_at IS NULL AND created_at < (SELECT created_at FROM contents WHERE id = $1)
			ORDER BY created_at DESC 
			LIMIT $2`
		args = []interface{}{pageToken, pageSize + 1}
	}

	rows, err := cd.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list contents: %w", err)
	}
	defer rows.Close()

	var contents []Content
	for rows.Next() {
		var content Content
		var publishedAt, createdAt, updatedAt time.Time
		var description, language, url, platformName sql.NullString
		var durationSeconds sql.NullInt32
		var deletedAt sql.NullTime

		err := rows.Scan(
			&content.ID,
			&content.Title,
			&description,
			&language,
			&durationSeconds,
			&publishedAt,
			&content.ContentType,
			&createdAt,
			&updatedAt,
			&url,
			&platformName,
			&deletedAt,
		)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan content row: %w", err)
		}

		content.Description = description.String
		content.Language = language.String
		content.ExternalURL = url.String
		content.PlatformName = platformName.String
		content.DurationSeconds = durationSeconds.Int32
		content.PublishedAt = publishedAt
		content.CreatedAt = createdAt
		content.UpdatedAt = updatedAt
		if deletedAt.Valid {
			content.DeletedAt = &deletedAt.Time
		}

		tags, err := cd.getContentTags(ctx, content.ID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get content tags: %w", err)
		}
		content.Tags = tags

		contents = append(contents, content)
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating over content rows: %w", err)
	}

	var nextPageToken string
	if len(contents) > int(pageSize) {
		contents = contents[:pageSize]
		nextPageToken = contents[len(contents)-1].ID
	}

	return contents, nextPageToken, nil
}

func (cd *ContentData) SearchContents(ctx context.Context, query string, pageSize int32, pageToken string) ([]Content, string, error) {
	if pageSize <= 0 {
		pageSize = 10
	}

	if pageSize > 100 {
		pageSize = 100
	}

	searchQuery := strings.TrimSpace(query)
	if searchQuery == "" {
		return []Content{}, "", nil
	}

	_, err := cd.db.ExecContext(ctx, "SET SESSION pg_trgm.similarity_threshold = 0.10")
	if err != nil {
		return nil, "", fmt.Errorf("failed to set similarity threshold: %w", err)
	}

	var sqlQuery string
	var args []interface{}

	if pageToken == "" {
		sqlQuery = `
			WITH content_with_tags AS (
				SELECT 
					c.id, c.title, c.description, c.language, c.duration_seconds, c.published_at, 
					c.content_type, c.created_at, c.updated_at, c.url, c.platform_name, c.deleted_at,
					STRING_AGG(t.name, ' ') as tag_text
				FROM contents c
				LEFT JOIN content_tags ct ON c.id = ct.content_id
				LEFT JOIN tags t ON ct.tag_id = t.id
				WHERE c.deleted_at IS NULL
				GROUP BY c.id, c.title, c.description, c.language, c.duration_seconds, c.published_at, 
					c.content_type, c.created_at, c.updated_at, c.url, c.platform_name, c.deleted_at
			)
			SELECT 
				id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name, deleted_at,
				GREATEST(
					SIMILARITY(LOWER(title), LOWER($1)),
					SIMILARITY(LOWER(description), LOWER($1)),
					SIMILARITY(LOWER(platform_name), LOWER($1)),
					COALESCE(SIMILARITY(LOWER(tag_text), LOWER($1)), 0)
				) as max_similarity
			FROM content_with_tags
			WHERE (
				LOWER(title) % LOWER($1) OR 
				LOWER(description) % LOWER($1) OR 
				LOWER(platform_name) % LOWER($1) OR
				LOWER(tag_text) % LOWER($1) OR
				title ILIKE $2 OR 
				description ILIKE $2 OR
				platform_name ILIKE $2 OR
				tag_text ILIKE $2
			)
			ORDER BY max_similarity DESC, created_at DESC 
			LIMIT $3`
		likeQuery := "%" + searchQuery + "%"
		args = []interface{}{searchQuery, likeQuery, pageSize + 1}
	} else {
		sqlQuery = `
			WITH content_with_tags AS (
				SELECT 
					c.id, c.title, c.description, c.language, c.duration_seconds, c.published_at, 
					c.content_type, c.created_at, c.updated_at, c.url, c.platform_name, c.deleted_at,
					STRING_AGG(t.name, ' ') as tag_text
				FROM contents c
				LEFT JOIN content_tags ct ON c.id = ct.content_id
				LEFT JOIN tags t ON ct.tag_id = t.id
				WHERE c.deleted_at IS NULL
				GROUP BY c.id, c.title, c.description, c.language, c.duration_seconds, c.published_at, 
					c.content_type, c.created_at, c.updated_at, c.url, c.platform_name, c.deleted_at
			)
			SELECT 
				id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name, deleted_at,
				GREATEST(
					SIMILARITY(LOWER(title), LOWER($1)),
					SIMILARITY(LOWER(description), LOWER($1)),
					SIMILARITY(LOWER(platform_name), LOWER($1)),
					COALESCE(SIMILARITY(LOWER(tag_text), LOWER($1)), 0)
				) as max_similarity
			FROM content_with_tags
			WHERE (
				LOWER(title) % LOWER($1) OR 
				LOWER(description) % LOWER($1) OR 
				LOWER(platform_name) % LOWER($1) OR
				LOWER(tag_text) % LOWER($1) OR
				title ILIKE $2 OR 
				description ILIKE $2 OR
				platform_name ILIKE $2 OR
				tag_text ILIKE $2
			)
			AND created_at < (SELECT created_at FROM content_with_tags WHERE id = $3)
			ORDER BY max_similarity DESC, created_at DESC 
			LIMIT $4`
		likeQuery := "%" + searchQuery + "%"
		args = []interface{}{searchQuery, likeQuery, pageToken, pageSize + 1}
	}

	rows, err := cd.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, "", fmt.Errorf("failed to search contents: %w", err)
	}
	defer rows.Close()

	var contents []Content
	for rows.Next() {
		var content Content
		var publishedAt, createdAt, updatedAt time.Time
		var description, language, url, platformName sql.NullString
		var durationSeconds sql.NullInt32
		var deletedAt sql.NullTime
		var maxSimilarity float64

		err := rows.Scan(
			&content.ID,
			&content.Title,
			&description,
			&language,
			&durationSeconds,
			&publishedAt,
			&content.ContentType,
			&createdAt,
			&updatedAt,
			&url,
			&platformName,
			&deletedAt,
			&maxSimilarity,
		)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan content row: %w", err)
		}

		content.Description = description.String
		content.Language = language.String
		content.ExternalURL = url.String
		content.PlatformName = platformName.String
		content.DurationSeconds = durationSeconds.Int32
		content.PublishedAt = publishedAt
		content.CreatedAt = createdAt
		content.UpdatedAt = updatedAt
		if deletedAt.Valid {
			content.DeletedAt = &deletedAt.Time
		}

		tags, err := cd.getContentTags(ctx, content.ID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to get content tags: %w", err)
		}
		content.Tags = tags

		contents = append(contents, content)
	}

	if err = rows.Err(); err != nil {
		return nil, "", fmt.Errorf("error iterating over search results: %w", err)
	}

	var nextPageToken string
	if len(contents) > int(pageSize) {
		contents = contents[:pageSize]
		nextPageToken = contents[len(contents)-1].ID
	}

	return contents, nextPageToken, nil
}

func (cd *ContentData) getContentTags(ctx context.Context, contentID string) ([]string, error) {
	query := `
		SELECT t.name 
		FROM tags t
		INNER JOIN content_tags ct ON t.id = ct.tag_id
		WHERE ct.content_id = $1
		ORDER BY t.name`

	rows, err := cd.db.QueryContext(ctx, query, contentID)
	if err != nil {
		return nil, fmt.Errorf("failed to query content tags: %w", err)
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tagName string
		if err := rows.Scan(&tagName); err != nil {
			return nil, fmt.Errorf("failed to scan tag name: %w", err)
		}
		tags = append(tags, tagName)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over tag rows: %w", err)
	}

	return tags, nil
}
