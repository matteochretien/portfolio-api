package token

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const JwtAccessTokenLifetime = 15 * time.Minute
const JwtRefreshTokenLifetime = 24 * time.Hour

var UnableToSignAccessTokenErr = fmt.Errorf("unable to sign access token")
var UnableToSignRefreshTokenErr = fmt.Errorf("unable to sign refresh token")
var UnableToValidateAccessTokenErr = fmt.Errorf("unable to validate access token")
var AccessTokenExpiredErr = fmt.Errorf("access token expired")
var UnableToValidateRefreshTokenErr = fmt.Errorf("unable to validate refresh token")
var RefreshTokenExpiredErr = fmt.Errorf("refresh token expired")
var InvalidTokenTypeErr = fmt.Errorf("invalid token type")

type JwtTokenStrategy struct {
	secret string
}

func NewJwtTokenStrategy(secret string) func() TokenStrategy {
	return func() TokenStrategy {
		return &JwtTokenStrategy{secret}
	}
}

func (j *JwtTokenStrategy) GenerateAccessToken(claims map[string]any) (string, error) {
	if claims == nil {
		claims = make(map[string]any)
	}
	claims["exp"] = time.Now().Add(JwtAccessTokenLifetime).Unix()
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	signedAccessToken, err := accessToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", UnableToSignAccessTokenErr, err)
	}

	return signedAccessToken, nil
}

func (j *JwtTokenStrategy) GenerateRefreshToken(claims map[string]any) (string, error) {
	if claims == nil {
		claims = make(map[string]any)
	}
	claims["exp"] = time.Now().Add(JwtRefreshTokenLifetime).Unix()
	claims["type"] = "refresh"
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))

	signedRefreshToken, err := refreshToken.SignedString([]byte(j.secret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", UnableToSignRefreshTokenErr, err)
	}

	return signedRefreshToken, nil
}

func (j *JwtTokenStrategy) ValidateAccessToken(token string) (map[string]any, error) {
	accessToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", UnableToValidateAccessTokenErr, err)
	}

	claims, ok := accessToken.Claims.(jwt.MapClaims)
	if !ok || !accessToken.Valid {
		return nil, fmt.Errorf("%w: %w", AccessTokenExpiredErr, err)
	}

	// If the token is not an access token, return an error
	if claims["type"] == "refresh" {
		return nil, InvalidTokenTypeErr
	}

	return claims, nil
}

func (j *JwtTokenStrategy) ValidateRefreshToken(token string) (map[string]any, error) {
	refreshToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", UnableToValidateRefreshTokenErr, err)
	}

	claims, ok := refreshToken.Claims.(jwt.MapClaims)
	if !ok || !refreshToken.Valid {
		return nil, fmt.Errorf("%w: %w", RefreshTokenExpiredErr, err)
	}

	if claims["type"] != "refresh" {
		return nil, InvalidTokenTypeErr
	}

	return claims, nil
}
