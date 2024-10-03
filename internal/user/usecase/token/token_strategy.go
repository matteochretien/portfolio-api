package token

type TokenStrategy interface {
	GenerateAccessToken(claims map[string]any) (string, error)
	GenerateRefreshToken(claims map[string]any) (string, error)
	ValidateAccessToken(token string) (map[string]any, error)
	ValidateRefreshToken(token string) (map[string]any, error)
}
