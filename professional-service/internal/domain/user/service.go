package user

import (
	"context"
	"crypto/rsa"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
	"os"
	"strings"
)

const publicKeyPath = "/config/rsa.public.key"

type Service interface {
	IsAuthenticated(ctx context.Context, auth string) bool
}

type serviceImpl struct {
	publicKey *rsa.PublicKey
}

func NewService() Service {
	publicKeyFile, err := os.ReadFile(publicKeyPath)
	if err != nil {
		logger.Error("Failed to open public key file: ", err)
	}
	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyFile)
	if err != nil {
		logger.Error("Failed to parse public key file: ", err)
	}
	return &serviceImpl{
		publicKey: publicKey,
	}
}

func (s *serviceImpl) IsAuthenticated(ctx context.Context, auth string) bool {
	accessToken := getAccessTokenFromAuthHeader(auth)
	if accessToken == "" {
		return false
	}
	return isValidAccessToken(accessToken, s.publicKey)
}

func getAccessTokenFromAuthHeader(auth string) string {
	authSplit := strings.Split(auth, " ")
	if len(authSplit) != 2 {
		return ""
	}
	authType := authSplit[0]
	accessToken := authSplit[1]
	if authType != "Bearer" {
		return ""
	}
	return accessToken
}

func isValidAccessToken(accessToken string, publicKey *rsa.PublicKey) bool {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return false
	}
	return token.Valid
}
