package v1

import (
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/project1/apigate/api/auth"
	repo "github.com/project1/apigate/storage/repo"

	"github.com/project1/apigate/config"
	"github.com/project1/apigate/pkg/logger"
	"github.com/project1/apigate/services"
)

type handlerV1 struct {
	log            logger.Logger
	serviceManager services.IServiceManager
	cfg            config.Config
	redisStorage   repo.RedisRepositoryStorage
	jwtHandler     auth.JWtHandler
}

// HandlerV1Config ...
type HandlerV1Config struct {
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	Cfg            config.Config
	Redis          repo.RedisRepositoryStorage
	jwtHandler     auth.JWtHandler
}

// New ...
func New(c *HandlerV1Config) *handlerV1 {
	return &handlerV1{
		log:            c.Logger,
		serviceManager: c.ServiceManager,
		cfg:            c.Cfg,
		redisStorage:   c.Redis,
		jwtHandler:     c.jwtHandler,
	}
}

// func CheckClaims(h *handlerV1, c *gin.Context) jwt.MapClaims {
// 	var (
// 		ErrUnauthorized = errors.New("unauthorized")
// 		authorization   JwtRequestModel
// 		claims          jwt.MapClaims
// 		err             error
// 	)

// 	authorization.Token = c.GetHeader("Authorization")
// 	if c.Request.Header.Get("Authorization") == "" {
// 		c.JSON(http.StatusUnauthorized, ErrUnauthorized)
// 		h.log.Error("Unauthorized request: ", logger.Error(ErrUnauthorized))
// 		return nil
// 	}

// 	h.jwtHandler.Token = authorization.Token
// 	claims, err = h.jwtHandler.ExtractClaims()
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, ErrUnauthorized)
// 		h.log.Error("Invalid token", logger.Error(ErrUnauthorized))
// 		return nil
// 	} 	
// 	return claims
// }

func CheckClaims(h *handlerV1, c *gin.Context) jwt.MapClaims {
	var (
	  ErrUnauthorized = errors.New("unauthorized")
	  authorization   JwtRequestModel
	  claims          jwt.MapClaims
	  err             error
	)
  
	authorization.Token = c.GetHeader("Authorization")
	if c.Request.Header.Get("Authorization") == "" {
	  c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	  h.log.Error("Unauthorized request:", logger.Error(err))
  
	}
	h.jwtHandler.Token = authorization.Token
	claims, err = h.jwtHandler.ExtractClaims()
	if err != nil {
	  c.JSON(http.StatusUnauthorized, ErrUnauthorized)
	  h.log.Error("token is invalid:", logger.Error(err))
	  return nil
	}
	return claims
  }