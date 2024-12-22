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
	testEmailID, _  = model.EmailIDFromString("1cde4608-a972-4ee1-98b2-0dbdf651d0d2")

	testUser       = model.User{Name: testName, EmailId: testEmailID}
	testSession    = model.Session{UserID: testUserID, DeviceID: testDeviceID, Version: testVersion}
	testAuthTokens = model.AuthTokens{AccessToken: testAccessToken, RefreshToken: testRefreshToken}

	testName     = "name"
	testVersion  = 1
	testPassword = "1234567"
	testHash     = "$2a$10$wLvXkSKzSEqAy5Et66.g4OGe7reW8mQMjczE26Xkx7cmG.4MmZDVO"

	// testUserID, testDeviceID, testVersion, refresh
	testRefreshToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMmI5OGVlODgtNzk3MC00ZTZmLWIzMjUtY2NmM2NlMTA5MDlmIiwiZGV2aWNlX2lkIjoiMWNkZTQ2MDgtYTk3Mi00ZWUxLTk4YjItMGRiZGY2NTRkMGQyIiwidHlwZSI6InJlZnJlc2giLCJ2ZXJzaW9uIjoxLCJleHAiOjQ1MTYyMzkwMjJ9.f9RvxXBf2E2hfTjndIHetBOaj88-D32s-X6Nju-GWZY"
	testAccessToken  = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMmI5OGVlODgtNzk3MC00ZTZmLWIzMjUtY2NmM2NlMTA5MDlmIiwiZGV2aWNlX2lkIjoiMWNkZTQ2MDgtYTk3Mi00ZWUxLTk4YjItMGRiZGY2NTRkMGQyIiwidmVyc2lvbiI6MSwiZXhwIjo0NTE2MjM5MDIyfQ.AplijZv6-ZjoB8nNRAHsw36DjEK9ukdYz0JMwBgcl2s"
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

type hashMatcher struct {
	password string
}

func (hm hashMatcher) Matches(x interface{}) bool {
	hashedPassword, ok := x.(string)
	if !ok {
		return false
	}

	err := comparePassword(hashedPassword, hm.password)
	return err == nil
}

func (hm hashMatcher) String() string {
	return "matches hashed password"
}
