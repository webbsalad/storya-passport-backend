package v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func TestService_Get(t *testing.T) {
	t.Parallel()

	type args struct {
		userID model.UserID
		user   model.User
	}

	type result struct {
		user model.User
		err  error
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
					GetUser(gomock.Any(), tc.args.userID).
					Return(tc.args.user, nil)

			},
			args: args{
				userID: testUserID,
				user:   testUser,
			},

			result: result{
				user: testUser,

				err: nil,
			},
		},
		{
			name: "user not found",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					GetUser(gomock.Any(), tc.args.userID).
					Return(model.User{}, model.ErrUserNotFound)

			},
			args: args{
				userID: testUserID,
				user:   testUser,
			},

			result: result{
				user: model.User{},
				err:  model.ErrUserNotFound,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := getTestDeps(t)

			tc.mocks(tc, deps)
			user, err := deps.Service.GetUser(deps.ctx, tc.args.userID)

			require.ErrorIs(t, err, tc.result.err)

			if err == nil {
				require.Equal(t, tc.result.user, user)

			}

		})
	}

}
