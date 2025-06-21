package v1

import (
	"context"
	"testing"

	mawjoodv1 "github.com/mosaibah/Mawjood/gen/go/packages/proto/v1"
	"github.com/mosaibah/Mawjood/packages/discovery/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetContent(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.GetContentRequest{
		Id: "550e8400-e29b-41d4-a716-446655440000",
	}

	resp, err := service.GetContent(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Verify the response matches our mock data
	if resp.Id != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected ID to be '550e8400-e29b-41d4-a716-446655440000', got '%s'", resp.Id)
	}

	if resp.Title != "Test Podcast" {
		t.Errorf("Expected title to be 'Test Podcast', got '%s'", resp.Title)
	}

	if resp.Description != "A test podcast description" {
		t.Errorf("Expected description to be 'A test podcast description', got '%s'", resp.Description)
	}

	if len(resp.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(resp.Tags))
	}

	if resp.Language != "en" {
		t.Errorf("Expected language to be 'en', got '%s'", resp.Language)
	}

	if resp.DurationSeconds != 3600 {
		t.Errorf("Expected duration to be 3600, got %d", resp.DurationSeconds)
	}

	if resp.ContentType != mawjoodv1.ContentType_CONTENT_TYPE_PODCAST {
		t.Errorf("Expected content type to be PODCAST, got %v", resp.ContentType)
	}

	if resp.PublishedAt != "2024-01-15T10:00:00Z" {
		t.Errorf("Expected published_at to be '2024-01-15T10:00:00Z', got '%s'", resp.PublishedAt)
	}

	if resp.Url != "https://example.com/podcast" {
		t.Errorf("Expected URL to be 'https://example.com/podcast', got '%s'", resp.Url)
	}

	if resp.PlatformName != "Test Platform" {
		t.Errorf("Expected platform name to be 'Test Platform', got '%s'", resp.PlatformName)
	}

	if resp.CreatedAt == "" {
		t.Error("Expected CreatedAt to not be empty")
	}

	if resp.UpdatedAt == "" {
		t.Error("Expected UpdatedAt to not be empty")
	}
}

func TestGetContent_DocumentaryType(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.GetContentRequest{
		Id: "550e8400-e29b-41d4-a716-446655440001",
	}

	resp, err := service.GetContent(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Verify documentary-specific fields
	if resp.Title != "Test Documentary" {
		t.Errorf("Expected title to be 'Test Documentary', got '%s'", resp.Title)
	}

	if resp.Language != "ar" {
		t.Errorf("Expected language to be 'ar', got '%s'", resp.Language)
	}

	if resp.ContentType != mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY {
		t.Errorf("Expected content type to be DOCUMENTARY, got %v", resp.ContentType)
	}

	if resp.DurationSeconds != 7200 {
		t.Errorf("Expected duration to be 7200, got %d", resp.DurationSeconds)
	}
}

func TestGetContent_NotFound(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	// Use a valid UUID format but one that doesn't exist in our mock
	req := &mawjoodv1.GetContentRequest{
		Id: "550e8400-e29b-41d4-a716-446655440999",
	}

	resp, err := service.GetContent(context.Background(), req)

	// Should return an error for non-existent content
	if err == nil {
		t.Fatal("Expected an error for non-existent content")
	}

	if resp != nil {
		t.Error("Expected response to be nil for non-existent content")
	}

	// Check that we got the correct error code
	statusErr, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if statusErr.Code() != codes.NotFound {
		t.Errorf("Expected NotFound error code, got %v", statusErr.Code())
	}
}

func TestListContents(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.ListContentsRequest{
		PageSize:  10,
		PageToken: "",
	}

	resp, err := service.ListContents(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Verify we got the expected number of contents from our mock
	if len(resp.Contents) != 2 {
		t.Errorf("Expected 2 contents, got %d", len(resp.Contents))
	}

	// Check first content
	if len(resp.Contents) > 0 {
		firstContent := resp.Contents[0]
		if firstContent.Title != "Listed Podcast 1" {
			t.Errorf("Expected first content title to be 'Listed Podcast 1', got '%s'", firstContent.Title)
		}

		if firstContent.Language != "en" {
			t.Errorf("Expected first content language to be 'en', got '%s'", firstContent.Language)
		}

		if firstContent.ContentType != mawjoodv1.ContentType_CONTENT_TYPE_PODCAST {
			t.Errorf("Expected first content type to be PODCAST, got %v", firstContent.ContentType)
		}

		if firstContent.DurationSeconds != 1800 {
			t.Errorf("Expected first content duration to be 1800, got %d", firstContent.DurationSeconds)
		}
	}

	// Check second content
	if len(resp.Contents) > 1 {
		secondContent := resp.Contents[1]
		if secondContent.Title != "Listed Documentary 1" {
			t.Errorf("Expected second content title to be 'Listed Documentary 1', got '%s'", secondContent.Title)
		}

		if secondContent.Language != "ar" {
			t.Errorf("Expected second content language to be 'ar', got '%s'", secondContent.Language)
		}

		if secondContent.ContentType != mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY {
			t.Errorf("Expected second content type to be DOCUMENTARY, got %v", secondContent.ContentType)
		}

		if secondContent.DurationSeconds != 5400 {
			t.Errorf("Expected second content duration to be 5400, got %d", secondContent.DurationSeconds)
		}
	}

	// Check pagination token (should be empty for our mock with default page size)
	if resp.NextPageToken != "" {
		t.Errorf("Expected empty next page token, got '%s'", resp.NextPageToken)
	}
}

func TestListContents_WithPagination(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.ListContentsRequest{
		PageSize:  1, // Request only 1 item to test pagination
		PageToken: "",
	}

	resp, err := service.ListContents(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Should return only 1 content due to page size
	if len(resp.Contents) != 1 {
		t.Errorf("Expected 1 content, got %d", len(resp.Contents))
	}

	// Should have a next page token
	if resp.NextPageToken != "next-page-token" {
		t.Errorf("Expected next page token to be 'next-page-token', got '%s'", resp.NextPageToken)
	}
}

func TestSearchContents_PodcastQuery(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.SearchContentsRequest{
		Query:     "podcast",
		PageSize:  10,
		PageToken: "",
	}

	resp, err := service.SearchContents(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Should return 1 podcast result
	if len(resp.Contents) != 1 {
		t.Errorf("Expected 1 content, got %d", len(resp.Contents))
	}

	if len(resp.Contents) > 0 {
		content := resp.Contents[0]
		if content.Title != "Found Podcast" {
			t.Errorf("Expected title to be 'Found Podcast', got '%s'", content.Title)
		}

		if content.ContentType != mawjoodv1.ContentType_CONTENT_TYPE_PODCAST {
			t.Errorf("Expected content type to be PODCAST, got %v", content.ContentType)
		}

		if content.Language != "en" {
			t.Errorf("Expected language to be 'en', got '%s'", content.Language)
		}
	}
}

func TestSearchContents_DocumentaryQuery(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.SearchContentsRequest{
		Query:     "documentary",
		PageSize:  10,
		PageToken: "",
	}

	resp, err := service.SearchContents(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Should return 1 documentary result
	if len(resp.Contents) != 1 {
		t.Errorf("Expected 1 content, got %d", len(resp.Contents))
	}

	if len(resp.Contents) > 0 {
		content := resp.Contents[0]
		if content.Title != "Found Documentary" {
			t.Errorf("Expected title to be 'Found Documentary', got '%s'", content.Title)
		}

		if content.ContentType != mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY {
			t.Errorf("Expected content type to be DOCUMENTARY, got %v", content.ContentType)
		}

		if content.Language != "ar" {
			t.Errorf("Expected language to be 'ar', got '%s'", content.Language)
		}
	}
}

func TestSearchContents_NoResults(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.SearchContentsRequest{
		Query:     "nonexistent",
		PageSize:  10,
		PageToken: "",
	}

	resp, err := service.SearchContents(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Should return no results
	if len(resp.Contents) != 0 {
		t.Errorf("Expected 0 contents, got %d", len(resp.Contents))
	}

	// Should have empty pagination token
	if resp.NextPageToken != "" {
		t.Errorf("Expected empty next page token, got '%s'", resp.NextPageToken)
	}
}

func TestSearchContents_MixedResults(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.SearchContentsRequest{
		Query:     "mixed",
		PageSize:  10,
		PageToken: "",
	}

	resp, err := service.SearchContents(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Should return 1 mixed result
	if len(resp.Contents) != 1 {
		t.Errorf("Expected 1 content, got %d", len(resp.Contents))
	}

	if len(resp.Contents) > 0 {
		content := resp.Contents[0]
		if content.Title != "Mixed Search Result 1" {
			t.Errorf("Expected title to be 'Mixed Search Result 1', got '%s'", content.Title)
		}

		if content.PlatformName != "Mixed Platform" {
			t.Errorf("Expected platform name to be 'Mixed Platform', got '%s'", content.PlatformName)
		}
	}
}
