package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// 自定义Claims
type MyClaims struct {
	Id string `json:"userid"` // 用户id
	jwt.RegisteredClaims
}

// 密钥
var (
	MySecret []byte = []byte("hello world")
	Prefix   string = "Bearer"
)

// 生成jwt
func GenToken(userid string) (tokenString string, err error) {

	claim := MyClaims{

		Id: userid,

		RegisteredClaims: jwt.RegisteredClaims{

			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour * time.Duration(1))), // 过期时间3小时

			IssuedAt: jwt.NewNumericDate(time.Now()), // 签发时间

			NotBefore: jwt.NewNumericDate(time.Now()), // 生效时间

		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim) // 使用HS256算法

	tokenString, err = token.SignedString(MySecret)

	return fmt.Sprintf("%s %s", Prefix, tokenString), err
}

func Secret() jwt.Keyfunc {

	return func(token *jwt.Token) (interface{}, error) {

		return MySecret, nil // 这是我的secret

	}

}

// 解析jwt
func ParseToken(tokenString string) (*MyClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, Secret())

	if err != nil {

		if ve, ok := err.(*jwt.ValidationError); ok {

			if ve.Errors&jwt.ValidationErrorMalformed != 0 {

				return nil, errors.New("that's not even a token")

			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {

				return nil, errors.New("token is expired")

			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {

				return nil, errors.New("token not active yet")

			} else {

				return nil, errors.New("couldn't handle this token")

			}

		}

	}

	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {

		return claims, nil

	}

	return nil, errors.New("couldn't handle this token")

}
