package sso

import (
	"context"
	"fmt"
	"grpc-sso/internal/db"
	"grpc-sso/internal/jwt"
	"grpc-sso/internal/messages"
	"grpc-sso/internal/validation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SSOService struct {
	UnimplementedSSOServiceServer
	jwts *jwt.JWTService
	us   *db.UserService
}

func NewSSOService(jwts *jwt.JWTService, us *db.UserService) *SSOService {
	return &SSOService{
		jwts: jwts,
		us:   us,
	}
}

func (s *SSOService) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Check password format
	if !validation.IsPasswordValid(req.GetPassword()) {
		return nil, status.Errorf(codes.InvalidArgument, messages.InvalidPassword)
	}

	// Check username format
	if !validation.IsUsernameValid(req.GetUsername()) {
		return nil, status.Errorf(codes.InvalidArgument, messages.InvalidUsername)
	}

	// Check email format
	if !validation.IsEmailValid(req.GetEmail()) {
		return nil, status.Errorf(codes.InvalidArgument, messages.InvalidEmail)
	}

	// Check unique email, username
	if us, _ := s.us.FindOneByEmail(req.GetEmail()); us != nil {
		return nil, status.Errorf(codes.AlreadyExists, messages.EmailAlreadyExists)
	}
	if us, _ := s.us.FindOneByUsername(req.GetUsername()); us != nil {
		return nil, status.Errorf(codes.AlreadyExists, messages.UsernameAlreadyExists)
	}

	user, err := s.us.CreateUser(req.GetEmail(), req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Unknown, messages.UnknownError)
	}

	return &RegisterResponse{
		UserId: fmt.Sprintf("%d", user.ID),
	}, nil
}

func (s *SSOService) Auth(ctx context.Context, req *AuthRequest) (*AuthResponse, error) {
	if req.GetUsername() == "" && req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, messages.UsernameOrEmailIsRequired)
	}

	// Impossible password
	if !validation.IsPasswordValid(req.GetPassword()) {
		return nil, status.Errorf(codes.NotFound, messages.AuthNotFound)
	}

	var user *db.User
	var err error
	if req.GetUsername() != "" {

		// Impossible username
		if !validation.IsUsernameValid(req.GetUsername()) {
			return nil, status.Errorf(codes.NotFound, messages.AuthNotFound)
		}

		user, err = s.us.AuthUserWithUsername(req.GetUsername(), req.GetPassword())
	} else if req.GetEmail() != "" {

		// Impossible email
		if !validation.IsEmailValid(req.GetEmail()) {
			return nil, status.Errorf(codes.NotFound, messages.AuthNotFound)
		}

		user, err = s.us.AuthUserWithEmail(req.GetEmail(), req.GetPassword())
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, messages.AuthNotFound)
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := s.jwts.GenerateAccessRefreshTokens(
		fmt.Sprintf("%d", user.ID),
		user.Username,
	)
	if err != nil {
		return nil, status.Errorf(codes.Unknown, messages.UnknownError)
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *SSOService) ValidateAccessToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	return s.validateTokenWithType(req.GetToken(), jwt.AccessTokenType)
}

func (s *SSOService) ValidateRefreshToken(ctx context.Context, req *ValidateTokenRequest) (*ValidateTokenResponse, error) {
	return s.validateTokenWithType(req.GetToken(), jwt.RefreshTokenType)
}

func (s *SSOService) validateTokenWithType(tokenString, tokenType string) (*ValidateTokenResponse, error) {
	claims, err := s.jwts.ValidateTokenWithType(tokenString, tokenType)
	if err != nil {
		return &ValidateTokenResponse{Valid: false}, nil
	}
	return &ValidateTokenResponse{
		Valid:     true,
		TokenType: claims.TokenType,
		ExpiresAt: claims.ExpiresAt.String(),
		UserId:    claims.UserID,
		Username:  claims.Username,
	}, nil

}

func (s *SSOService) RefreshTokens(ctx context.Context, req *RefreshTokensRequest) (*RefreshTokensResponse, error) {
	claims, err := s.jwts.ValidateTokenWithType(req.GetRefreshToken(), jwt.RefreshTokenType)
	if err != nil {
		return &RefreshTokensResponse{}, status.Errorf(codes.Unauthenticated, messages.Unauthenticated)
	}

	newAccessToken, newRefreshToken, err := s.jwts.GenerateAccessRefreshTokens(claims.UserID, claims.Username)
	if err != nil {
		return &RefreshTokensResponse{}, status.Errorf(codes.Unknown, messages.UnknownError)
	}
	return &RefreshTokensResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}, nil
}
