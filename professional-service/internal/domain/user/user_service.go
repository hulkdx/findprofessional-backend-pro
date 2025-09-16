package user

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hulkdx/findprofessional-backend-pro/professional-service/internal/utils/logger"
)

const publicKeyPath = "/config/rsa.public.key"
const baseUrl = "http://user-service:8080"
const loginUrl = baseUrl + "/auth/login"

type Service interface {
	IsAuthenticated(ctx context.Context, auth string) bool
	Login(ctx context.Context, email string, password string) (string, error)
	GetAuthenticatedUserId(ctx context.Context, auth string) (int64, error)
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

func getAccessToken(accessToken string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	return jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
}

func isValidAccessToken(accessToken string, publicKey *rsa.PublicKey) bool {
	token, err := getAccessToken(accessToken, publicKey)
	if err != nil {
		return false
	}
	return token.Valid
}

func (s *serviceImpl) Login(ctx context.Context, email, password string) (string, error) {
	loginReq := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, email, password)

	byt, err := httpRequest(ctx, http.MethodPost, loginUrl, strings.NewReader(loginReq))
	return string(byt), err
}

func httpRequest(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return io.ReadAll(res.Body)
}

func (s *serviceImpl) GetAuthenticatedUserId(ctx context.Context, auth string) (int64, error) {
	accessTokenString := getAccessTokenFromAuthHeader(auth)
	accessToken, err := getAccessToken(accessTokenString, s.publicKey)
	if err != nil {
		return -2, err
	}
	subject, err := accessToken.Claims.GetSubject()
	if err != nil {
		return -2, err
	}
	idInt, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return -2, err
	}
	return idInt, nil
}
