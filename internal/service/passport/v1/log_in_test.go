package v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func TestService_LogIn(t *testing.T) {
	t.Parallel()

	type args struct {
		session  model.Session
		name     string
		password string
		hash     string
	}

	type result struct {
		authTokens model.AuthTokens
		err        error
	}

	type testCase struct {
		name   string
		args   args
		mocks  func(tc testCase, deps *serviceTestDeps)
		result result
	}

	testCases := []testCase{
		{
			name: "success",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					GetPasswordHash(gomock.Any(), tc.args.name).
					Return(testHash, nil)
				deps.passportRepository.EXPECT().
					GetSessionInfo(gomock.Any(), tc.args.name).
					Return(testSession, nil)

			},
			args: args{
				name:     testName,
				session:  testSession,
				password: testPassword,
				hash:     testHash,
			},
			result: result{
				authTokens: testAuthTokens,
				err:        nil,
			},
		},
		{
			name: "hash not found",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					GetPasswordHash(gomock.Any(), tc.args.name).
					Return("", model.ErrUserNotFound)

			},
			args: args{
				name:     testName,
				session:  model.Session{},
				password: testPassword,
				hash:     "",
			},
			result: result{
				authTokens: model.AuthTokens{},
				err:        model.ErrUserNotFound,
			},
		},
		{
			name: "session not found",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					GetPasswordHash(gomock.Any(), tc.args.name).
					Return(testHash, nil)
				deps.passportRepository.EXPECT().
					GetSessionInfo(gomock.Any(), tc.args.name).
					Return(model.Session{}, model.ErrUserNotFound)

			},
			args: args{
				name:     testName,
				session:  model.Session{},
				password: testPassword,
				hash:     testHash,
			},
			result: result{
				authTokens: model.AuthTokens{},
				err:        model.ErrUserNotFound,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := getTestDeps(t)

			tc.mocks(tc, deps)
			authTokens, err := deps.Service.LogIn(deps.ctx, tc.args.name, tc.args.password)

			require.ErrorIs(t, err, tc.result.err)

			if err == nil {
				session, _, _ := extractTokenClaims(authTokens.AccessToken, deps.Service.config.JWTSecret)
				resSession, _, _ := extractTokenClaims(tc.result.authTokens.RefreshToken, deps.Service.config.JWTSecret)
				require.Equal(t, resSession.UserID, session.UserID)
				require.Equal(t, resSession.DeviceID, session.DeviceID)

			}

		})
	}

}
