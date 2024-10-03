package token

type RefreshAccessTokenUseCase struct {
	tokenStrategy TokenStrategy
}

func NewRefreshAccessTokenUseCase(tokenStrategy TokenStrategy) *RefreshAccessTokenUseCase {
	return &RefreshAccessTokenUseCase{tokenStrategy: tokenStrategy}
}

func (u *RefreshAccessTokenUseCase) Execute(refreshToken string) (string, error) {
	parsedRefreshToken, err := u.tokenStrategy.ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", err
	}

	accessToken, err := u.tokenStrategy.GenerateAccessToken(map[string]any{
		"email": parsedRefreshToken["email"],
	})
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
