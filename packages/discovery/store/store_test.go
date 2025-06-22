package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreInterface(t *testing.T) {
	var _ Interface = &ContentData{}
	assert.True(t, true, "ContentData successfully implements the Interface")
}

func TestGetContent_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()
	contentID := "550e8400-e29b-41d4-a716-446655440000"

	publishedAt := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	createdAt := time.Now()
	updatedAt := time.Now()

	contentRows := sqlmock.NewRows([]string{
		"id", "title", "description", "language", "duration_seconds",
		"published_at", "content_type", "created_at", "updated_at", "url", "platform_name",
	}).AddRow(
		contentID, "Test Content", "A test description", "en", 3600,
		publishedAt, "podcast", createdAt, updatedAt, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Test Platform",
	)

	mock.ExpectQuery(`SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name FROM contents WHERE id = \$1 AND deleted_at IS NULL`).
		WithArgs(contentID).
		WillReturnRows(contentRows)

	tagRows := sqlmock.NewRows([]string{"name"}).
		AddRow("technology").
		AddRow("podcast")

	mock.ExpectQuery(`SELECT t\.name FROM tags t INNER JOIN content_tags ct ON t\.id = ct\.tag_id WHERE ct\.content_id = \$1 ORDER BY t\.name`).
		WithArgs(contentID).
		WillReturnRows(tagRows)

	content, err := store.GetContent(ctx, contentID)

	require.NoError(t, err)
	require.NotNil(t, content)

	assert.Equal(t, contentID, content.ID)
	assert.Equal(t, "Test Content", content.Title)
	assert.Equal(t, "A test description", content.Description)
	assert.Equal(t, "en", content.Language)
	assert.Equal(t, int32(3600), content.DurationSeconds)
	assert.Equal(t, publishedAt, content.PublishedAt)
	assert.Equal(t, "podcast", content.ContentType)
	assert.Equal(t, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", content.ExternalURL)
	assert.Equal(t, "Test Platform", content.PlatformName)
	assert.Len(t, content.Tags, 2)
	assert.Contains(t, content.Tags, "technology")
	assert.Contains(t, content.Tags, "podcast")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetContent_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()
	contentID := "550e8400-e29b-41d4-a716-446655440999"

	mock.ExpectQuery(`SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name FROM contents WHERE id = \$1 AND deleted_at IS NULL`).
		WithArgs(contentID).
		WillReturnError(sql.ErrNoRows)

	content, err := store.GetContent(ctx, contentID)

	assert.Error(t, err)
	assert.Nil(t, content)
	assert.Contains(t, err.Error(), "not found")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListContents_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()

	createdAt1 := time.Now()
	createdAt2 := time.Now().Add(-time.Hour)

	contentRows := sqlmock.NewRows([]string{
		"id", "title", "description", "language", "duration_seconds",
		"published_at", "content_type", "created_at", "updated_at", "url", "platform_name",
	}).AddRow(
		"id1", "Content 1", "Description 1", "en", 1800,
		time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), "podcast", createdAt1, createdAt1, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Platform 1",
	).AddRow(
		"id2", "Content 2", "Description 2", "ar", 3600,
		time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC), "documentary", createdAt2, createdAt2, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Platform 2",
	)

	mock.ExpectQuery(`SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name FROM contents WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT \$1`).
		WithArgs(11).
		WillReturnRows(contentRows)

	tagRows1 := sqlmock.NewRows([]string{"name"}).AddRow("tech")
	mock.ExpectQuery(`SELECT t\.name FROM tags t INNER JOIN content_tags ct ON t\.id = ct\.tag_id WHERE ct\.content_id = \$1 ORDER BY t\.name`).
		WithArgs("id1").
		WillReturnRows(tagRows1)

	tagRows2 := sqlmock.NewRows([]string{"name"}).AddRow("science")
	mock.ExpectQuery(`SELECT t\.name FROM tags t INNER JOIN content_tags ct ON t\.id = ct\.tag_id WHERE ct\.content_id = \$1 ORDER BY t\.name`).
		WithArgs("id2").
		WillReturnRows(tagRows2)

	contents, nextPageToken, err := store.ListContents(ctx, 10, "")

	require.NoError(t, err)
	assert.Len(t, contents, 2)
	assert.Empty(t, nextPageToken)

	assert.Equal(t, "id1", contents[0].ID)
	assert.Equal(t, "Content 1", contents[0].Title)
	assert.Equal(t, "en", contents[0].Language)
	assert.Len(t, contents[0].Tags, 1)
	assert.Contains(t, contents[0].Tags, "tech")

	assert.Equal(t, "id2", contents[1].ID)
	assert.Equal(t, "Content 2", contents[1].Title)
	assert.Equal(t, "ar", contents[1].Language)
	assert.Len(t, contents[1].Tags, 1)
	assert.Contains(t, contents[1].Tags, "science")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchContents_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()
	searchQuery := "podcast"

	mock.ExpectExec(`SET SESSION pg_trgm\.similarity_threshold = 0\.10`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	createdAt := time.Now()

	searchRows := sqlmock.NewRows([]string{
		"id", "title", "description", "language", "duration_seconds",
		"published_at", "content_type", "created_at", "updated_at", "url", "platform_name", "max_similarity",
	}).AddRow(
		"search-id", "Found Podcast", "A podcast found by search", "en", 2700,
		time.Date(2024, 1, 15, 14, 0, 0, 0, time.UTC), "podcast", createdAt, createdAt, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Search Platform", 0.8,
	)

	mock.ExpectQuery(`WITH content_with_tags AS \(.*\) SELECT .* FROM content_with_tags WHERE .* ORDER BY max_similarity DESC, created_at DESC LIMIT \$3`).
		WithArgs(searchQuery, "%"+searchQuery+"%", 11).
		WillReturnRows(searchRows)

	tagRows := sqlmock.NewRows([]string{"name"}).AddRow("podcast").AddRow("search")
	mock.ExpectQuery(`SELECT t\.name FROM tags t INNER JOIN content_tags ct ON t\.id = ct\.tag_id WHERE ct\.content_id = \$1 ORDER BY t\.name`).
		WithArgs("search-id").
		WillReturnRows(tagRows)

	contents, nextPageToken, err := store.SearchContents(ctx, searchQuery, 10, "")

	require.NoError(t, err)
	assert.Len(t, contents, 1)
	assert.Empty(t, nextPageToken)

	content := contents[0]
	assert.Equal(t, "search-id", content.ID)
	assert.Equal(t, "Found Podcast", content.Title)
	assert.Equal(t, "A podcast found by search", content.Description)
	assert.Equal(t, "podcast", content.ContentType)
	assert.Len(t, content.Tags, 2)
	assert.Contains(t, content.Tags, "podcast")
	assert.Contains(t, content.Tags, "search")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSearchContents_EmptyQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()

	contents, nextPageToken, err := store.SearchContents(ctx, "", 10, "")

	require.NoError(t, err)
	assert.Empty(t, contents)
	assert.Empty(t, nextPageToken)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestListContents_WithPagination(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()

	createdAt := time.Now()

	contentRows := sqlmock.NewRows([]string{
		"id", "title", "description", "language", "duration_seconds",
		"published_at", "content_type", "created_at", "updated_at", "url", "platform_name",
	}).AddRow(
		"id1", "Content 1", "Description 1", "en", 1800,
		time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), "podcast", createdAt, createdAt, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Platform 1",
	).AddRow(
		"id2", "Content 2", "Description 2", "ar", 3600,
		time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC), "documentary", createdAt, createdAt, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Platform 2",
	).AddRow(
		"id3", "Content 3", "Description 3", "en", 900,
		time.Date(2024, 1, 17, 8, 0, 0, 0, time.UTC), "podcast", createdAt, createdAt, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Platform 3",
	)

	mock.ExpectQuery(`SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name FROM contents WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT \$1`).
		WithArgs(3).
		WillReturnRows(contentRows)

	tagRows1 := sqlmock.NewRows([]string{"name"}).AddRow("tag1")
	mock.ExpectQuery(`SELECT t\.name FROM tags t INNER JOIN content_tags ct ON t\.id = ct\.tag_id WHERE ct\.content_id = \$1 ORDER BY t\.name`).
		WithArgs("id1").
		WillReturnRows(tagRows1)

	tagRows2 := sqlmock.NewRows([]string{"name"}).AddRow("tag2")
	mock.ExpectQuery(`SELECT t\.name FROM tags t INNER JOIN content_tags ct ON t\.id = ct\.tag_id WHERE ct\.content_id = \$1 ORDER BY t\.name`).
		WithArgs("id2").
		WillReturnRows(tagRows2)

	tagRows3 := sqlmock.NewRows([]string{"name"}).AddRow("tag3")
	mock.ExpectQuery(`SELECT t\.name FROM tags t INNER JOIN content_tags ct ON t\.id = ct\.tag_id WHERE ct\.content_id = \$1 ORDER BY t\.name`).
		WithArgs("id3").
		WillReturnRows(tagRows3)

	contents, nextPageToken, err := store.ListContents(ctx, 2, "")

	require.NoError(t, err)
	assert.Len(t, contents, 2)
	assert.Equal(t, "id2", nextPageToken)

	assert.NoError(t, mock.ExpectationsWereMet())
}
