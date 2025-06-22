package mock

import (
	"context"
	"time"

	"github.com/mosaibah/Mawjood/packages/cms/store"
)

type MockContentData struct{}

func (m *MockContentData) CreateContent(ctx context.Context, content store.Content) (*store.Content, error) {
	content.ID = "550e8400-e29b-41d4-a716-446655440000"
	content.CreatedAt = time.Now()
	content.UpdatedAt = time.Now()
	return &content, nil
}

func (m *MockContentData) GetContent(ctx context.Context, id string) (*store.Content, error) {
	return &store.Content{
		ID:              "550e8400-e29b-41d4-a716-446655440000",
		Title:           "Test Content",
		Description:     "Test Description",
		Tags:            []string{"test", "mock"},
		Language:        "en",
		DurationSeconds: 3600,
		PublishedAt:     time.Now(),
		ContentType:     "podcast",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
		PlatformName:    "Test Platform",
	}, nil
}

func (m *MockContentData) UpdateContent(ctx context.Context, content store.Content) (*store.Content, error) {
	content.UpdatedAt = time.Now()
	return &content, nil
}

func (m *MockContentData) DeleteContent(ctx context.Context, id string) error {
	return nil
}

func (m *MockContentData) ListContents(ctx context.Context, pageSize int32, pageToken string) ([]store.Content, string, error) {
	return []store.Content{
		{
			ID:              "550e8400-e29b-41d4-a716-446655440000",
			Title:           "Test Content 1",
			Description:     "Test Description 1",
			Tags:            []string{"test", "mock"},
			Language:        "en",
			DurationSeconds: 3600,
			PublishedAt:     time.Now(),
			ContentType:     "podcast",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
			PlatformName:    "Test Platform",
		},
		{
			ID:              "550e8400-e29b-41d4-a716-446655440001",
			Title:           "Test Content 2",
			Description:     "Test Description 2",
			Tags:            []string{"test", "documentary"},
			Language:        "ar",
			DurationSeconds: 7200,
			PublishedAt:     time.Now(),
			ContentType:     "documentary",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
			PlatformName:    "Test Platform",
		},
	}, "", nil
}

func (m *MockContentData) SearchContents(ctx context.Context, query string, pageSize int32, pageToken string) ([]store.Content, string, error) {
	return []store.Content{
		{
			ID:              "550e8400-e29b-41d4-a716-446655440000",
			Title:           "Matching Content",
			Description:     "Content that matches the search query",
			Tags:            []string{"search", "test"},
			Language:        "en",
			DurationSeconds: 1800,
			PublishedAt:     time.Now(),
			ContentType:     "podcast",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguGh",
			PlatformName:    "Search Platform",
		},
	}, "", nil
}
