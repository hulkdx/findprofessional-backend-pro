package user

import (
	"strings"
)

type Service interface {
	IsAuthenticated(auth string) bool
}

type serviceImpl struct {
}

func NewService() Service {
	return &serviceImpl{}
}

func (s *serviceImpl) IsAuthenticated(auth string) bool {
	authSplit := strings.Split(auth, " ")
	if len(authSplit) != 2 {
		return false
	}
	authType := authSplit[0]
	accessToken := authSplit[1]
	if authType != "Bearer" {
		return false
	}
	return isValidAccessToken(accessToken)
}

func isValidAccessToken(accessToken string) bool {
	return true
}
