package router

import (
	v1 "github.com/afyi/sytalk/api/v1"
	"github.com/afyi/sytalk/middleware"
	"github.com/gin-gonic/gin"
)

func Init(r *gin.Engine) {

	apiv1 := r.Group("/api/v1")

	{
		// 获取说说列表
		apiv1.GET("/art", v1.GetArticleList)

	}

	apiv1.Use(middleware.JWTAuthMiddleware())

	{
		// 新增说说
		apiv1.POST("/art", v1.InsertArticle)

		// 修改说说
		apiv1.PUT("/art/:id", v1.UpdateArticle)

		// 删除说说
		apiv1.DELETE("/art/:id", v1.DeleteArticle)

	}

}
