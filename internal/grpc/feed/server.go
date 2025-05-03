package feed

import (
	"context"
	"fmt"

	"github.com/VariableSan/gia-feed/pkg/validator"
	feedv1 "github.com/VariableSan/gia-protos/gen/go/feed"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Feed interface {
	SaveFeed(
		ctx context.Context,
		title string,
		content string,
	) (int64, error)
}

type serverAPI struct {
	feedv1.UnimplementedFeedServiceServer
	feed Feed
}

func Register(gRPC *grpc.Server, feed Feed) {
	feedv1.RegisterFeedServiceServer(
		gRPC,
		&serverAPI{feed: feed},
	)
}

func (s *serverAPI) CreateFeed(
	ctx context.Context,
	req *feedv1.CreateFeedRequest,
) (*feedv1.CreateFeedResponse, error) {
	if err := validator.ValidateCreateFeedRequest(req); err != nil {
		return nil, err
	}

	_, err := s.feed.SaveFeed(ctx, req.GetTitle(), req.GetContent())
	if err != nil {
		fmt.Println(err)
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &feedv1.CreateFeedResponse{}, nil
}
