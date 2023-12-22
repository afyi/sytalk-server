package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/afyi/sytalk/model"
	"github.com/golang-jwt/jwt/v4"
)

// 自定义Claims
type MyClaims struct {
	Id   string `json:"userid"` // 用户id
	Name string `json:"name"`
	jwt.RegisteredClaims
}

// 密钥
var (
	MySecret []byte = []byte("hello world")
	Prefix   string = "Bearer "
)

// 生成jwt
func GenToken(user *model.UserClaims) (tokenString string, err error) {

	claim := MyClaims{

		Id: user.Id,

		Name: user.Name,

		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时

			IssuedAt: jwt.NewNumericDate(time.Now()), // 签发时间

			NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间

		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法

	tokenString, err = token.SignedString(MySecret)

	return fmt.Sprintf("%s%s", Prefix, tokenString), err
}

// 解析jwt
func ParseToken(tokenString string) (*MyClaims, error) {

	// 解析类似 "Bearer xxxxxxxxxxxxxxx" 类型的token
	newTokenString := strings.Replace(tokenString, Prefix, "", -1)

	claims := &MyClaims{}

	_, err := jwt.ParseWithClaims(newTokenString, claims, func(t *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})

	// 若token只是过期claims是有数据的，若token无法解析claims无数据
	return claims, err
}

func RenewToken(claims *MyClaims) (string, error) {

	userClaims := &model.UserClaims{Id: claims.Id, Name: claims.Name}

	// 若token过期不超过10分钟则给它续签
	if time.Now().Unix()-claims.ExpiresAt.Unix() < 600 {

		return GenToken(userClaims)

	}

	return "", errors.New("登陆已过期")
}
