package v1

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/webbsalad/storya-passport-backend/internal/model"
)

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

func TestService_Register(t *testing.T) {
	t.Parallel()

	type args struct {
		session  model.Session
		emailID  model.EmailID
		name     string
		password string
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
					Register(gomock.Any(), tc.args.emailID, tc.args.name, hashMatcher{password: tc.args.password}).
					Return(tc.args.session, nil)
			},
			args: args{
				name:     testName,
				session:  testSession,
				password: testPassword,
				emailID:  testEmailID,
			},
			result: result{
				authTokens: testAuthTokens,
				err:        nil,
			},
		},
		{
			name: "email not confirmed",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					Register(gomock.Any(), tc.args.emailID, tc.args.name, hashMatcher{password: tc.args.password}).
					Return(tc.args.session, model.ErrEmailNotConfirmed)
			},
			args: args{
				name:     testName,
				session:  model.Session{},
				password: testPassword,
				emailID:  testEmailID,
			},
			result: result{
				authTokens: model.AuthTokens{},
				err:        model.ErrEmailNotConfirmed,
			},
		},
		{
			name: "user already exist",
			mocks: func(tc testCase, deps *serviceTestDeps) {
				deps.passportRepository.EXPECT().
					Register(gomock.Any(), tc.args.emailID, tc.args.name, hashMatcher{password: tc.args.password}).
					Return(tc.args.session, model.ErrUserAlreadyExist)
			},
			args: args{
				name:     testName,
				session:  model.Session{},
				password: testPassword,
				emailID:  testEmailID,
			},
			result: result{
				authTokens: model.AuthTokens{},
				err:        model.ErrUserAlreadyExist,
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			deps := getTestDeps(t)

			tc.mocks(tc, deps)
			authTokens, err := deps.Service.Register(deps.ctx, tc.args.emailID, tc.args.name, tc.args.password)

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
