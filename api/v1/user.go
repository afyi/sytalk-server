package v1

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/afyi/sytalk/database"
	"github.com/afyi/sytalk/model"
	"github.com/afyi/sytalk/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func Login(ctx *gin.Context) {

	defer close(ctx)

	name := strings.Trim(ctx.DefaultPostForm("name", ""), "")
	pass := strings.Trim(ctx.DefaultPostForm("pass", ""), "")

	if name == "" || pass == "" {
		panic("用户名或者密码不能为空")
	}

	cli, err := database.Connect("sytalk", "user", ctx1)

	if err != nil {
		// 防止数据库信息泄露
		panic("数据库连接错误!")
	}

	defer cli.Close(ctx1)

	var user model.User

	err = cli.Find(ctx1, bson.M{"name": name}).One(&user)

	if err != nil {
		panic("用户名密码错误或者不存在")
	}

	if user.GetSysMd5(pass, user.Salt) != user.Password {
		panic("用户名密码错误或者不存在")
	}

	// 然后授权JWT
	token, err := service.GenToken(user.Id.String())

	if err != nil {
		fmt.Println(err)
		panic("令牌生成错误")
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "ok",
		"data": map[string]string{
			"name":     user.Name,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"token":    token,
		},
	})

}
