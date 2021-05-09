package Interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/paypay3/tukecholl-api/pkg/apperrors"
)

func TransmitErrorWithStatus() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		appErr := apperrors.AsAppError(err)
		if appErr == nil {
			log.Print("failed to type assertion for appError")
			return nil, status.Error(codes.Internal, "予期しないエラーが発生しました")
		}

		st, err := appErr.Status()
		if err != nil {
			log.Printf("failed to convert to gRPC status: %+v", err)
			return nil, status.Error(codes.Internal, "予期しないエラーが発生しました")
		}

		return nil, st.Err()
	}
}
