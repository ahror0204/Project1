package casbin

import (
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin" 
	"github.com/project1/apigate/api/auth"
	"github.com/project1/apigate/config"
	"github.com/project1/apigate/api/model"
)



type JwtRoleStruct struct {
	enforce    *casbin.Enforcer
	conf       config.Config
	jwtHandler auth.JWtHandler
}

func NewJwtRoleStruct(e *casbin.Enforcer, c config.Config, jwtHandler auth.JWtHandler) gin.HandlerFunc {
	conf := &JwtRoleStruct{
		enforce:    e,
		conf:       c,
		jwtHandler: jwtHandler,
	}

	return func(c *gin.Context) {
		allow, err := conf.CheckPermission(c.Request)
		if err != nil {
			v, _ := err.(jwt.ValidationError)
			if v.Errors == jwt.ValidationErrorExpired {
				conf.RequireRefresh(c)
			} else {
				conf.RequirePermission(c)
			}
		} else if !allow {
			conf.RequirePermission(c)
		}
	}
}

func (a *JwtRoleStruct) RequireRefresh(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, model.ResponseError{
		Error: model.ServerError{
			Status: "UNAUTHORIZED",
			Message: "Token is expired",
		},
	})
	c.AbortWithStatus(401)
}

func (a *JwtRoleStruct) RequirePermission(c *gin.Context) {
	c.AbortWithStatus(403)
}

func (a *JwtRoleStruct) CheckPermission(r *http.Request) (bool, error) {
	role, err := a.GerRole(r)
	if err != nil {
		return false, err
	}
	method := r.Method
	path := r.URL.Path

	allowed, err := a.enforce.Enforce(role, path, method)
	if err != nil {
		panic(err)
	}
	return allowed, nil
}

func (a *JwtRoleStruct) GerRole(r *http.Request) (string, error) {
	var (
		role   string
		claims jwt.MapClaims
		err    error
	)
	jwtToken := r.Header.Get("Authorization")
	if jwtToken == "" {
		return "unauthorized", nil
	} else if strings.Contains(jwtToken, "Basic") { // basic - default token
		return "unauthorized", nil
	}

	a.jwtHandler.Token = jwtToken
	claims, err = a.jwtHandler.ExtractClaims()
	if err != nil {
		return "", err
	}

	if claims["role"].(string) == "authorized" {
		role = "authorized"
	} else {
		role = "unknown"
	}
	return role, nil
}
