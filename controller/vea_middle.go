package controller

import (
	"ai-smart/model"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

type VeaUserJwtInfo struct {
	jwt.StandardClaims
	Uid int64 `json:"uid"`
}

func genVeaJwt(uid int64) (string, error) {
	var ans VeaUserJwtInfo
	now := time.Now().Unix()
	secret := ""
	seconds := int64(123)

	ans.ExpiresAt = now + seconds
	ans.Uid = uid
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = ans
	return token.SignedString([]byte(secret))
}
func GenVeaJwt(pid, uid, fid int64) (string, error) {
	return genVeaJwt(uid)
}

func getVeaUserinfoFromToken(s, secret string) (*VeaUserJwtInfo, error) {
	var ans VeaUserJwtInfo
	_, err := jwt.ParseWithClaims(s, &ans, func(t *jwt.Token) (interface{}, error) { return []byte(secret), nil })
	if err != nil {
		return nil, errors.Wrapf(err, "get token parser err")
	}

	return &ans, ans.Valid()
}

func VeaJwtMiddle(needLogin bool) gin.HandlerFunc {
	secret := ""
	return func(c *gin.Context) {
		s := c.GetHeader("Authorization")
		user, err := getVeaUserinfoFromToken(s, secret)
		if err == nil {
			c.Set("uid", user.Uid)
			c.Set("user", user)
			return
		}

		if needLogin {
			ginJsonResponse(c, model.NeedLoginRsp)
		}
	}
}

func getVeaJwtUser(c *gin.Context) (*VeaUserJwtInfo, error) {
	v := c.Value("user")
	ans, ok := v.(*VeaUserJwtInfo)
	if ok {
		return ans, nil
	}

	return nil, errors.Errorf("未登录，type:%T,value:%+v", v, v)
}
