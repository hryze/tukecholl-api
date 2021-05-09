package Interceptor

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"github.com/paypay3/tukecholl-api/account/config"
	"github.com/paypay3/tukecholl-api/pkg/apperrors"
)

func Logging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		appErr := apperrors.AsAppError(err)
		if appErr == nil {
			log.Print("failed to type assertion for appError")
			return nil, err
		}

		// Output log if Debug flag is true in development environment.
		if config.Env.Log.Debug {
			log.Printf("%+v", appErr)
			return nil, err
		}

		// Transfer log of ERROR level or higher in production environment.
		if appErr.IsLevelError() || appErr.IsLevelCritical() {
			log.Printf("%+v", appErr)
			return nil, err
		}

		return nil, err
	}
}
