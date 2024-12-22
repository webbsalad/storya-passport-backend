package v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func TestService_Delete(t *testing.T) {
	t.Parallel()

	type args struct {
		userID  model.UserID
		emailID model.EmailID
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
					Delete(gomock.Any(), tc.args.userID, tc.args.emailID).
					Return(nil)
			},
			args: args{
				userID:  testUserID,
				emailID: testEmailID,
			},
			result: result{
				err: nil,
			},
		},

		{
			name: "user not found",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					Delete(gomock.Any(), tc.args.userID, tc.args.emailID).
					Return(model.ErrUserNotFound)
			},
			args: args{
				userID:  testUserID,
				emailID: testEmailID,
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
			err := deps.Service.Delete(deps.ctx, tc.args.userID, tc.args.emailID)

			require.ErrorIs(t, err, tc.result.err)

		})
	}
}
