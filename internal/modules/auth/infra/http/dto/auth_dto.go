package dto

import (
	"github.com/boilerplate/internal/domain"
	"github.com/boilerplate/internal/modules/auth/application"
	userDto "github.com/boilerplate/internal/modules/user/infra/http/dto"
)

type LoginReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// TokenRes is the token data in the response
type TokenRes struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// LoginRes is the complete login response with both token and user
type LoginRes struct {
	Token *TokenRes       `json:"token"`
	User  userDto.UserRes `json:"user"`
}

// LoginResFromPairAndUser converts both token pair and user to LoginRes
func AuthResFromPairAndUser(pair *application.TokenPair, user *domain.User) LoginRes {
	return LoginRes{
		Token: &TokenRes{
			AccessToken:  pair.AccessToken,
			RefreshToken: pair.RefreshToken,
			TokenType:    pair.TokenType,
			ExpiresIn:    pair.ExpiresIn,
		},
		User: userDto.UserResFromDomain(user),
	}
}
