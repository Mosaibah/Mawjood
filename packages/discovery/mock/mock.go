package mock

import (
	"context"
	"fmt"
	"time"

	"github.com/mosaibah/Mawjood/packages/discovery/store"
)

type MockContentData struct{}

func (m *MockContentData) GetContent(ctx context.Context, id string) (*store.Content, error) {
	switch id {
	case "550e8400-e29b-41d4-a716-446655440000":
		return &store.Content{
			ID:              "550e8400-e29b-41d4-a716-446655440000",
			Title:           "Test Podcast",
			Description:     "A test podcast description",
			Tags:            []string{"test", "podcast"},
			Language:        "en",
			DurationSeconds: 3600,
			PublishedAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			ContentType:     "podcast",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
			PlatformName:    "Test Platform",
		}, nil
	case "550e8400-e29b-41d4-a716-446655440001":
		return &store.Content{
			ID:              "550e8400-e29b-41d4-a716-446655440001",
			Title:           "Test Documentary",
			Description:     "A test documentary description",
			Tags:            []string{"test", "documentary"},
			Language:        "ar",
			DurationSeconds: 7200,
			PublishedAt:     time.Date(2024, 1, 16, 12, 0, 0, 0, time.UTC),
			ContentType:     "documentary",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
			PlatformName:    "Documentary Platform",
		}, nil
	default:
		return nil, fmt.Errorf("content with ID %s not found", id)
	}
}

func (m *MockContentData) ListContents(ctx context.Context, pageSize int32, pageToken string) ([]store.Content, string, error) {
	// Return mock list of contents
	contents := []store.Content{
		{
			ID:              "550e8400-e29b-41d4-a716-446655440000",
			Title:           "Listed Podcast 1",
			Description:     "First podcast in the list",
			Tags:            []string{"podcast", "tech"},
			Language:        "en",
			DurationSeconds: 1800,
			PublishedAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			ContentType:     "podcast",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
			PlatformName:    "Podcast Platform",
		},
		{
			ID:              "550e8400-e29b-41d4-a716-446655440001",
			Title:           "Listed Documentary 1",
			Description:     "First documentary in the list",
			Tags:            []string{"documentary", "science"},
			Language:        "ar",
			DurationSeconds: 5400,
			PublishedAt:     time.Date(2024, 1, 16, 14, 0, 0, 0, time.UTC),
			ContentType:     "documentary",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
			PlatformName:    "Documentary Platform",
		},
	}

	// Simple pagination simulation
	if pageToken == "" {
		// First page
		if pageSize > 0 && int(pageSize) < len(contents) {
			return contents[:pageSize], "next-page-token", nil
		}
		return contents, "", nil
	}

	// Subsequent pages (return empty for simplicity)
	return []store.Content{}, "", nil
}

func (m *MockContentData) SearchContents(ctx context.Context, query string, pageSize int32, pageToken string) ([]store.Content, string, error) {
	// Return search results based on query
	if query == "podcast" {
		return []store.Content{
			{
				ID:              "550e8400-e29b-41d4-a716-446655440000",
				Title:           "Found Podcast",
				Description:     "A podcast that matches the search",
				Tags:            []string{"podcast", "search"},
				Language:        "en",
				DurationSeconds: 2700,
				PublishedAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
				ContentType:     "podcast",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
				PlatformName:    "Search Platform",
			},
		}, "", nil
	} else if query == "documentary" {
		return []store.Content{
			{
				ID:              "550e8400-e29b-41d4-a716-446655440001",
				Title:           "Found Documentary",
				Description:     "A documentary that matches the search",
				Tags:            []string{"documentary", "search"},
				Language:        "ar",
				DurationSeconds: 6300,
				PublishedAt:     time.Date(2024, 1, 16, 15, 0, 0, 0, time.UTC),
				ContentType:     "documentary",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
				PlatformName:    "Search Platform",
			},
		}, "", nil
	} else if query == "nonexistent" {
		// Return empty results for this query
		return []store.Content{}, "", nil
	}

	// Default: return mixed results
	return []store.Content{
		{
			ID:              "550e8400-e29b-41d4-a716-446655440000",
			Title:           "Mixed Search Result 1",
			Description:     "First result for mixed search",
			Tags:            []string{"mixed", "search"},
			Language:        "en",
			DurationSeconds: 1200,
			PublishedAt:     time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
			ContentType:     "podcast",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			ExternalURL:     "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
			PlatformName:    "Mixed Platform",
		},
	}, "", nil
}
