package v1

import (
	"context"
	"math/rand"
	"strconv"
	"strings"
	"unicode"

	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/project1/apigate/api/auth"
	pb "github.com/project1/apigate/genproto"
	l "github.com/project1/apigate/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"google.golang.org/protobuf/encoding/protojson"

	"crypto/tls"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	// "github.com/google/uuid"
	gomail "gopkg.in/mail.v2"
)

// @Summary Register User
// @Description This api uses for registration new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserReqBody true "user body"
// @Success 200 {string} Success
// @Router /v1/users/register [post]
func (h handlerV1) RegisterUser(c *gin.Context) {
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

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	//----------------------------------------------------------------
	status, err := h.serviceManager.UserService().CheckField(ctx, &pb.UserCheckRequest{
		Field: "username",
		Value: body.UserName,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed while calling CheckField function whith USERNAME", l.Error(err))
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed set to redis", l.Error(err))
		return
	}

	if !status.Response {
		status2, err := h.serviceManager.UserService().CheckField(ctx, &pb.UserCheckRequest{
			Field: "email",
			Value: body.Email,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed while calling CheckField function with EMAIL", l.Error(err))
			return
		}

		if status2.Response {
			c.JSON(http.StatusConflict, gin.H{
				"error": "user_name already in use",
			})
			h.log.Error("User already exists", l.Error(err))
			return
		}
	} else {
		c.JSON(http.StatusConflict, gin.H{
			"error": "email already in use",
		})
		h.log.Error("User already exists", l.Error(err))
		return
	}

	// verifyPassword...
	err = verifyPassword(body.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed while calling CheckField function whith USERNAME", l.Error(err))
		return
	}

	//Hashing the user Password

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("error while hashing password", l.Error(err))
		return
	}
	body.Password = string(hashedPassword)

	min := 99999
	max := 1000000
	rand.Seed(time.Now().UnixNano())
	Code := rand.Intn(max-min) + min

	verCode := strconv.Itoa(Code)

	body.EmailCode = verCode

	SendEmail(body.Email, verCode)

	setBodyRedis, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed set to redis IN REGISTER FUNC  1", l.Error(err))
		return
	}

	err = h.redisStorage.Set(body.Email, string(setBodyRedis))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed set to redis IN REGISTER FUNC  2", l.Error(err))
		return
	}
}

// @Summary Send Email Code
// @Description This api uses for sendin email code to user
// @Tags users
// @Accept json
// @Produce json
// @Param user body EmailVer true "user body"
// @Success 200 {string} Success
// @Router /v1/users/verification [post]
func (h handlerV1) VerifyUser(c *gin.Context) {

	var mailData EmailVer

	err := c.ShouldBindJSON(&mailData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json in VerifyUser func", l.Error(err))
		return
	}
	mailData.Email = strings.TrimSpace(mailData.Email)
	mailData.Email = strings.ToLower(mailData.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	bod, err := redis.String(h.redisStorage.Get(mailData.Email))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		h.log.Error("failed get from redis", l.Error(err))
		return
	}

	var redisBody *pb.User

	err = json.Unmarshal([]byte(bod), &redisBody)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed while using CreateUser func in verify user func!!!", l.Error(err))
		return
	}
	//------------------------------------------------------------------------------------------

	h.jwtHandler = auth.JWtHandler{
		Sub:  redisBody.Id,
		Iss:  "client",
		Role: "authorized",
		Log:  h.log,
		SigningKey: h.cfg.SigningKey,
	}

	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error while generating jwt",
		})
		h.log.Error("error while generating jwt tokens", l.Error(err))
		return
	}

	
//--------------------------------------------------------------------------------------
	if mailData.EmailCode == redisBody.EmailCode {

		createVal, err := h.serviceManager.UserService().CreateUser(ctx, redisBody)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			h.log.Error("failed while using CreateUser func in verify user func!!!", l.Error(err))
			return
		}

		var frontResp = CreateUserReqBody {
			Id: createVal.Id,
			FirstName: createVal.FirstName,
			LastName: createVal.LastName,
			Email: createVal.Email,
			Bio: createVal.Bio,
			PhoneNumbers: createVal.PhoneNumbers,
			Status: createVal.Status,
			UserName: createVal.UserName,
			Password: createVal.Password,
			RefreshToken: refresh,
			AccessToken: access,
		}

		c.JSON(http.StatusCreated, frontResp)
	}
}

func SendEmail(email, code string) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "ahrorahrorovnt@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", email)
	// id,err := uuid.NewUUID()
	// if err != nil {
	//   fmt.Println(err)
	// }
	// Set E-Mail subject
	m.SetHeader("code:", "Verification code")

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", code)

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "ahrorahrorovnt@gmail.com", "qmxlgijkvuuoacrh")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

}

func verifyPassword(password string) error {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	const minPassLength = 8
	const maxPassLength = 32
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if !numberPresent {
		appendError("atleast one numeric character required")
	}
	if !specialCharPresent {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength))
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}
	return nil
}
