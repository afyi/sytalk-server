package middleware

import (
	"net/http"
	"strings"

	"github.com/afyi/sytalk/service"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 客户端携带Token有三种方式 1. 放在请求头 2. 放在请求体 3. 放在URI
		// 这里选2，放在Header的Authorization中，并用Bearer开头，就是普通的JWT
		// token 格式如下 "Bearer xfdsafdsafdsafdsa"
		auth := ctx.Request.Header.Get("Authorization")

		if len(auth) == 0 {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "未登陆无权限",
			})
			return
		}

		// 校检token，只要出错直接拒绝请求
		claims, err := service.ParseToken(auth)

		if err != nil {

			// 判定是否为过期
			if strings.Contains(err.Error(), "token is expired") {
				// 若过期，则续签
				newToken, _ := service.RenewToken(claims)

				if newToken != "" {
					// 续签成功，给返回头设置一个newtoken字体
					ctx.Header("newtoken", newToken)
					ctx.Request.Header.Set("Authorization", newToken)
					ctx.Next()
					return
				}
			}
			// 其它错误直接返回拒绝请求
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"code":    0,
				"message": "未登陆无权限",
			})
			return
		}
		// 将当前的用户信息保存到上下文对象中
		ctx.Set("user", claims)
		ctx.Next()
	}
}
