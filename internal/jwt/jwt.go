package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secretKey       string
	accessTokenExp  int
	refreshTokenExp int
}

func NewJWTService(secretKey string, accessTokenExpiry, refreshTokenExpiry int) *JWTService {
	return &JWTService{
		secretKey:       secretKey,
		accessTokenExp:  accessTokenExpiry,
		refreshTokenExp: refreshTokenExpiry,
	}
}

type Claims struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

const (
	AccessTokenType  = "access"
	RefreshTokenType = "refresh"
)

func (m *JWTService) GenerateAccessRefreshTokens(userID, username string) (string, string, error) {

	// Generate access token
	accessToken, err := m.GenerateToken(userID, username, AccessTokenType)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token
	refreshToken, err := m.GenerateToken(userID, username, RefreshTokenType)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (m *JWTService) GenerateToken(userID, username string, tokenType string) (string, error) {
	var expirationSec int
	if tokenType == AccessTokenType {
		expirationSec = m.accessTokenExp
	} else if tokenType == RefreshTokenType {
		expirationSec = m.refreshTokenExp
	} else {
		return "", errors.New("wrong tokenType")
	}

	issuedTime := time.Now()
	expirationTime := issuedTime.Add(time.Duration(expirationSec) * time.Second)

	claims := &Claims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "go-grpc-sso",
			IssuedAt:  jwt.NewNumericDate(issuedTime),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (m *JWTService) ValidateTokenWithType(tokenString, tokenType string) (*Claims, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.TokenType != tokenType {
		return nil, errors.New("invalid token type")
	}
	return claims, nil
}

func (m *JWTService) ValidateToken(tokenString string) (*Claims, error) {

	// Validate signed by us
	token, err := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(m.secretKey), nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func (m *JWTService) RefreshTokens(refreshTokenString string) (string, string, error) {
	claims, err := m.ValidateTokenWithType(refreshTokenString, RefreshTokenType)
	if err != nil {
		return "", "", err
	}

	return m.GenerateAccessRefreshTokens(claims.UserID, claims.Username)
}
