package model

type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}

func NewAuthTokensFromStrings(accessToken, refreshToken string) AuthTokens {
	return AuthTokens{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
