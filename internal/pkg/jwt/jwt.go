package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/oxiginedev/sabipass/config"
	"github.com/oxiginedev/sabipass/internal/models"
)

const defaultTokenExpiry = 1 * time.Hour

var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type ValidatedToken struct {
	UserID    string
	ExpiresIn int64
}

type TokenManager interface {
	GenerateToken(user *models.User) (Token, error)
	ValidateToken(tokenString string) (*ValidatedToken, error)
}

type jwtTokenManager struct {
	secretKey string
	expiry    time.Duration
}

func NewJwtTokenManager(cfg *config.Config) TokenManager {
	if cfg.Auth.JWT.Expiry == 0 {
		cfg.Auth.JWT.Expiry = defaultTokenExpiry
	}

	return &jwtTokenManager{
		secretKey: cfg.Auth.JWT.SecretKey,
		expiry:    cfg.Auth.JWT.Expiry,
	}
}

func (j *jwtTokenManager) GenerateToken(user *models.User) (Token, error) {
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(j.expiry).Unix(),
		"iat": time.Now().Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := jwtToken.SignedString([]byte(j.secretKey))
	if err != nil {
		return Token{}, err
	}

	return Token{
		AccessToken: accessToken,
		ExpiresIn:   claims["exp"].(int64),
	}, nil
}

func (j *jwtTokenManager) ValidateToken(tokenString string) (*ValidatedToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("[jwt]: unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && verr.Errors == jwt.ValidationErrorExpired {
			return nil, ErrTokenExpired
		}
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		userId := claims["sub"].(string)
		expiresIn := claims["exp"].(float64)

		return &ValidatedToken{
			UserID:    userId,
			ExpiresIn: int64(expiresIn),
		}, nil
	}

	return nil, ErrInvalidToken
}
