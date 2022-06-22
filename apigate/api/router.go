package api

import (
	v1 "github.com/project1/apigate/api/handlers/v1"
	"github.com/project1/apigate/config"
	"github.com/project1/apigate/pkg/logger"
	"github.com/project1/apigate/services"
	repo "github.com/project1/apigate/storage/repo"

	"github.com/gin-gonic/gin"

	_ "github.com/project1/apigate/api/docs"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

type Option struct {
	Conf           config.Config
	Logger         logger.Logger
	ServiceManager services.IServiceManager
	RedisRepo      repo.RedisRepositoryStorage
}

func New(option Option) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	handlerV1 := v1.New(&v1.HandlerV1Config{
		Logger:         option.Logger,
		ServiceManager: option.ServiceManager,
		Cfg:            option.Conf,
		Redis:          option.RedisRepo,
	})

	api := router.Group("v1")
	api.POST("/users/verification", handlerV1.VerifyUser)
	api.POST("/users/register", handlerV1.RegisterUser)
	api.POST("/users", handlerV1.CreateUser)
	api.GET("/users/:id", handlerV1.GetUser)
	api.GET("/users/all", handlerV1.GetAllUser)
	api.PUT("/usersupdate/:id", handlerV1.UpdateUser)
	api.GET("/users/list", handlerV1.UserList)
	// api.DELETE("/users/:id", handlerV1.DeleteUser)

	url := ginSwagger.URL("swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
