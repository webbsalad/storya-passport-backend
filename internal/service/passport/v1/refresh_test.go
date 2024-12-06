package v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func TestService_Refresh(t *testing.T) {
	t.Parallel()

	type args struct {
		token    string
		userID   model.UserID
		deviceID model.DeviceID
		version  int
		session  model.Session
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
				updatedSession := tc.args.session
				updatedSession.Version += 1

				deps.passportRepository.EXPECT().
					GetTokenVersion(gomock.Any(), tc.args.userID, tc.args.deviceID).
					Return(tc.args.version, nil)
				deps.passportRepository.EXPECT().
					UpdateTokenVersion(gomock.Any(), updatedSession).
					Return(nil)

			},
			args: args{
				token:    testRefreshToken,
				userID:   testUserID,
				deviceID: testDeviceID,
				version:  testVersion,
				session:  testSession,
			},
			result: result{
				authTokens: testAuthTokens,
				err:        nil,
			},
		},
		{
			name:  "expired token",
			mocks: func(tc testCase, deps *serviceTestDeps) {},
			args: args{
				token:    testExpiredToken,
				userID:   model.UserID{},
				deviceID: model.DeviceID{},
				version:  0,
				session:  model.Session{},
			},
			result: result{
				authTokens: model.AuthTokens{},
				err:        model.ErrExpiredToken,
			},
		},
		{
			name:  "wrong token type",
			mocks: func(tc testCase, deps *serviceTestDeps) {},
			args: args{
				token:    testAccessToken,
				userID:   model.UserID{},
				deviceID: model.DeviceID{},
				version:  0,
				session:  model.Session{},
			},
			result: result{
				authTokens: model.AuthTokens{},
				err:        model.ErrWrongTokenType,
			},
		},
		{
			name: "expired session",
			mocks: func(tc testCase, deps *serviceTestDeps) {

				deps.passportRepository.EXPECT().
					GetTokenVersion(gomock.Any(), tc.args.userID, tc.args.deviceID).
					Return(tc.args.version+1, nil)

			},
			args: args{
				token:    testRefreshToken,
				userID:   testUserID,
				deviceID: testDeviceID,
				version:  testVersion,
				session:  model.Session{},
			},
			result: result{
				authTokens: model.AuthTokens{},
				err:        model.ErrExpiredSession,
			},
		},
		{
			name: "device not found",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				updatedSession := tc.args.session
				updatedSession.Version += 1

				deps.passportRepository.EXPECT().
					GetTokenVersion(gomock.Any(), tc.args.userID, tc.args.deviceID).
					Return(tc.args.version, nil)
				deps.passportRepository.EXPECT().
					UpdateTokenVersion(gomock.Any(), updatedSession).
					Return(model.ErrDeviceNotFound)

			},
			args: args{
				token:    testRefreshToken,
				userID:   testUserID,
				deviceID: testDeviceID,
				version:  testVersion,
				session:  testSession,
			},
			result: result{
				authTokens: model.AuthTokens{},
				err:        model.ErrDeviceNotFound,
			},
		},
		{
			name: "user not found",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					GetTokenVersion(gomock.Any(), tc.args.userID, tc.args.deviceID).
					Return(0, model.ErrUserNotFound)

			},
			args: args{
				token:    testRefreshToken,
				userID:   testUserID,
				deviceID: testDeviceID,
				version:  0,
				session:  testSession,
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
			authTokens, err := deps.Service.RefreshToken(deps.ctx, tc.args.token)

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
