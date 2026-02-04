package interceptor

import (
	"context"
	"fmt"
	"gotemplate/pkg/lg"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewRecoveryInterceptor(log lg.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				log.CtxError(ctx, "panic recovered",
					lg.Any(lg.ErrorFieldName, fmt.Sprintf("panic: %+v", r)),
					lg.Str(lg.StackFieldName, stack),
				)
				err = status.Errorf(codes.Internal, "internal server error")
			}
		}()

		return handler(ctx, req)
	}
}
