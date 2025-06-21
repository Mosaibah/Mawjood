package v1

import (
	"context"
	"log"
	"time"

	mawjoodv1 "github.com/mosaibah/Mawjood/gen/go/packages/proto/v1"

	"github.com/mosaibah/Mawjood/packages/discovery/store"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DiscoveryService struct {
	mawjoodv1.UnimplementedDiscoveryServiceServer
	store store.Interface
}

func New(store store.Interface) *DiscoveryService {
	return &DiscoveryService{store: store}
}

func (ds *DiscoveryService) GetContent(ctx context.Context, req *mawjoodv1.GetContentRequest) (*mawjoodv1.Content, error) {
	log.Printf("GetContent started - ID: %s", req.Id)

	// Use protocol buffer validation
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Call store to get content
	content, err := ds.store.GetContent(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get content: %v", err)
	}

	log.Printf("GetContent completed successfully - ID: %s", content.ID)

	// Convert store content to proto response
	return ds.storeContentToProto(content), nil
}

func (ds *DiscoveryService) ListContents(ctx context.Context, req *mawjoodv1.ListContentsRequest) (*mawjoodv1.ListContentsResponse, error) {
	log.Printf("ListContents started")

	// Use protocol buffer validation
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Call store to list contents
	contents, nextPageToken, err := ds.store.ListContents(ctx, req.PageSize, req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list contents: %v", err)
	}

	// Convert store contents to proto contents
	protoContents := make([]*mawjoodv1.Content, len(contents))
	for i, content := range contents {
		protoContents[i] = ds.storeContentToProto(&content)
	}

	log.Printf("ListContents completed successfully - count: %d", len(contents))

	return &mawjoodv1.ListContentsResponse{
		Contents:      protoContents,
		NextPageToken: nextPageToken,
	}, nil
}

func (ds *DiscoveryService) SearchContents(ctx context.Context, req *mawjoodv1.SearchContentsRequest) (*mawjoodv1.SearchContentsResponse, error) {
	log.Printf("SearchContents started - query: %s", req.Query)

	// Use protocol buffer validation
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	// Call store to search contents
	contents, nextPageToken, err := ds.store.SearchContents(ctx, req.Query, req.PageSize, req.PageToken)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to search contents: %v", err)
	}

	// Convert store contents to proto contents
	protoContents := make([]*mawjoodv1.Content, len(contents))
	for i, content := range contents {
		protoContents[i] = ds.storeContentToProto(&content)
	}

	log.Printf("SearchContents completed successfully - count: %d", len(contents))

	return &mawjoodv1.SearchContentsResponse{
		Contents:      protoContents,
		NextPageToken: nextPageToken,
	}, nil
}

// Helper function to convert string to proto ContentType
func (ds *DiscoveryService) stringToProtoContentType(contentType string) mawjoodv1.ContentType {
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
func (ds *DiscoveryService) storeContentToProto(content *store.Content) *mawjoodv1.Content {
	var publishedAt string
	if !content.PublishedAt.IsZero() {
		publishedAt = content.PublishedAt.Format(time.RFC3339)
	}

	return &mawjoodv1.Content{
		Id:              content.ID,
		Title:           content.Title,
		Description:     content.Description,
		Tags:            content.Tags,
		Language:        content.Language,
		DurationSeconds: content.DurationSeconds,
		PublishedAt:     publishedAt,
		ContentType:     ds.stringToProtoContentType(content.ContentType),
		CreatedAt:       content.CreatedAt.Format(time.RFC3339),
		UpdatedAt:       content.UpdatedAt.Format(time.RFC3339),
		Url:             content.ExternalURL,
		PlatformName:    content.PlatformName,
	}
}
