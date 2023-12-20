package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/afyi/sytalk/service"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		// 客户端携带Token有三种方式 1. 放在请求头 2. 放在请求体 3. 放在URI
		// 这里选2，放在Header的Authorization中，并用Bearer开头，就是普通的JWT
		// token 格式如下 "Bearer xfdsafdsafdsafdsa"
		authHeader := ctx.Request.Header.Get("Authorization")

		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
				ctx.JSON(http.StatusOK, gin.H{"code": err, "message": "无效的Token"})
				ctx.Abort()
			}
		}()

		// 如果存在，则把token按空格分隔开
		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) < 2 {
			panic(40001)
		}

		// 直接验证
		userid, err := service.ParseToken(parts[1])

		if err != nil {
			panic(40001)
		}

		// 将当前的用户信息保存到上下文对象中
		ctx.Set("userid", userid)
		ctx.Next()
	}
}
