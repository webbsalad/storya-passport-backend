package passport

import (
	desc "github.com/webbsalad/storya-passport-backend/internal/pb/github.com/webbsalad/storya-passport-backend/passport"
	service "github.com/webbsalad/storya-passport-backend/internal/service/passport"
)

type Implementation struct {
	desc.UnimplementedPassportServiceServer

	PassportService service.Service
}

func NewImplementation(passportService service.Service) desc.PassportServiceServer {
	return &Implementation{
		PassportService: passportService,
	}
}
