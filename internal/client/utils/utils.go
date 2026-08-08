package utils

import (
	"context"

	"github.com/artni96/GophKeeper/internal/client/config"
	"google.golang.org/grpc/metadata"
)

func PrepareMDContext(ctx context.Context, cfg *config.Config) context.Context {
	md := metadata.Pairs("authorization", cfg.Token)
	ctx = metadata.NewOutgoingContext(ctx, md)
	return ctx
}
