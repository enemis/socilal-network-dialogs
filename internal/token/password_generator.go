package token

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"social-network-dialogs/internal/config"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type PasswordGenerator interface {
	CompareHashAndPassword(userId uuid.UUID, hashedPassword, password string) bool
	PasswordHash(userId uuid.UUID, password string) (string, error)
	ParseToken(accessToken string) (uuid.UUID, error)
}

type Claims struct {
	jwt.StandardClaims
}

type PasswordGeneratorInstance struct {
	config *config.Config
}

func NewPasswordGenerator(config *config.Config) *PasswordGeneratorInstance {
	return &PasswordGeneratorInstance{config: config}
}

func (s *PasswordGeneratorInstance) CompareHashAndPassword(userId uuid.UUID, hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(s.buildPasswordString(userId, password))) == nil
}

func (s *PasswordGeneratorInstance) PasswordHash(userId uuid.UUID, password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(s.buildPasswordString(userId, password)), 14)
	return string(bytes), err
}

func (s *PasswordGeneratorInstance) buildPasswordString(userId uuid.UUID, password string) string {
	return strings.Join([]string{password, s.config.Salt, userId.String()}, "")
}

func (s *PasswordGeneratorInstance) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.config.SigningKey), nil
	})

	if err = token.Claims.Valid(); err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return uuid.Nil, errors.New("token claims are not of type Claims")
	}

	return uuid.Parse(claims.Issuer)
}
