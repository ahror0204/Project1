package v1

import (
	"context"

	"net/http"
	"strconv"
	"time"

	pb "github.com/project1/apigate/genproto"
	l "github.com/project1/apigate/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"google.golang.org/protobuf/encoding/protojson"
)

type Media struct {
	Id   string `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Type string `protobuf:"bytes,2,opt,name=type,proto3" json:"type"`
	Link string `protobuf:"bytes,3,opt,name=link,proto3" json:"link"`
}

type Post struct {
	Id          string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	Name        string  `protobuf:"bytes,2,opt,name=name,proto3" json:"name"`
	Description string  `protobuf:"bytes,3,opt,name=description,proto3" json:"description"`
	UserId      string  `protobuf:"bytes,4,opt,name=user_id,json=userId,proto3" json:"user_id"`
	Medias      []Media `protobuf:"bytes,5,rep,name=medias,proto3" json:"medias"`
}

type Address struct {
	City       string `protobuf:"bytes,1,opt,name=city,proto3" json:"city"`
	Country    string `protobuf:"bytes,2,opt,name=country,proto3" json:"country"`
	District   string `protobuf:"bytes,3,opt,name=district,proto3" json:"district"`
	PostalCode int64  `protobuf:"varint,4,opt,name=postal_code,json=postalCode,proto3" json:"postal_code"`
}

type CreateUserReqBody struct {
	Id           string    `protobuf:"bytes,1,opt,name=id,proto3" json:"id"`
	FirstName    string    `protobuf:"bytes,2,opt,name=first_name,json=firstName,proto3" json:"first_name"`
	LastName     string    `protobuf:"bytes,3,opt,name=last_name,json=lastName,proto3" json:"last_name"`
	Email        string    `protobuf:"bytes,4,opt,name=email,proto3" json:"email"`
	Bio          string    `protobuf:"bytes,5,opt,name=bio,proto3" json:"bio"`
	PhoneNumbers []string  `protobuf:"bytes,6,rep,name=phone_numbers,json=phoneNumbers,proto3" json:"phone_numbers"`
	Address      []Address `protobuf:"bytes,7,rep,name=address,proto3" json:"address"`
	Status       string    `protobuf:"bytes,8,opt,name=status,proto3" json:"status"`
	CreatedAt    string    `protobuf:"bytes,9,opt,name=created_at,json=createdAt,proto3" json:"created_at"`
	UpdatedAt    string    `protobuf:"bytes,10,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at"`
	DeletedAt    string    `protobuf:"bytes,11,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at"`
	Posts        []Post    `protobuf:"bytes,12,rep,name=posts,proto3" json:"posts"`
}

// @Summary Create user
// @Description This api uses for creating new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserReqBody true "user body"
// @Success 200 {string} Success
// @Router /v1/users [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        pb.User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().CreateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	bodyByte, err := json.Marshal(response)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed set to redis", l.Error(err))
		return
	}

	err = h.redisStorage.Set(body.FirstName, string(bodyByte))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed set to redis", l.Error(err))
		return
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Get user by id
// @Description This api uses for getting user by id
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {string} CreateUserReqBody
// @Router /v1/users/{id} [get]
func (h *handlerV1) GetUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	guid := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetUserById(
		ctx, &pb.GetUserByIdRequest{
			Id: guid,
		})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	// redisValue, err := h.redisStorage.Get(response.FirstName)

	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	h.log.Error("failed set to redis", l.Error(err))
	// 	return
	// }

	// fmt.Printf(string(redisValue))

	c.JSON(http.StatusOK, response)
}

// @Summary Get all users
// @Description This api uses for getting all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {string} CreateUserReqBody
// @Router /v1/users/all [get]
func (h *handlerV1) GetAllUser(c *gin.Context) {

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().GetAllUser(
		ctx, &pb.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed get all users", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// @ListUsers returns list of users
// @Summary Update user by id
// @Description This api uses for updating user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body CreateUserReqBody true "user body"
// @Success 200 {string} CreateUserReqBody
// @Router /v1/usersupdate/{id} [put]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	var (
		body        pb.User
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}
	body.Id = c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UpdateUser(ctx, &body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to update user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// UserList get limited users

// @ListUsers returns list of users
// @Summary Get users list
// @Description This api uses for getting users list
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int true "limit"
// @Param page query int true "page"
// @Success 200 {string} CreateUserReqBody
// @Router /v1/users/list [get]
func (h *handlerV1) UserList(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	limitValue, _ := strconv.ParseInt(limit, 10, 64)
	pageValue, _ := strconv.ParseInt(page, 10, 64)

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UserList(
		ctx,
		&pb.UserListRequest{
			Limit: limitValue,
			Page:  pageValue,
		},
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("ailed to list users", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// // DeleteUser deletes user by id
// // route /v1/users/{id} [delete]
// func (h *handlerV1) DeleteUser(c *gin.Context) {
// 	var jspbMarshal protojson.MarshalOptions
// 	jspbMarshal.UseProtoNames = true

// 	guid := c.Param("id")
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
// 	defer cancel()

// 	response, err := h.serviceManager.UserService().Delete(
// 		ctx, &pb.ByIdReq{
// 			Id: guid,
// 		})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"error": err.Error(),
// 		})
// 		h.log.Error("failed to delete user", l.Error(err))
// 		return
// 	}

// 	c.JSON(http.StatusOK, response)
// }
