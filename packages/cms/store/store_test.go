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

func TestCreateContent_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()

	content := Content{
		Title:           "Test Podcast",
		Description:     "A test podcast description",
		Tags:            []string{"technology", "podcast"},
		Language:        "en",
		DurationSeconds: 3600,
		PublishedAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
		ContentType:     "podcast",
		ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
		PlatformName:    "Test Platform",
	}

	contentID := "550e8400-e29b-41d4-a716-446655440000"
	createdAt := time.Now()
	updatedAt := time.Now()

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO contents \(title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10\) RETURNING id, created_at, updated_at`).
		WithArgs(
			content.Title, content.Description, content.Language, content.DurationSeconds,
			content.PublishedAt, content.ContentType, sqlmock.AnyArg(), sqlmock.AnyArg(),
			content.ExternalURL, content.PlatformName,
		).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
			AddRow(contentID, createdAt, updatedAt))

	for _, tag := range content.Tags {
		tagID := "tag-id-" + tag
		mock.ExpectQuery(`INSERT INTO tags \(id, name\) VALUES \(gen_random_uuid\(\), \$1\) ON CONFLICT \(name\) DO UPDATE SET name = EXCLUDED\.name RETURNING id`).
			WithArgs(tag).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tagID))

		mock.ExpectExec(`INSERT INTO content_tags \(content_id, tag_id\) VALUES \(\$1, \$2\) ON CONFLICT \(content_id, tag_id\) DO NOTHING`).
			WithArgs(contentID, tagID).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	result, err := store.CreateContent(ctx, content)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, contentID, result.ID)
	assert.Equal(t, content.Title, result.Title)
	assert.Equal(t, content.Description, result.Description)
	assert.Equal(t, content.Language, result.Language)
	assert.Equal(t, content.DurationSeconds, result.DurationSeconds)
	assert.Equal(t, content.ContentType, result.ContentType)
	assert.Equal(t, content.Tags, result.Tags)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateContent_DatabaseError(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()

	content := Content{
		Title:       "Test Content",
		Description: "Test Description",
		Language:    "en",
		ContentType: "podcast",
	}

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO contents`).
		WillReturnError(sql.ErrConnDone)

	mock.ExpectRollback()

	result, err := store.CreateContent(ctx, content)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to insert content")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateContent_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()

	contentID := "550e8400-e29b-41d4-a716-446655440000"
	content := Content{
		ID:              contentID,
		Title:           "Updated Podcast",
		Description:     "Updated description",
		Tags:            []string{"updated", "technology"},
		Language:        "ar",
		DurationSeconds: 7200,
		PublishedAt:     time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC),
		ContentType:     "documentary",
		ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
		PlatformName:    "Updated Platform",
	}

	createdAt := time.Now().Add(-time.Hour)
	updatedAt := time.Now()

	mock.ExpectBegin()

	mock.ExpectQuery(`UPDATE contents SET title = \$1, description = \$2, language = \$3, duration_seconds = \$4, published_at = \$5, content_type = \$6, updated_at = \$7, url = \$8, platform_name = \$9 WHERE id = \$10 AND deleted_at IS NULL RETURNING created_at, updated_at`).
		WithArgs(
			content.Title, content.Description, content.Language, content.DurationSeconds,
			content.PublishedAt, content.ContentType, sqlmock.AnyArg(),
			content.ExternalURL, content.PlatformName, contentID,
		).
		WillReturnRows(sqlmock.NewRows([]string{"created_at", "updated_at"}).
			AddRow(createdAt, updatedAt))

	mock.ExpectExec(`DELETE FROM content_tags WHERE content_id = \$1`).
		WithArgs(contentID).
		WillReturnResult(sqlmock.NewResult(0, 2))

	for _, tag := range content.Tags {
		tagID := "tag-id-" + tag
		mock.ExpectQuery(`INSERT INTO tags \(id, name\) VALUES \(gen_random_uuid\(\), \$1\) ON CONFLICT \(name\) DO UPDATE SET name = EXCLUDED\.name RETURNING id`).
			WithArgs(tag).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(tagID))

		mock.ExpectExec(`INSERT INTO content_tags \(content_id, tag_id\) VALUES \(\$1, \$2\) ON CONFLICT \(content_id, tag_id\) DO NOTHING`).
			WithArgs(contentID, tagID).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	mock.ExpectCommit()

	result, err := store.UpdateContent(ctx, content)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, contentID, result.ID)
	assert.Equal(t, content.Title, result.Title)
	assert.Equal(t, content.Description, result.Description)
	assert.Equal(t, content.Tags, result.Tags)
	assert.Equal(t, createdAt.Format(time.RFC3339), result.CreatedAt.Format(time.RFC3339))

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateContent_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()

	contentID := "550e8400-e29b-41d4-a716-446655440999"
	content := Content{
		ID:          contentID,
		Title:       "Nonexistent Content",
		Description: "This content doesn't exist",
		Language:    "en",
		ContentType: "podcast",
	}

	mock.ExpectBegin()

	mock.ExpectQuery(`UPDATE contents SET`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			sqlmock.AnyArg(), contentID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectRollback()

	result, err := store.UpdateContent(ctx, content)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not found")

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteContent_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()
	contentID := "550e8400-e29b-41d4-a716-446655440000"

	mock.ExpectExec(`UPDATE contents SET deleted_at = \$1, updated_at = \$1 WHERE id = \$2 AND deleted_at IS NULL`).
		WithArgs(sqlmock.AnyArg(), contentID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = store.DeleteContent(ctx, contentID)

	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteContent_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	store := New(db)
	ctx := context.Background()
	contentID := "550e8400-e29b-41d4-a716-446655440999"

	mock.ExpectExec(`UPDATE contents SET deleted_at = \$1, updated_at = \$1 WHERE id = \$2 AND deleted_at IS NULL`).
		WithArgs(sqlmock.AnyArg(), contentID).
		WillReturnResult(sqlmock.NewResult(0, 0))

	err = store.DeleteContent(ctx, contentID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found or already deleted")

	assert.NoError(t, mock.ExpectationsWereMet())
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
		"published_at", "content_type", "created_at", "updated_at", "url", "platform_name", "deleted_at",
	}).AddRow(
		contentID, "Test Content", "A test description", "en", 3600,
		publishedAt, "podcast", createdAt, updatedAt, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Test Platform", nil,
	)

	mock.ExpectQuery(`SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name, deleted_at FROM contents WHERE id = \$1 AND deleted_at IS NULL`).
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
	assert.Nil(t, content.DeletedAt)

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
		"published_at", "content_type", "created_at", "updated_at", "url", "platform_name", "deleted_at",
	}).AddRow(
		"id1", "Content 1", "Description 1", "en", 1800,
		time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC), "podcast", createdAt1, createdAt1, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Platform 1", nil,
	).AddRow(
		"id2", "Content 2", "Description 2", "ar", 3600,
		time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC), "documentary", createdAt2, createdAt2, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", "Platform 2", nil,
	)

	mock.ExpectQuery(`SELECT id, title, description, language, duration_seconds, published_at, content_type, created_at, updated_at, url, platform_name, deleted_at FROM contents WHERE deleted_at IS NULL ORDER BY created_at DESC LIMIT \$1`).
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
