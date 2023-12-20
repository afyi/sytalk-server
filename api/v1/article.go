package v1

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/afyi/sytalk/database"
	"github.com/afyi/sytalk/model"
	"github.com/gin-gonic/gin"
	"github.com/mileusna/useragent"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetArticleList(ctx *gin.Context) {

	defer close(ctx)

	// 当前页
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))

	// 每页条数
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	var article []model.Article

	cli, err := database.Connect("sytalk", "article", ctx1)

	if err != nil {
		panic("数据库连接错误!")
	}

	defer cli.Close(ctx1)

	total, err := cli.Find(ctx1, bson.M{}).Count()

	if err != nil {
		panic(err)
	}

	lastPage := math.Ceil(float64(total) / float64(pageSize))

	if total > 0 {
		// 起点
		offset := int64((page - 1) * pageSize)
		// 数量
		limit := int64(pageSize)

		if err := cli.Find(ctx1, bson.M{}).Sort("-createdAt").Skip(offset).Limit(limit).All(&article); err != nil {
			panic("查询数据异常")
		}
	}

	// 总页数至少为1
	if lastPage == 0 {
		lastPage = 1
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "ok",
		"data": map[string]interface{}{
			"data":     &article,
			"total":    total,
			"page":     page,
			"pageSize": pageSize,
			"lastPage": lastPage,
		},
	})
}

func UpdateArticle(ctx *gin.Context) {

	defer close(ctx)

	// 查询完记得释放
	var article model.Article

	if err := ctx.Bind(&article); err != nil {
		panic(err.Error())
	}

	cli, err := database.Connect("sytalk", "article", ctx1)

	if err != nil {
		panic("数据库连接错误!")
	}

	defer cli.Close(ctx1)

	// 把id类型转成objectid
	article.Id, err = primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		panic(err.Error())
	}

	ua := useragent.Parse(ctx.Request.UserAgent())

	// 更新的数据，因为只会更新部分数据，这里没有用到replaceOne方法，所以字段应该为真实字段，而非模型映射字段
	newData := bson.M{"$set": bson.M{"atContentHtml": article.Html, "atContentMd": article.Md, "os": fmt.Sprintf("%s %s", ua.OS, ua.OSVersion), "ua": fmt.Sprintf("%s %s", ua.Name, ua.Version), "updatedAt": time.Now()}}

	if err = cli.UpdateOne(ctx1, bson.M{"_id": article.Id}, newData); err != nil {
		panic(err.Error())
	}

	// 读取一遍新增数据，返回给前端
	if err = cli.Find(ctx1, bson.M{"_id": article.Id}).One(&article); err != nil {
		panic(err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "ok",
		"data":    &article,
	})
}

func InsertArticle(ctx *gin.Context) {

	// 查询完记得释放
	defer close(ctx)

	var article model.Article

	if err := ctx.ShouldBind(&article); err != nil {
		panic(err.Error())
	}

	cli, err := database.Connect("sytalk", "article", ctx1)

	if err != nil {
		// 防止数据库信息泄露
		panic("数据库连接错误!")
	}

	defer cli.Close(ctx1)

	ua := useragent.Parse(ctx.Request.UserAgent())

	// 更新的数据，因为只会更新部分数据，这里没有用到replaceOne方法，所以字段应该为真实字段，而非模型映射字段
	newData := &model.Article{
		Html:      article.Html,
		Md:        article.Md,
		Os:        fmt.Sprintf("%s %s", ua.OS, ua.OSVersion),
		Ua:        fmt.Sprintf("%s %s", ua.Name, ua.Version),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	res, err := cli.InsertOne(ctx1, newData)

	if err != nil {
		panic(err.Error())
	}

	// 读取一遍新增数据，返回给前端
	if err = cli.Find(ctx1, bson.M{"_id": res.InsertedID}).One(&article); err != nil {
		panic(err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "ok",
		"data":    &article,
	})
}

func DeleteArticle(ctx *gin.Context) {

	// 查询完记得释放
	defer close(ctx)

	// 把id类型转成objectid
	aid, err := primitive.ObjectIDFromHex(ctx.Param("id"))

	if err != nil {
		panic("参数错误")
	}

	cli, err := database.Connect("sytalk", "article", ctx1)

	if err != nil {
		panic("数据库连接错误!")
	}

	defer cli.Close(ctx1)

	// 删除数据
	err = cli.Remove(ctx1, bson.M{"_id": aid})

	// 不存在也算已删除
	if err != nil && err.Error() != "mongo: no documents in result" {
		panic(err.Error())
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    1,
		"message": "ok",
	})
}
