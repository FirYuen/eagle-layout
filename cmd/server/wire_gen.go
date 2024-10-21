// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-eagle/eagle-layout/internal/dal"
	"github.com/go-eagle/eagle-layout/internal/dal/cache"
	"github.com/go-eagle/eagle-layout/internal/repository"
	"github.com/go-eagle/eagle-layout/internal/server"
	"github.com/go-eagle/eagle-layout/internal/service"
	"github.com/go-eagle/eagle/pkg/app"
	"github.com/go-eagle/eagle/pkg/log"
	"github.com/go-eagle/eagle/pkg/redis"
	"github.com/go-eagle/eagle/pkg/transport/grpc"
	"github.com/go-eagle/eagle/pkg/transport/http"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

func InitApp(cfg *app.Config) (*app.App, func(), error) {
	dbClient, cleanup, err := dal.Init()
	if err != nil {
		return nil, nil, err
	}
	client, cleanup2, err := redis.Init()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	userCache := cache.NewUserCache(client)
	userRepo := repository.NewUserRepo(dbClient, userCache)
	greeterServiceServer := service.NewGreeterServiceServer(userRepo)
	httpServer := server.NewHTTPServer(cfg, greeterServiceServer)
	grpcServer := server.NewGRPCServer(cfg)
	appApp := newApp(cfg, httpServer, grpcServer)
	return appApp, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

func newApp(cfg *app.Config, hs *http.Server, gs *grpc.Server) *app.App {
	return app.New(app.WithName(cfg.Name), app.WithVersion(cfg.Version), app.WithLogger(log.GetLogger()), app.WithServer(

		hs,

		gs,
	),
	)
}
