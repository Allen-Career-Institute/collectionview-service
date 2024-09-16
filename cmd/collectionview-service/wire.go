//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"collectionview-service/internal/cache"
	"collectionview-service/internal/conf"
	"collectionview-service/internal/mongo"
	"collectionview-service/internal/server"
	"collectionview-service/internal/service"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// wireApp init kratos application.
func wireApp(*conf.Server, *conf.Data, log.Logger, *conf.Redis) (*kratos.App, func(), error) {
	panic(wire.Build(server.ProviderSet, mongo.ProviderSet, cache.ProviderSet, service.ProviderSet, newApp))
}
