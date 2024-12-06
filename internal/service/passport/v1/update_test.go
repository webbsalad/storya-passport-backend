package v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

func TestService_Update(t *testing.T) {
	t.Parallel()

	type args struct {
		userID   model.UserID
		name     string
		password string
		user     model.User
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
					UpdateUser(gomock.Any(), tc.args.userID, tc.args.name, hashMatcher{password: tc.args.password}).
					Return(tc.args.user, nil)
			},
			args: args{
				name:     testName,
				password: testPassword,
				userID:   testUserID,
				user:     testUser,
			},
			result: result{
				user: testUser,
				err:  nil,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := getTestDeps(t)

			tc.mocks(tc, deps)
			user, err := deps.Service.UpdateUser(deps.ctx, tc.args.userID, tc.args.name, tc.args.password)

			require.ErrorIs(t, err, tc.result.err)

			if err == nil {
				require.Equal(t, tc.result.user, user)
			}
		})
	}
}
