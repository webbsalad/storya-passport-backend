package v1

import (
	"context"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/webbsalad/storya-passport-backend/internal/config"
	"github.com/webbsalad/storya-passport-backend/internal/model"
	mock_repository "github.com/webbsalad/storya-passport-backend/internal/repository/passport/mock"
)

var (
	testUserID, _   = model.UserIDFromString("2b98ee88-7970-4e6f-b325-ccf3ce10909f")
	testDeviceID, _ = model.DeviceIDFromString("1cde4608-a972-4ee1-98b2-0dbdf654d0d2")
	//wrongUserID, _   = model.UserIDFromString("00000000-000-0000-0000-0000000000000")
	//wrongDeviceID, _ = model.DeviceIDFromString("00000000-000-0000-0000-0000000000000")
	testEmailID, _ = model.EmailIDFromString("1cde4608-a972-4ee1-98b2-0dbdf651d0d2")

	testUser    = model.User{Name: "name", EmailId: testEmailID}
	testVersion = 1
	//testSecret       = "test secret"

	// testUserID, testDeviceID, testVersion, refresh
	testRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMmI5OGVlODgtNzk3MC00ZTZmLWIzMjUtY2NmM2NlMTA5MDlmIiwiZGV2aWNlX2lkIjoiMWNkZTQ2MDgtYTk3Mi00ZWUxLTk4YjItMGRiZGY2NTRkMGQyIiwidHlwZSI6InJlZnJlc2giLCJ2ZXJzaW9uIjoxLCJleHAiOjQ1MTYyMzkwMjJ9.f9RvxXBf2E2hfTjndIHetBOaj88-D32s-X6Nju-GWZY"
	testExpiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMmI5OGVlODgtNzk3MC00ZTZmLWIzMjUtY2NmM2NlMTA5MDlmIiwiZGV2aWNlX2lkIjoiMWNkZTQ2MDgtYTk3Mi00ZWUxLTk4YjItMGRiZGY2NTRkMGQyIiwidHlwZSI6InJlZnJlc2giLCJ2ZXJzaW9uIjoxLCJleHAiOjYyMzkwMjJ9.01hYQjFM6WAMeUegW7v6usUBkbgVGIJaco_TTsp-_ic"
)

type serviceTestDeps struct {
	Service *Service

	ctx                context.Context
	passportRepository *mock_repository.MockRepository
}

func getTestDeps(t *testing.T) *serviceTestDeps {
	ctrl := gomock.NewController(t)
	passportRepository := mock_repository.NewMockRepository(ctrl)

	return &serviceTestDeps{
		Service: &Service{
			passportRepository: passportRepository,
			config: config.Config{
				JWTSecret: "test secret",
			},
		},

		ctx:                context.Background(),
		passportRepository: passportRepository,
	}
}
