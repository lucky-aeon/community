package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
	"time"
	"xhyovo.cn/community/pkg/constant"
	"xhyovo.cn/community/pkg/result"
	services "xhyovo.cn/community/server/service"
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
	if claims.ID < 1 {
		result.Err("id 不正确").Json(ctx)
		ctx.Abort()
		return
	}

	var blackService = services.BlacklistService{}

	// 判断token 黑名单
	exist := blackService.ExistToken(token)
	if exist {
		result.Err("你已涉嫌违规社区文化，token 失效，请重新登陆").Json(ctx)
		ctx.Abort()
		return
	}

	// 判断用户黑名单
	var userService = services.UserService{}
	if userService.IsBlack(claims.ID) {
		result.Err("你已涉嫌违规社区文化，已被纳入小黑屋，如误封请联系我：xhyQAQ250").Json(ctx)
		ctx.Abort()
		return
	}

	ctx.Set(AUTHORIZATION, claims.ID)

	// 新增代码：检查 token 的剩余有效期
	now := time.Now()
	exp := claims.ExpiresAt.Time
	if exp.Sub(now) < 30*time.Minute {
		// 生成新的 token
		newToken, err := GenerateToken(claims.ID, claims.Name)
		if err == nil {
			// 在响应头中设置新的 token
			ctx.Header(AUTHORIZATION, newToken)
		}
	}

	ctx.Next()
}
func GenerateToken(id int, name string) (string, error) {
	// 初始化
	iJwtCustomClaims := JwtCustomClaims{
		ID:   id,
		Name: name,
		RegisteredClaims: jwt.RegisteredClaims{
			// 设置过期时间 在当前基础上 添加两个小时后 过期
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constant.Token_TTl)),
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
