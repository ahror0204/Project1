package v1

import (
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
