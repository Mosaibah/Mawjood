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
	GetContent(ctx context.Context, id string) (*Content, error)
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
}

func (cd *ContentData) GetContent(ctx context.Context, id string) (*Content, error) {
	// Get the content details (only non-deleted content)
	getContentQuery := `
		SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name
		FROM contents 
		WHERE id = $1 AND deleted_at IS NULL`

	var content Content
	var publishedAt, createdAt, updatedAt time.Time
	var description, language, url, platformName sql.NullString
	var durationSeconds sql.NullInt32

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
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("content with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get content: %w", err)
	}

	// Handle nullable fields
	content.Description = description.String
	content.Language = language.String
	content.ExternalURL = url.String
	content.PlatformName = platformName.String
	content.DurationSeconds = durationSeconds.Int32
	content.PublishedAt = publishedAt
	content.CreatedAt = createdAt
	content.UpdatedAt = updatedAt

	// Get associated tags
	tags, err := cd.getContentTags(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get content tags: %w", err)
	}
	content.Tags = tags

	return &content, nil
}

func (cd *ContentData) ListContents(ctx context.Context, pageSize int32, pageToken string) ([]Content, string, error) {
	// Default page size if not specified
	if pageSize <= 0 {
		pageSize = 10
	}

	// Maximum page size limit
	if pageSize > 100 {
		pageSize = 100
	}

	// Build the query with pagination
	var query string
	var args []interface{}

	if pageToken == "" {
		// First page (only non-deleted content)
		query = `
			SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name
			FROM contents 
			WHERE deleted_at IS NULL
			ORDER BY created_at DESC 
			LIMIT $1`
		args = []interface{}{pageSize + 1} // Get one extra to determine if there's a next page
	} else {
		// Subsequent pages - use cursor-based pagination (only non-deleted content)
		query = `
			SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name
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
		)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan content row: %w", err)
		}

		// Handle nullable fields
		content.Description = description.String
		content.Language = language.String
		content.ExternalURL = url.String
		content.PlatformName = platformName.String
		content.DurationSeconds = durationSeconds.Int32
		content.PublishedAt = publishedAt
		content.CreatedAt = createdAt
		content.UpdatedAt = updatedAt

		// Get associated tags
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

	// Determine next page token
	var nextPageToken string
	if len(contents) > int(pageSize) {
		// Remove the extra record and use the last record's ID as the next page token
		contents = contents[:pageSize]
		nextPageToken = contents[len(contents)-1].ID
	}

	return contents, nextPageToken, nil
}

func (cd *ContentData) SearchContents(ctx context.Context, query string, pageSize int32, pageToken string) ([]Content, string, error) {
	// Default page size if not specified
	if pageSize <= 0 {
		pageSize = 10
	}

	// Maximum page size limit
	if pageSize > 100 {
		pageSize = 100
	}

	// Sanitize the search query for full-text search
	searchQuery := strings.TrimSpace(query)
	if searchQuery == "" {
		return []Content{}, "", nil
	}

	// Build the search query with pagination
	var sqlQuery string
	var args []interface{}

	if pageToken == "" {
		// First page - search in title and description using trigram similarity (only non-deleted content)
		sqlQuery = `
			SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name
			FROM contents 
			WHERE deleted_at IS NULL AND (title ILIKE $1 OR description ILIKE $1)
			ORDER BY 
				CASE 
					WHEN title ILIKE $1 THEN 1
					WHEN description ILIKE $1 THEN 2
					ELSE 3
				END,
				created_at DESC 
			LIMIT $2`
		likeQuery := "%" + searchQuery + "%"
		args = []interface{}{likeQuery, pageSize + 1}
	} else {
		// Subsequent pages with cursor-based pagination (only non-deleted content)
		sqlQuery = `
			SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name
			FROM contents 
			WHERE deleted_at IS NULL AND (title ILIKE $1 OR description ILIKE $1) 
			AND created_at < (SELECT created_at FROM contents WHERE id = $2)
			ORDER BY 
				CASE 
					WHEN title ILIKE $1 THEN 1
					WHEN description ILIKE $1 THEN 2
					ELSE 3
				END,
				created_at DESC 
			LIMIT $3`
		likeQuery := "%" + searchQuery + "%"
		args = []interface{}{likeQuery, pageToken, pageSize + 1}
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
		)
		if err != nil {
			return nil, "", fmt.Errorf("failed to scan content row: %w", err)
		}

		// Handle nullable fields
		content.Description = description.String
		content.Language = language.String
		content.ExternalURL = url.String
		content.PlatformName = platformName.String
		content.DurationSeconds = durationSeconds.Int32
		content.PublishedAt = publishedAt
		content.CreatedAt = createdAt
		content.UpdatedAt = updatedAt

		// Get associated tags
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

	// Determine next page token
	var nextPageToken string
	if len(contents) > int(pageSize) {
		// Remove the extra record and use the last record's ID as the next page token
		contents = contents[:pageSize]
		nextPageToken = contents[len(contents)-1].ID
	}

	return contents, nextPageToken, nil
}

// Helper function to get tags for a specific content
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
