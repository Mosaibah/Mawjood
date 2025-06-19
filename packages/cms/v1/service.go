package v1

import (
	"context"
	// "fmt"
	"time"

	mawjoodv1 "mawjood/gen/go/packages/proto/v1"
	"github.com/mosaibah/Mawjood/packages/cms/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CMSService struct {
	mawjoodv1.UnimplementedCMSServiceServer
	store store.Interface
}

func New(store store.Interface) *CMSService {
	return &CMSService{store: store}
}

func (cs *CMSService) CreateContent(ctx context.Context, req *mawjoodv1.CreateContentRequest) (*mawjoodv1.Content, error) {
	// Validate required fields
	if req.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if req.ContentType == mawjoodv1.ContentType_CONTENT_TYPE_UNSPECIFIED {
		return nil, status.Error(codes.InvalidArgument, "content_type is required")
	}

	// Parse published_at if provided
	var publishedAt time.Time
	var err error
	if req.PublishedAt != "" {
		publishedAt, err = time.Parse(time.RFC3339, req.PublishedAt)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid published_at format: %v", err)
		}
	}

	// Convert content type to string
	contentType := "podcast"
	if req.ContentType == mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY {
		contentType = "documentary"
	}

	// Create content struct for store
	content := store.Content{
		Title:          req.Title,
		Description:    req.Description,
		Tags:           req.Tags,
		Language:       req.Language,
		DurationSeconds: req.DurationSeconds,
		PublishedAt:    publishedAt,
		ContentType:    contentType,
		ExternalURL:    req.Url,
		PlatformName:   req.PlatformName,
	}

	// Call store to create content
	createdContent, err := cs.store.CreateContent(ctx, content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create content: %v", err)
	}

	// Convert back to proto response
	response := &mawjoodv1.Content{
		Id:             createdContent.ID,
		Title:          createdContent.Title,
		Description:    createdContent.Description,
		Tags:           createdContent.Tags,
		Language:       createdContent.Language,
		DurationSeconds: createdContent.DurationSeconds,
		PublishedAt:    createdContent.PublishedAt.Format(time.RFC3339),
		ContentType:    req.ContentType, // Convert back to proto enum
		CreatedAt:      createdContent.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      createdContent.UpdatedAt.Format(time.RFC3339),
		Url:    createdContent.ExternalURL,
		PlatformName:   createdContent.PlatformName,
	}

	return response, nil
}
