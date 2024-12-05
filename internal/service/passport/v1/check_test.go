package v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func TestService_Check(t *testing.T) {
	t.Parallel()

	type args struct {
		token    string
		version  int
		deviceID model.DeviceID
		userID   model.UserID
	}

	type result struct {
		userID    model.UserID
		deviceID  model.DeviceID
		tokenType string
		err       error
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
					GetTokenVersion(gomock.Any(), tc.args.userID, tc.args.deviceID).
					Return(tc.args.version, nil)

			},
			args: args{
				token:    testRefreshToken,
				version:  testVersion,
				userID:   testUserID,
				deviceID: testDeviceID,
			},

			result: result{
				userID:    testUserID,
				deviceID:  testDeviceID,
				tokenType: "refresh",
				err:       nil,
			},
		},
		{
			name:  "expired token",
			mocks: func(tc testCase, deps *serviceTestDeps) {},
			args: args{
				token: testExpiredToken,
			},

			result: result{
				err: model.ErrExpiredToken,
			},
		},
		{
			name: "expired session",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					GetTokenVersion(gomock.Any(), tc.args.userID, tc.args.deviceID).
					Return(tc.args.version-1, nil)

			},
			args: args{
				token:    testRefreshToken,
				version:  testVersion,
				userID:   testUserID,
				deviceID: testDeviceID,
			},

			result: result{
				userID:    testUserID,
				deviceID:  testDeviceID,
				tokenType: "refresh",
				err:       model.ErrExpiredSession,
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
				version:  testVersion,
				userID:   testUserID,
				deviceID: testDeviceID,
			},

			result: result{
				err: model.ErrUserNotFound,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := getTestDeps(t)

			tc.mocks(tc, deps)
			userID, deviceID, tokenType, err := deps.Service.CheckToken(deps.ctx, tc.args.token)

			require.ErrorIs(t, err, tc.result.err)

			if err == nil {
				require.Equal(t, tc.result.userID, userID)
				require.Equal(t, tc.result.deviceID, deviceID)
				require.Equal(t, tc.result.tokenType, tokenType)
			}

		})
	}
}
