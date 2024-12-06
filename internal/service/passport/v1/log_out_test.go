package v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func TestService_LogOut(t *testing.T) {
	t.Parallel()

	type args struct {
		deviceID model.DeviceID
		userID   model.UserID
	}

	type result struct {
		err error
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
					LogOut(gomock.Any(), tc.args.userID, tc.args.deviceID).
					Return(nil)
			},
			args: args{
				userID:   testUserID,
				deviceID: testDeviceID,
			},
			result: result{
				err: nil,
			},
		},
		{
			name: "device not found",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					LogOut(gomock.Any(), tc.args.userID, tc.args.deviceID).
					Return(model.ErrDeviceNotFound)
			},
			args: args{
				userID:   testUserID,
				deviceID: testDeviceID,
			},
			result: result{
				err: model.ErrDeviceNotFound,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := getTestDeps(t)

			tc.mocks(tc, deps)
			err := deps.Service.LogOut(deps.ctx, tc.args.userID, tc.args.deviceID)

			require.ErrorIs(t, err, tc.result.err)

		})
	}
}
