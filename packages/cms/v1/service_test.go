package v1

import (
	"context"
	"testing"

	mawjoodv1 "github.com/mosaibah/Mawjood/gen/go/packages/proto/v1"
	"github.com/mosaibah/Mawjood/packages/cms/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateContent(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.CreateContentRequest{
		Title:           "Test Podcast",
		Description:     "A test podcast description",
		Tags:            []string{"test", "podcast", "technology"},
		Language:        "en",
		DurationSeconds: 3600,
		PublishedAt:     "2024-01-15T10:00:00Z",
		ContentType:     mawjoodv1.ContentType_CONTENT_TYPE_PODCAST,
		Url:             "https://example.com/podcast",
		PlatformName:    "Test Platform",
	}

	resp, err := service.CreateContent(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Verify the response matches our request
	if resp.Id != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected ID to be '550e8400-e29b-41d4-a716-446655440000', got '%s'", resp.Id)
	}

	if resp.Title != "Test Podcast" {
		t.Errorf("Expected title to be 'Test Podcast', got '%s'", resp.Title)
	}

	if resp.Description != "A test podcast description" {
		t.Errorf("Expected description to be 'A test podcast description', got '%s'", resp.Description)
	}

	if len(resp.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(resp.Tags))
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

	if resp.CreatedAt == "" {
		t.Error("Expected CreatedAt to not be empty")
	}

	if resp.UpdatedAt == "" {
		t.Error("Expected UpdatedAt to not be empty")
	}
}

func TestUpdateContent(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.UpdateContentRequest{
		Id:              "550e8400-e29b-41d4-a716-446655440000",
		Title:           "Updated Podcast",
		Description:     "An updated podcast description",
		Tags:            []string{"updated", "podcast", "technology"},
		Language:        "ar",
		DurationSeconds: 7200,
		PublishedAt:     "2024-01-16T12:00:00Z",
		ContentType:     mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY,
		Url:             "https://example.com/updated-podcast",
		PlatformName:    "Updated Platform",
	}

	resp, err := service.UpdateContent(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// Verify the response matches our request
	if resp.Id != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected ID to be '550e8400-e29b-41d4-a716-446655440000', got '%s'", resp.Id)
	}

	if resp.Title != "Updated Podcast" {
		t.Errorf("Expected title to be 'Updated Podcast', got '%s'", resp.Title)
	}

	if resp.Description != "An updated podcast description" {
		t.Errorf("Expected description to be 'An updated podcast description', got '%s'", resp.Description)
	}

	if resp.Language != "ar" {
		t.Errorf("Expected language to be 'ar', got '%s'", resp.Language)
	}

	if resp.DurationSeconds != 7200 {
		t.Errorf("Expected duration to be 7200, got %d", resp.DurationSeconds)
	}

	if resp.ContentType != mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY {
		t.Errorf("Expected content type to be DOCUMENTARY, got %v", resp.ContentType)
	}
}

func TestDeleteContent(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.DeleteContentRequest{
		Id: "550e8400-e29b-41d4-a716-446655440000",
	}

	resp, err := service.DeleteContent(context.Background(), req)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response to not be nil")
	}

	// DeleteContent should return an empty response on success
	// The fact that we got here without error means the deletion was successful
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
		if firstContent.Title != "Test Content 1" {
			t.Errorf("Expected first content title to be 'Test Content 1', got '%s'", firstContent.Title)
		}

		if firstContent.Language != "en" {
			t.Errorf("Expected first content language to be 'en', got '%s'", firstContent.Language)
		}

		if firstContent.ContentType != mawjoodv1.ContentType_CONTENT_TYPE_PODCAST {
			t.Errorf("Expected first content type to be PODCAST, got %v", firstContent.ContentType)
		}
	}

	// Check second content
	if len(resp.Contents) > 1 {
		secondContent := resp.Contents[1]
		if secondContent.Title != "Test Content 2" {
			t.Errorf("Expected second content title to be 'Test Content 2', got '%s'", secondContent.Title)
		}

		if secondContent.Language != "ar" {
			t.Errorf("Expected second content language to be 'ar', got '%s'", secondContent.Language)
		}

		if secondContent.ContentType != mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY {
			t.Errorf("Expected second content type to be DOCUMENTARY, got %v", secondContent.ContentType)
		}
	}

	// Check pagination token (should be empty for our mock)
	if resp.NextPageToken != "" {
		t.Errorf("Expected empty next page token, got '%s'", resp.NextPageToken)
	}
}

func TestImportFromExternal(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.ImportRequest{
		Url: "https://www.youtube.com/watch?v=example",
	}

	resp, err := service.ImportFromExternal(context.Background(), req)

	// This function is not implemented yet, so we expect an Unimplemented error
	if err == nil {
		t.Fatal("Expected an error for unimplemented function")
	}

	if resp != nil {
		t.Error("Expected response to be nil for unimplemented function")
	}

	// Check that we got the correct error code
	statusErr, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if statusErr.Code() != codes.Unimplemented {
		t.Errorf("Expected Unimplemented error code, got %v", statusErr.Code())
	}

	if statusErr.Message() != "ImportFromExternal is not yet implemented" {
		t.Errorf("Expected specific error message, got '%s'", statusErr.Message())
	}
}

func TestCreateContent_InvalidPublishedAt(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.CreateContentRequest{
		Title:           "Test Podcast",
		Description:     "A test podcast description",
		Tags:            []string{"test"},
		Language:        "en",
		DurationSeconds: 3600,
		PublishedAt:     "invalid-date-format",
		ContentType:     mawjoodv1.ContentType_CONTENT_TYPE_PODCAST,
		Url:             "https://example.com/podcast",
		PlatformName:    "Test Platform",
	}

	resp, err := service.CreateContent(context.Background(), req)

	// Should return an error for invalid date format
	if err == nil {
		t.Fatal("Expected an error for invalid date format")
	}

	if resp != nil {
		t.Error("Expected response to be nil for invalid request")
	}

	// Check that we got the correct error code
	statusErr, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if statusErr.Code() != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument error code, got %v", statusErr.Code())
	}
}

func TestUpdateContent_InvalidPublishedAt(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.UpdateContentRequest{
		Id:              "550e8400-e29b-41d4-a716-446655440000",
		Title:           "Updated Podcast",
		Description:     "An updated podcast description",
		Tags:            []string{"updated"},
		Language:        "ar",
		DurationSeconds: 7200,
		PublishedAt:     "invalid-date-format",
		ContentType:     mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY,
		Url:             "https://example.com/updated-podcast",
		PlatformName:    "Updated Platform",
	}

	resp, err := service.UpdateContent(context.Background(), req)

	// Should return an error for invalid date format
	if err == nil {
		t.Fatal("Expected an error for invalid date format")
	}

	if resp != nil {
		t.Error("Expected response to be nil for invalid request")
	}

	// Check that we got the correct error code
	statusErr, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if statusErr.Code() != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument error code, got %v", statusErr.Code())
	}
}
