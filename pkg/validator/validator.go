package validator

import (
	feedv1 "github.com/VariableSan/gia-protos/gen/go/feed"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(s any) error {
	if err := validate.Struct(s); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		if len(validationErrors) > 0 {
			return status.Errorf(codes.InvalidArgument, "validation failed: %v", validationErrors[0])
		}
	}
	return nil
}

type CreateFeedRequestValidator struct {
	Title   string `validate:"required"`
	Content string `validate:"required"`
}

func ValidateCreateFeedRequest(req *feedv1.CreateFeedRequest) error {
	return Validate(CreateFeedRequestValidator{
		Title:   req.GetTitle(),
		Content: req.GetContent(),
	})
}
