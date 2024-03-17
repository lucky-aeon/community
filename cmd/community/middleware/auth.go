package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"time"
	"xhyovo.cn/community/pkg/result"
)

const AUTHORIZATION = "Authorization"

var stStringKey = []byte(viper.GetString("jwt.StringKey"))

type JwtCustomClaims struct {
	ID   int
	Name string
	jwt.RegisteredClaims
}

func GetUserId(ctx *gin.Context) int {

	return ctx.Value(AUTHORIZATION).(int)
}

func Auth(ctx *gin.Context) {
	token := ctx.GetHeader(AUTHORIZATION)
	if len(token) == 0 {
		token, _ = ctx.Cookie(AUTHORIZATION)
	}
	claims, err := ParseToken(token)
	if err != nil {
		result.Err(err.Error()).Json(ctx)
		ctx.Abort()
		return
	}
	ctx.Set(AUTHORIZATION, claims.ID)
	ctx.Next()
}
func GenerateToken(id int, name string) (string, error) {
	// 初始化
	iJwtCustomClaims := JwtCustomClaims{
		ID:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间 在当前基础上 添加一个小时后 过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			// 颁发时间 也就是生成时间
			IssuedAt: jwt.NewNumericDate(time.Now()),
			//主题
			Subject: "Token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, iJwtCustomClaims)

	return token.SignedString(stStringKey)
}

// ParseToken 解析token
func ParseToken(tokenStr string) (JwtCustomClaims, error) {

	iJwtCustomClaims := JwtCustomClaims{}
	if tokenStr == "" {
		return iJwtCustomClaims, errors.New("token为空")
	}
	token, err := jwt.ParseWithClaims(tokenStr, &iJwtCustomClaims, func(token *jwt.Token) (interface{}, error) {
		return stStringKey, nil
	})

	if err != nil || !token.Valid {
		err = errors.New("invalid Token")
	}
	return iJwtCustomClaims, err
}
