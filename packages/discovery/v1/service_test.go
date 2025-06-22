package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "550e8400-e29b-41d4-a716-446655440000", resp.Id)
	assert.Equal(t, "Test Podcast", resp.Title)
	assert.Equal(t, "A test podcast description", resp.Description)
	assert.Len(t, resp.Tags, 2)
	assert.Equal(t, "en", resp.Language)
	assert.Equal(t, int32(3600), resp.DurationSeconds)
	assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_PODCAST, resp.ContentType)
	assert.Equal(t, "2024-01-15T10:00:00Z", resp.PublishedAt)
	assert.Equal(t, "https://youtu.be/mcrAH6g7CFk?si=vMHT2MSD6kAPlguG", resp.Url)
	assert.Equal(t, "Test Platform", resp.PlatformName)
	assert.NotEmpty(t, resp.CreatedAt)
	assert.NotEmpty(t, resp.UpdatedAt)
}

func TestGetContent_DocumentaryType(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.GetContentRequest{
		Id: "550e8400-e29b-41d4-a716-446655440001",
	}

	resp, err := service.GetContent(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, "Test Documentary", resp.Title)
	assert.Equal(t, "ar", resp.Language)
	assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY, resp.ContentType)
	assert.Equal(t, int32(7200), resp.DurationSeconds)
}

func TestGetContent_NotFound(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.GetContentRequest{
		Id: "550e8400-e29b-41d4-a716-446655440999",
	}

	resp, err := service.GetContent(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)

	statusErr, ok := status.FromError(err)
	require.True(t, ok, "Expected gRPC status error")
	assert.Equal(t, codes.NotFound, statusErr.Code())
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
		assert.Equal(t, "Listed Podcast 1", firstContent.Title)
		assert.Equal(t, "en", firstContent.Language)
		assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_PODCAST, firstContent.ContentType)
		assert.Equal(t, int32(1800), firstContent.DurationSeconds)
	}

	if assert.Greater(t, len(resp.Contents), 1) {
		secondContent := resp.Contents[1]
		assert.Equal(t, "Listed Documentary 1", secondContent.Title)
		assert.Equal(t, "ar", secondContent.Language)
		assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY, secondContent.ContentType)
		assert.Equal(t, int32(5400), secondContent.DurationSeconds)
	}

	assert.Empty(t, resp.NextPageToken)
}

func TestListContents_WithPagination(t *testing.T) {
	mockStore := &mock.MockContentData{}
	service := New(mockStore)

	req := &mawjoodv1.ListContentsRequest{
		PageSize:  1,
		PageToken: "",
	}

	resp, err := service.ListContents(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Len(t, resp.Contents, 1)
	assert.Equal(t, "next-page-token", resp.NextPageToken)
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

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Len(t, resp.Contents, 1)

	if assert.NotEmpty(t, resp.Contents) {
		content := resp.Contents[0]
		assert.Equal(t, "Found Podcast", content.Title)
		assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_PODCAST, content.ContentType)
		assert.Equal(t, "en", content.Language)
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

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Len(t, resp.Contents, 1)

	if assert.NotEmpty(t, resp.Contents) {
		content := resp.Contents[0]
		assert.Equal(t, "Found Documentary", content.Title)
		assert.Equal(t, mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY, content.ContentType)
		assert.Equal(t, "ar", content.Language)
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

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Empty(t, resp.Contents)
	assert.Empty(t, resp.NextPageToken)
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

	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Len(t, resp.Contents, 1)

	if assert.NotEmpty(t, resp.Contents) {
		content := resp.Contents[0]
		assert.Equal(t, "Mixed Search Result 1", content.Title)
		assert.Equal(t, "Mixed Platform", content.PlatformName)
	}
}
