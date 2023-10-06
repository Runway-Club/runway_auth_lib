package jwt

import (
	"github.com/Runway-Club/auth_lib/domain"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
	"time"
)

type DummyJwtGenerator struct {
	issuer string
	exp    int64
	secret []byte
}

func (d DummyJwtGenerator) GenerateToken(auth *domain.Auth, payload map[string]interface{}) (string, error) {
	payload["id"] = auth.Id
	payload["username"] = auth.Username
	payload["role_id"] = auth.RoleId
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.MapClaims{
		"payload": payload,
		"exp":     time.Now().UnixMilli() + d.exp,
		"iat":     time.Now().UnixMilli(),
		"iss":     d.issuer,
	})
	tokenString, err := token.SignedString(d.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (d DummyJwtGenerator) VerifyToken(token string) (*domain.Auth, map[string]interface{}, error) {
	parsedToken, err := jwtlib.Parse(token, func(token *jwtlib.Token) (interface{}, error) {
		return d.secret, nil
	})
	if parsedToken == nil {
		return nil, nil, domain.ErrInvalidToken
	}
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}
	claims := parsedToken.Claims.(jwtlib.MapClaims)
	var payload map[string]interface{}
	err = mapstructure.Decode(claims["payload"], &payload)
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}
	auth := &domain.Auth{
		Id:       payload["id"].(string),
		Username: payload["username"].(string),
		RoleId:   payload["role_id"].(string),
	}
	return auth, payload, nil
}

func NewDummyJwtGenerator(issuer string, exp int64, secret string) *DummyJwtGenerator {
	return &DummyJwtGenerator{
		issuer: issuer,
		exp:    exp,
		secret: []byte(secret),
	}
}
