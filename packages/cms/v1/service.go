package v1

import (
	"context"
	// "fmt"
	"time"

	mawjoodv1 "mawjood/gen/go/packages/proto/v1"
	"github.com/mosaibah/Mawjood/packages/cms/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
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
	contentType := cs.protoContentTypeToString(req.ContentType)

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
	return cs.storeContentToProto(createdContent), nil
}

func (cs *CMSService) UpdateContent(ctx context.Context, req *mawjoodv1.UpdateContentRequest) (*mawjoodv1.Content, error) {
	// Validate required fields
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
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
	contentType := cs.protoContentTypeToString(req.ContentType)

	// Create content struct for store
	content := store.Content{
		ID:             req.Id,
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

	// Call store to update content
	updatedContent, err := cs.store.UpdateContent(ctx, content)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update content: %v", err)
	}

	// Convert back to proto response
	return cs.storeContentToProto(updatedContent), nil
}

func (cs *CMSService) DeleteContent(ctx context.Context, req *mawjoodv1.DeleteContentRequest) (*emptypb.Empty, error) {
	// Validate required fields
	if req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	// Call store to delete content
	err := cs.store.DeleteContent(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete content: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (cs *CMSService) ListContents(ctx context.Context, req *mawjoodv1.ListContentsRequest) (*mawjoodv1.ListContentsResponse, error) {
	// Call store to list contents
	contents, nextPageToken, err := cs.store.ListContents(ctx, req.PageSize, req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list contents: %v", err)
	}

	// Convert store contents to proto contents
	protoContents := make([]*mawjoodv1.Content, len(contents))
	for i, content := range contents {
		protoContents[i] = cs.storeContentToProto(&content)
	}

	return &mawjoodv1.ListContentsResponse{
		Contents:      protoContents,
		NextPageToken: nextPageToken,
	}, nil
}

func (cs *CMSService) ImportFromExternal(ctx context.Context, req *mawjoodv1.ImportRequest) (*mawjoodv1.ImportResponse, error) {
	// This is a stretch goal implementation - for now, return unimplemented
	return nil, status.Error(codes.Unimplemented, "ImportFromExternal is not yet implemented")
}

// Helper function to convert proto ContentType to string
func (cs *CMSService) protoContentTypeToString(contentType mawjoodv1.ContentType) string {
	switch contentType {
	case mawjoodv1.ContentType_CONTENT_TYPE_PODCAST:
		return "podcast"
	case mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY:
		return "documentary"
	default:
		return "podcast" // default fallback
	}
}

// Helper function to convert string to proto ContentType
func (cs *CMSService) stringToProtoContentType(contentType string) mawjoodv1.ContentType {
	switch contentType {
	case "podcast":
		return mawjoodv1.ContentType_CONTENT_TYPE_PODCAST
	case "documentary":
		return mawjoodv1.ContentType_CONTENT_TYPE_DOCUMENTARY
	default:
		return mawjoodv1.ContentType_CONTENT_TYPE_PODCAST // default fallback
	}
}

// Helper function to convert store.Content to proto Content
func (cs *CMSService) storeContentToProto(content *store.Content) *mawjoodv1.Content {
	var publishedAt string
	if !content.PublishedAt.IsZero() {
		publishedAt = content.PublishedAt.Format(time.RFC3339)
	}

	return &mawjoodv1.Content{
		Id:             content.ID,
		Title:          content.Title,
		Description:    content.Description,
		Tags:           content.Tags,
		Language:       content.Language,
		DurationSeconds: content.DurationSeconds,
		PublishedAt:    publishedAt,
		ContentType:    cs.stringToProtoContentType(content.ContentType),
		CreatedAt:      content.CreatedAt.Format(time.RFC3339),
		UpdatedAt:      content.UpdatedAt.Format(time.RFC3339),
		Url:           content.ExternalURL,
		PlatformName:   content.PlatformName,
	}
}
