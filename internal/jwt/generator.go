package jwt

import (
	"github.com/Runway-Club/auth_lib/domain"
	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/mitchellh/mapstructure"
	_ "github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"time"
)

type JwtGenerator struct {
	secret []byte
	exp    int64
	issuer string
}

func NewJwtGenerator() *JwtGenerator {
	secret := viper.GetString("runway_auth.jwt.secret")
	if secret == "" {
		panic("[required config] jwt.secret")
	}
	exp := viper.GetInt64("runway_auth.jwt.exp")
	if exp == 0 {
		panic("[required config] jwt.exp")
	}
	issuer := viper.GetString("runway_auth.jwt.issuer")
	return &JwtGenerator{
		secret: []byte(secret),
		exp:    exp,
		issuer: issuer,
	}
}

func (j JwtGenerator) GenerateToken(auth *domain.Auth, payload map[string]interface{}) (string, error) {
	payload["id"] = auth.Id
	payload["username"] = auth.Username
	payload["role_id"] = auth.RoleId
	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, jwtlib.MapClaims{
		"payload": payload,
		"exp":     j.exp*1000 + time.Now().UnixMilli(),
		"iat":     time.Now().UnixMilli(),
		"iss":     j.issuer,
	})
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j JwtGenerator) VerifyToken(token string) (*domain.Auth, map[string]interface{}, error) {
	// Bearer process
	if len(token) > 7 && token[0:7] == "Bearer " {
		token = token[7:]
	}
	parsedToken, err := jwtlib.Parse(token, func(token *jwtlib.Token) (interface{}, error) {
		return j.secret, nil
	})
	if parsedToken == nil {
		return nil, nil, domain.ErrInvalidToken
	}
	if err != nil {
		return nil, nil, domain.ErrInvalidToken
	}
	if !parsedToken.Valid {
		return nil, nil, domain.ErrInvalidToken
	}
	claims, ok := parsedToken.Claims.(jwtlib.MapClaims)
	if !ok {
		return nil, nil, domain.ErrInvalidToken
	}
	// check issuer
	issuer, ok := claims["iss"].(string)
	if !ok {
		return nil, nil, domain.ErrInvalidToken
	}
	if issuer != j.issuer {
		return nil, nil, domain.ErrInvalidIssuer
	}
	// check expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, nil, domain.ErrInvalidToken
	}
	if exp < float64(time.Now().UnixMilli()) {
		return nil, nil, domain.ErrExpiredToken
	}
	// get payload
	payload, ok := claims["payload"].(map[string]interface{})
	if !ok {
		return nil, nil, domain.ErrInvalidToken
	}

	parsedAuth := &domain.Auth{}
	// map[string]interface{} to domain.Auth parse
	mapstructure.Decode(claims["payload"], &parsedAuth)

	if parsedAuth.Id == "" {
		return nil, nil, domain.ErrInvalidToken
	}

	return parsedAuth, payload, nil
}
