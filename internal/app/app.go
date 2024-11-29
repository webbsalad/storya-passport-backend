package app

import (
	"github.com/webbsalad/storya-passport-backend/internal/api/passport"
	"github.com/webbsalad/storya-passport-backend/internal/config"
	pb "github.com/webbsalad/storya-passport-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
	"github.com/webbsalad/storya-passport-backend/internal/repository/passport/pg"

	v1 "github.com/webbsalad/storya-passport-backend/internal/service/passport/v1"
	"go.uber.org/fx"
)

func NewApp() *fx.App {
	return fx.New(
		fx.Provide(
			config.NewConfig,
			initDB,
		),
		grpcOptions(),
		serviceOptions(),
	)
}

func serviceOptions() fx.Option {
	return fx.Options(
		fx.Provide(
			passport.NewImplementation,
			pg.NewRepository,
			v1.NewService,
		),
		fx.Invoke(
			pb.RegisterPassportServiceServer,
		),
	)
}
