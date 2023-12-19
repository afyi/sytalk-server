package middleware

import (
	"net/http"
	"strings"

	"github.com/afyi/sytalk/service"
	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1. 放在请求头 2. 放在请求体 3. 放在URI
		// 这里选2，放在Header的Authorization中，并用Bearer开头，就是普通的JWT
		// token 格式如下 "Bearer xfdsafdsafdsafdsa"
		authHeader := c.Request.Header.Get("Authorization")
		// 这里直接判定是否存在
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code": 40001,
				"msg":  "请求头中auth为空",
			})
			c.Abort()
			return
		}
		// 如果存在，则把token按空格分隔开
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 && parts[0] != "Bearer" {
			c.JSON(http.StatusOK, gin.H{
				"code": 40002,
				"msg":  "令牌格式不正确",
			})
			c.Abort()
			return
		}

		user, err := service.ParseToken(parts[1])

		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 40003,
				"msg":  "无效的Token",
			})
			c.Abort()
			return
		}
		// 将当前的用户信息保存到上下文对象中
		c.Set("user", user)
		c.Next()
	}
}
