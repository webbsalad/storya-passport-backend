package v1

import (
	"github.com/webbsalad/storya-passport-backend/internal/config"
	passport_service "github.com/webbsalad/storya-passport-backend/internal/service/passport"
	"github.com/webbsalad/storya-passport-backend/repository/passport"
)

type Service struct {
	passportRepository passport.Repository
	config             config.Config
}

func NewService(passportRepository passport.Repository, config config.Config) passport_service.Service {
	return &Service{
		passportRepository: passportRepository,
		config:             config,
	}
}
