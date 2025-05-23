package feed

import (
	"context"
	"fmt"

	"github.com/VariableSan/gia-feed/internal/domain/models"
	"github.com/VariableSan/gia-feed/internal/grpc/middleware"
	"github.com/VariableSan/gia-feed/pkg/validator"
	feedv1 "github.com/VariableSan/gia-protos/gen/go/feed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type FeedService interface {
	GetFeed(
		ctx context.Context,
		id string,
	) (*models.Feed, error)
	CreateFeed(
		ctx context.Context,
		title string,
		content string,
		authorID string,
	) (string, error)
}

type serverAPI struct {
	feedv1.UnimplementedFeedServiceServer
	feedService FeedService
}

func Register(gRPC *grpc.Server, feedService FeedService) {
	feedv1.RegisterFeedServiceServer(
		gRPC,
		&serverAPI{feedService: feedService},
	)
}

func (s *serverAPI) GetFeed(
	ctx context.Context,
	req *feedv1.GetFeedRequest,
) (*feedv1.GetFeedResponse, error) {
	if err := validator.ValidateGetFeedRequest(req); err != nil {
		return nil, err
	}

	feed, err := s.feedService.GetFeed(ctx, req.GetId())
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &feedv1.GetFeedResponse{
		Feed: &feedv1.Feed{
			Id:               feed.ID,
			AuthorId:         "",
			Title:            feed.Title,
			Content:          feed.Content,
			Status:           0,
			ModeratorComment: "",
			CreatedAt:        &timestamppb.Timestamp{},
			UpdatedAt:        &timestamppb.Timestamp{},
		},
	}, nil
}

func (s *serverAPI) CreateFeed(
	ctx context.Context,
	req *feedv1.CreateFeedRequest,
) (*feedv1.CreateFeedResponse, error) {
	userID, ok := ctx.Value(middleware.UserIDKey).(string)
	if !ok || userID == "" {
		fmt.Println("user ID not found")
		return nil, status.Error(codes.Unauthenticated, "internal error")
	}

	if err := validator.ValidateCreateFeedRequest(req); err != nil {
		return nil, err
	}

	id, err := s.feedService.CreateFeed(ctx, req.GetTitle(), req.GetContent(), userID)
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &feedv1.CreateFeedResponse{
		Id: id,
	}, nil
}
