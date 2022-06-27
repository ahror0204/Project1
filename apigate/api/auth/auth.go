package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/project1/apigate/pkg/logger"
)

type JWtHandler struct {
	Sub       string
	Iss       string
	Exp       string
	Iat       string
	Aud       []string
	Role      string
	SigningKey string
	Log       logger.Logger
	Token     string
}

// Genereting Access and Refresh Tokens
func (jwtHandler *JWtHandler) GenerateAuthJWT() (access, refresh string, err error) {
	var (
		accessToken  *jwt.Token
		refreshToken *jwt.Token
		claims       jwt.MapClaims
	)
	accessToken = jwt.New(jwt.SigningMethodHS256)
	refreshToken = jwt.New(jwt.SigningMethodHS256)

	claims = accessToken.Claims.(jwt.MapClaims)
	claims["iss"] = jwtHandler.Iss
	claims["sub"] = jwtHandler.Sub
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["iat"] = time.Now().Unix() 
	claims["aud"] = jwtHandler.Aud
	claims["role"] = jwtHandler.Role
	claims["signingkey"] = jwtHandler.SigningKey
	claims["token"] = jwtHandler.Token
	  
	access, err = accessToken.SignedString([]byte(jwtHandler.SigningKey))
	if err != nil {
		jwtHandler.Log.Error("error genereting accessToken", logger.Error(err))
		return
	}

	refresh, err = refreshToken.SignedString([]byte(jwtHandler.SigningKey))

	if err != nil {
		jwtHandler.Log.Error("error genereting accessToken", logger.Error(err))
		return
	}

	return access, refresh, nil
}

// ExtractClaims ...
func (jwtHandler *JWtHandler) ExtractClaims() (jwt.MapClaims, error) {
	var (
		token *jwt.Token
		err error
	)

	token, err = jwt.Parse(jwtHandler.Token, func(t *jwt.Token) (interface{}, error) {
		return []byte(jwtHandler.SigningKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		jwtHandler.Log.Error("invalid jwt token")
		return nil, err
	}
	return claims, nil
}