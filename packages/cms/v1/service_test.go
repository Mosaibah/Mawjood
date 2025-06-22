package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		Url:             "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
		PlatformName:    "Test Platform",
	}

	resp, err := service.CreateContent(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", resp.Id)
	assert.Equal(t, "Test Podcast", resp.Title)
	assert.Equal(t, "A test podcast description", resp.Description)
	assert.Len(t, resp.Tags, 3)
	assert.Equal(t, "en", resp.Language)
	assert.Equal(t, int32(3600), resp.DurationSeconds)
	assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_PODCAST, resp.ContentType)
	assert.NotEmpty(t, resp.CreatedAt)
	assert.NotEmpty(t, resp.UpdatedAt)
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
		Url:             "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
		PlatformName:    "Updated Platform",
	}

	resp, err := service.UpdateContent(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", resp.Id)
	assert.Equal(t, "Updated Podcast", resp.Title)
	assert.Equal(t, "An updated podcast description", resp.Description)
	assert.Equal(t, "ar", resp.Language)
	assert.Equal(t, int32(7200), resp.DurationSeconds)
	assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY, resp.ContentType)
}

func TestDeleteContent(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.DeleteContentRequest{
		Id: "550e8400-e29b-41d4-a716-446655440000",
	}

	resp, err := service.DeleteContent(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestListContents(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.ListContentsRequest{
		PageSize:  10,
		PageToken: "",
	}

	resp, err := service.ListContents(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Len(t, resp.Contents, 2)

	if assert.NotEmpty(t, resp.Contents) {
		firstContent := resp.Contents[0]
		assert.Equal(t, "Test Content 1", firstContent.Title)
		assert.Equal(t, "en", firstContent.Language)
		assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_PODCAST, firstContent.ContentType)
	}

	if assert.Greater(t, len(resp.Contents), 1) {
		secondContent := resp.Contents[1]
		assert.Equal(t, "Test Content 2", secondContent.Title)
		assert.Equal(t, "ar", secondContent.Language)
		assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY, secondContent.ContentType)
	}

	assert.Empty(t, resp.NextPageToken)
}

func TestImportFromExternal(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.ImportRequest{
		Url: "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
	}

	resp, err := service.ImportFromExternal(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	statusErr, ok := status.FromError(err)
	require.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.Unimplemented, statusErr.Code())
	assert.Equal(t, "ImportFromExternal is not yet implemented", statusErr.Message())
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
		Url:             "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
		PlatformName:    "Test Platform",
	}

	resp, err := service.CreateContent(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	statusErr, ok := status.FromError(err)
	require.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
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
		Url:             "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG",
		PlatformName:    "Updated Platform",
	}

	resp, err := service.UpdateContent(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	statusErr, ok := status.FromError(err)
	require.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.InvalidArgument, statusErr.Code())
}
