package absences

import (
	"context"
	proteiservice1 "github.com/cockyafspb/contracts/gen/go/proteiservice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Absences interface {
	GetUser(ctx context.Context, email string) (fullName string, err error)
}

type serverAPI struct {
	proteiservice1.UnimplementedAbsencesServer
	absences Absences
}

func Register(gRPC *grpc.Server, absences Absences) {
	proteiservice1.RegisterAbsencesServer(gRPC, &serverAPI{absences: absences})
}

func (s *serverAPI) GetUser(
	ctx context.Context,
	req *proteiservice1.GetUserRequest,
) (*proteiservice1.GetUserResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	fullName, err := s.absences.GetUser(ctx, req.GetEmail())
	if err != nil {
		// TODO: ...
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &proteiservice1.GetUserResponse{
		Ok:       true,
		FullName: fullName,
	}, nil
}
